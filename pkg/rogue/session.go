package rogue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jpoz/conveyor"

	"github.com/teamreviso/code/pkg/background/wire"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	"github.com/teamreviso/code/pkg/service/document"
	"github.com/teamreviso/code/pkg/service/timeline"
	"github.com/teamreviso/code/pkg/storage/dynamo"
	rogueV3 "github.com/teamreviso/code/rogue/v3"
)

const loadedEvent = "{\"type\":\"event\", \"event\":\"loaded\"}"
const pingEvent = "{\"type\":\"event\", \"event\":\"ping\"}"

type AuthEvent struct {
	Type     string `json:"type"`
	AuthorID string `json:"authorID"`
}

type Event struct {
	Type  string                 `json:"type"`
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

type Session struct {
	ID       string
	docID    string
	document *models.Document
	user     *models.User
	doc      *rogueV3.Rogue
	conn     *websocket.Conn
	store    *DocStore
	realtime *Realtime
	cancel   context.CancelFunc
	docLog   *Logger
	log      *slog.Logger
	query    *query.Query
	lastPong time.Time
	authorID string

	pongMutex sync.Mutex
	connMutex sync.Mutex
	docMutex  sync.Mutex
}

func NewSession(
	ctx context.Context,
	user *models.User,
	conn *websocket.Conn,
	store *DocStore,
	docID string,
) (*Session, error) {
	q := env.Query(ctx)

	doc, err := query.GetEditableDocumentForUser(q, docID, user.ID)
	if err != nil {
		return nil, err
	}

	sessionId := uuid.NewString()

	session := &Session{
		ID:       sessionId,
		docID:    docID,
		document: doc,
		user:     user,
		conn:     conn,
		store:    store,
		cancel:   nil,
		query:    q,
		lastPong: time.Now(),
	}

	session.realtime = NewRealtime(
		store.Redis,
		store.Query,
		docID,
		session.UserID(),
		session.UserName(),
		session.HighlightColor(),
	)

	log := log.With("docID", docID)
	session.docLog = NewLogger(log, store.S3, docID, fmt.Sprintf("%s|%s", user.ID[0:5], sessionId[0:5]))
	session.docLog.DeactivatePassthrough() // disable passthrough remove this if you want the log to go to stdout
	session.log = env.SLog(ctx).With("docID", docID, "sessionID", sessionId)

	go session.Keepalive()

	return session, nil
}

type MsgType struct {
	Type string `json:"type"`
}

type Subscribe struct {
	Type     string `json:"type"`
	DocID    string `json:"docID"`
	AuthorID string `json:"authorID,omitempty"`
}

// {"id":["cr8vxds2",153],"char":"!","parentId":["cr8vxds2",152],"side":1}
type Operation struct {
	Type string `json:"type"`
	Op   string `json:"op"`
}

type Op struct {
	ID       ID `json:"id"`
	ParentID ID `json:"parentId"`
}

type ID struct {
	Author    string
	LamportID int
}

func IDFromString(s string) (ID, error) {
	var id ID
	err := json.Unmarshal([]byte(s), &id)
	return id, err
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var temp []interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if len(temp) != 2 {
		return fmt.Errorf("expected 2 elements in the array, got %d", len(temp))
	}

	var ok bool
	if id.Author, ok = temp[0].(string); !ok {
		return fmt.Errorf("first element is not a string")
	}

	floatValue, ok := temp[1].(float64) // JSON numbers are floats
	if !ok {
		return fmt.Errorf("second element is not an integer")
	}
	id.LamportID = int(floatValue)

	return nil
}

// MarshalJSON custom marshaler for ID
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{id.Author, id.LamportID})
}

func (s *Session) UserID() string {
	return s.user.ID
}

func (s *Session) UserName() string {
	return s.user.ShortName()
}

func (s *Session) HighlightColor() string {
	if s.user == nil {
		return constants.DefaultHighlightColor
	}
	return s.user.HighlightColor()
}

func (s *Session) DeactivateDocLogger() {
	s.docLog.Deactivate()
}

func (s *Session) Keepalive() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.log.Info("ping. last pong", "duration", time.Now().Sub(s.lastPong).Seconds())
			err := s.writeMessage([]byte(pingEvent))
			if err != nil {
				return
			}
		}
	}
}

func (s *Session) handleSubscribe(ctx context.Context, msg []byte) error {
	s.docMutex.Lock()
	defer s.docMutex.Unlock()
	sub := Subscribe{}
	err := json.Unmarshal(msg, &sub)
	if err != nil {
		s.log.Error("error unmarshalling subscribe", "error", err)
		return err
	}

	s.log.Info("handling subscribe", "docID", sub.DocID, "authorID", sub.AuthorID, "message", string(msg), "userID", s.UserID())

	if sub.DocID != s.docID {
		s.log.Error("invalid docID: %s", "docID", sub.DocID)
		return fmt.Errorf("invalid docID: %s", sub.DocID)
	}

	if sub.AuthorID != "" {
		valid, err := document.ValidateAuthorID(ctx, sub.AuthorID, sub.DocID, s.UserID())
		if err != nil || !valid {
			s.log.Warn("invalid authorID", "authorID", sub.AuthorID, "userID", s.UserID(), "docID", sub.DocID)
			sub.AuthorID = ""
		}
	}

	if sub.AuthorID == "" {
		sub.AuthorID, err = document.NewAuthorID(ctx, sub.DocID, s.user.ID)
		if err != nil {
			return fmt.Errorf("error creating authorID: %w", err)
		}
	}

	start := time.Now()

	_, doc, err := s.store.GetCurrentDoc(ctx, sub.DocID)
	if err != nil {
		s.log.Error("error getting doc", "error", err)
		return err
	}

	s.doc = doc
	s.authorID = sub.AuthorID

	// Give the client a unique author id
	authBytes, err := json.Marshal(AuthEvent{
		Type:     "auth",
		AuthorID: sub.AuthorID,
	})
	if err != nil {
		s.log.Error("error marshalling auth event", "error", err)
		return err
	}
	err = s.writeMessage(authBytes)
	if err != nil {
		s.log.Error("error writing auth event", "error", err)
		return err
	}

	// Send the snapshot
	snapshot, err := doc.NewSnapshotOp()
	if err != nil {
		s.log.Error("error creating snapshot operation", "error", err)
		return err
	}
	opBytes, err := json.Marshal(snapshot)
	if err != nil {
		s.log.Error("error marshalling snapshot", "error", err)
		return err
	}
	err = s.writeMessage(opBytes)
	if err != nil {
		s.log.Error("error marshalling snapshot", "error", err)
		return err
	}

	// Send the loaded event
	err = s.writeMessage([]byte(loadedEvent))
	if err != nil {
		s.log.Error("error writing loaded event", "error", err)
		return err
	}

	// Connect to the realtime channel
	s.cancel = s.realtime.Subscribe(ctx, sub.AuthorID, s.onRealtimeMessage)

	// Send the active cursors
	cursors, err := s.realtime.CurrentRealtimeOperations(ctx)
	if err != nil {
		s.log.Error("error getting cursors", "error", err)
		return err
	}
	for _, cursor := range cursors {
		cursorBytes, err := json.Marshal(cursor)
		if err != nil {
			s.log.Error("error marshalling cursor", "error", err)
			return err
		}
		err = s.writeMessage(cursorBytes)
		if err != nil {
			s.log.Error("error writing cursor", "error", err)
			return err
		}
	}

	elapsed := time.Since(start)
	s.log.Info("subscribed for docID", "docID", sub.DocID, "elapsed", elapsed)

	return nil
}

func (s *Session) onRealtimeMessage(msg []byte) {
	s.docMutex.Lock()
	defer s.docMutex.Unlock()

	err := s.writeMessage(msg)
	if err != nil {
		s.log.Error("failed to write message", "err", err)
		return
	}
}

func (s *Session) handleOp(ctx context.Context, msg []byte) error {
	s.docMutex.Lock()
	defer s.docMutex.Unlock()

	log.Debugf("[%s] handling op: %s", s.docID, string(msg))
	var op Operation
	err := json.Unmarshal(msg, &op)
	if err != nil {
		return fmt.Errorf("error unmarshalling op: %s", err)
	}

	seq, err := s.store.AddDeltaLog(ctx, s.docID, op.Op)
	if err != nil {
		return fmt.Errorf("s.store.AddDeltaLog(ctx, %s, %s): %w", s.docID, op.Op, err)
	}

	if seq%1000 == 0 {
		_, err = env.Background(ctx).Enqueue(ctx, &wire.SnapshotRogue{
			DocId: s.docID,
		})
		if err != nil {
			log.Errorf("[messaging] error enqueueing job: %s", err)
		}
	}

	log.Debugf("[%s] merged op: %s", s.docID, string(msg))

	err = s.realtime.PublishOp(s.docID, []byte(op.Op))
	if err != nil {
		log.Errorf("error publishing op: %s", err)
		return fmt.Errorf("error publishing op: %s", err)
	}

	// if the doc has been updated in the last minute, update the document record
	if s.document.UpdatedAt.Before(time.Now().Add(-1 * time.Minute)) {
		docTbl := s.query.Document
		_, err = docTbl.
			Where(docTbl.ID.Eq(s.docID)).
			Updates(map[string]interface{}{
				"updated_at": time.Now(),
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) handleCursorUpdate(ctx context.Context, msg []byte) error {
	return s.realtime.PublishCursorUpdate(ctx, msg)
}

func (s *Session) handleEvent(ctx context.Context, msg []byte) error {
	var e Event
	err := json.Unmarshal(msg, &e)
	if err != nil {
		return fmt.Errorf("error unmarshalling event: %s", err)
	}

	if e.Event == "pong" {
		s.pongMutex.Lock()
		defer s.pongMutex.Unlock()
		s.lastPong = time.Now()
	}

	if e.Event == "paste" {
		var ok bool
		var val interface{}
		var before string

		if val, ok = e.Data["contentAddressBefore"]; ok {
			if before, ok = val.(string); !ok {
				log.Errorf("error getting before content address from paste event: %s %#v", err, e.Data)
			}
		} else {
			log.Errorf("error getting before content address key from paste event: %s %#v", err, e.Data)
		}

		var after string
		if val, ok = e.Data["contentAddressAfter"]; ok {
			if after, ok = val.(string); !ok {
				log.Errorf("error getting after content address from paste event: %s %#v", err, e.Data)
			}
		} else {
			log.Errorf("error getting after content address key from paste event: %s %#v", err, e.Data)
		}

		timeline.CreateTimelineEvent(ctx, &dynamo.TimelineEvent{
			UserID:   s.UserID(),
			AuthorID: s.authorID,
			DocID:    s.docID,
			Event: &models.TimelineEventPayload{
				Payload: &models.TimelineEventPayload_Paste{
					Paste: &models.TimelinePaste{
						ContentAddressBefore: before,
						ContentAddressAfter:  after,
					},
				},
			},
		})
	}

	return nil
}

func (s *Session) Message(ctx context.Context, msg []byte) error {
	t := MsgType{}

	err := json.Unmarshal(msg, &t)
	if err != nil {
		s.log.Error("rogue session could not parse msg", "msg", string(msg), "error", err)
		return fmt.Errorf("rogue session could not parse msg %s: %w", string(msg), err)
	}

	err = env.Redis(ctx).Set(ctx, fmt.Sprintf(constants.DocUserLastMessageKey, s.docID, s.user.ID), time.Now().Unix(), time.Minute).Err()
	if err != nil {
		s.log.Error("error setting last message time", "error", err)
	}

	if t.Type == "subscribe" {
		s.docLog.Info(fmt.Sprintf("-> %s", msg))
		return s.handleSubscribe(ctx, msg)
	} else if t.Type == "op" {
		s.docLog.Info(fmt.Sprintf("-> %s", msg))
		return s.handleOp(ctx, msg)
	} else if t.Type == "cursor" {
		return s.handleCursorUpdate(ctx, msg)
	} else if t.Type == "event" {
		return s.handleEvent(ctx, msg)
	}

	s.log.Error("unknown message type: %s", "type", t.Type)
	return fmt.Errorf("unknown message type: %s", t.Type)
}

func (s *Session) writeMessage(msg []byte) error {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()
	go func() {
		if strings.HasPrefix(string(msg), `{"type":"cursor"`) {
			return
		}
		s.docLog.Info(fmt.Sprintf("<- %s", msg))
	}()
	return s.conn.WriteMessage(websocket.TextMessage, msg)
}

func (s *Session) Close(ctx context.Context) error {
	s.log.Info("closing session for docID", "docID", s.docID)
	s.docLog.Close()
	if s.cancel != nil {
		s.log.Info("cancelling session for docID", "docID", s.docID)
		s.cancel()
		s.cancel = nil
	}

	if s.doc == nil {
		return nil
	}

	lastMsg, err := env.Redis(ctx).Get(ctx, fmt.Sprintf(constants.DocUserLastMessageKey, s.docID, s.user.ID)).Int64()
	if err != nil {
		s.log.Error("error getting last message time", "error", err)
	}

	env.Background(ctx).Enqueue(ctx, &wire.SummarizeSession{
		SessionId:       s.ID,
		DocId:           s.docID,
		UserId:          s.user.ID,
		LastMessageTime: lastMsg,
	}, conveyor.Delay(10*time.Second))

	s.doc = nil // free up doc memory
	return s.conn.Close()
}
