package rogue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
	"github.com/teamreviso/code/pkg/query"
	v3 "github.com/teamreviso/code/rogue/v3"
)

const PresenceCheckInterval = 5 * time.Second

type Realtime struct {
	docID    string
	userID   string
	authorID string
	name     string
	color    string

	rds *redis.Client
	qry *query.Query
}

type DocCursorOperation struct {
	Type     string  `json:"type,omitempty"` // newCursor | cursor
	UserID   string  `json:"userID,omitempty"`
	AuthorID string  `json:"authorID,omitempty"`
	Name     string  `json:"name,omitempty"`
	Color    string  `json:"color,omitempty"`
	Range    []v3.ID `json:"range,omitempty"`
	Editing  bool    `json:"editing"`
}

func NewRealtime(rds *redis.Client, qry *query.Query, docID, userID, name, color string) *Realtime {
	return &Realtime{
		docID:  docID,
		userID: userID,
		name:   name,
		color:  color,

		rds: rds,
		qry: qry,
	}
}

func (rt *Realtime) Subscribe(ctx context.Context, authorID string, onMsg func(msg []byte)) context.CancelFunc {
	log := env.Log(ctx)
	ctx, cancel := context.WithCancel(ctx)
	pubsub := SubscribeToDoc(ctx, rt.rds, rt.docID)
	incoming := pubsub.Channel()
	rt.authorID = authorID

	go func() {
		defer func() {
			bctx := context.Background()
			log.Info("🛑 closing realtime listener", "doc", rt.docID)
			rt.Disconnect(bctx, log)
			pubsub.Unsubscribe(bctx) // not sure if this is needed
			pubsub.Close()
		}()

		// Publish the new cursor
		err := rt.PublishNewDocCursor(ctx)
		if err != nil {
			log.Warn("error publishing new cursor", "doc", rt.docID, "err", err)
		}

		err = AddAuthorToActiveConnections(ctx, rt.rds, rt.docID, rt.userID, rt.authorID)
		if err != nil {
			log.Warn("error adding author to connections set", "doc", rt.docID, "err", err, "author", rt.authorID)
		}

		for {
			err = ExtendAuthorLastCursor(ctx, rt.rds, rt.docID, rt.userID, rt.authorID)
			timer := time.NewTimer(PresenceCheckInterval - time.Second)

			select {
			case rm := <-incoming: // Message from redis
				log.Debug("received message", "doc", rt.docID)
				onMsg([]byte(rm.Payload))
			case <-timer.C:
				log.Debug("resetting presance", "doc", rt.docID)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

func (rt *Realtime) Disconnect(ctx context.Context, log *log.Logger) {
	err := RemoveAuthorFromActiveConnections(ctx, rt.rds, rt.docID, rt.userID, rt.authorID)
	if err != nil {
		log.Error("error removing author from connections set", "doc", rt.docID, "err", err, "author", rt.authorID)
	}

	err = RemoveAuthorLastCursor(ctx, rt.rds, rt.docID, rt.userID, rt.authorID)
	if err != nil {
		log.Error("error deleting cursor", "doc", rt.docID, "err", err, "author", rt.authorID)
	}

	err = rt.PublishDeleteDocCursor(ctx)
	if err != nil {
		log.Error("error deleting cursor", "doc", rt.docID, "err", err, "author", rt.authorID)
	}
}

func (rt *Realtime) CursorForUser(ctx context.Context, userID, authorID string) (*DocCursorOperation, error) {
	bts, err := GetAuthorLastCursor(ctx, rt.rds, rt.docID, userID, authorID)
	if err != nil {
		return nil, err
	}

	cursor := DocCursorOperation{}
	err = json.Unmarshal(bts, &cursor)
	if err != nil {
		log.Errorf("error unmarshalling cursor: %s %s %s", userID, authorID, err)
		return nil, err
	}

	return &cursor, nil
}

func (rt *Realtime) CurrentRealtimeOperations(ctx context.Context) ([]DocCursorOperation, error) {
	log := env.Log(ctx)
	members, err := CurrentActiveConnections(ctx, rt.rds, rt.docID)
	if err != nil {
		return nil, fmt.Errorf("error getting active connections: %s", err)
	}

	users := []string{}
	authors := map[string]string{}
	for _, member := range members {
		parts := strings.Split(member, ":")
		userID := parts[0]
		authorID := parts[1]

		users = append(users, userID)
		authors[authorID] = userID
	}

	userTbl := rt.qry.User
	dbUsers, err := userTbl.Where(userTbl.ID.In(users...)).Find()
	if err != nil {
		return nil, fmt.Errorf("error querying users: %s", err)
	}

	userMap := map[string]*models.User{}
	for _, user := range dbUsers {
		userMap[user.ID] = user
	}

	cursors := []DocCursorOperation{}
	for authorID, userID := range authors {
		user, ok := userMap[userID]
		if !ok {
			log.Warnf("user no longer exists: %s %s. Removing connection", userID, authorID)
			ierr := RemoveAuthorFromActiveConnections(ctx, rt.rds, rt.docID, userID, authorID)
			if ierr != nil {
				log.Errorf("error removing author from connections set: %s %s %s", userID, authorID, ierr)
			}
			continue
		}

		bts, err := GetAuthorLastCursor(ctx, rt.rds, rt.docID, userID, authorID)
		if err != nil {
			log.Warn("no cursor!!! Removing connection", "docID", rt.docID, "userID", userID, "authorID", authorID, "err", err)
			ierr := RemoveAuthorFromActiveConnections(ctx, rt.rds, rt.docID, userID, authorID)
			if ierr != nil {
				log.Errorf("error removing author from connections set: %s %s %s", userID, authorID, ierr)
			}
			continue
		}

		// Unmarshal the last operation
		cursor := DocCursorOperation{}
		err = json.Unmarshal(bts, &cursor)
		if err != nil {
			log.Errorf("error unmarshalling cursor: %s %s %s", userID, authorID, err)
			continue
		}

		// Update the cursor with the user's name and highlight color
		cursor.Type = "newCursor"
		cursor.UserID = userID
		cursor.AuthorID = authorID
		cursor.Name = user.ShortName()
		cursor.Color = user.HighlightColor()

		cursors = append(cursors, cursor)
	}

	return cursors, nil
}

func (rt *Realtime) PublishNewDocCursor(ctx context.Context) error {
	docCursor := DocCursorOperation{
		Type:     "newCursor",
		UserID:   rt.userID,
		AuthorID: rt.authorID,
		Name:     rt.name,
		Color:    rt.color,
	}
	op, err := json.Marshal(docCursor)
	if err != nil {
		return fmt.Errorf("error marshalling doc cursor: %s", err)
	}

	err = AddAuthorLastCursor(ctx, rt.rds, rt.docID, rt.userID, rt.authorID, op)
	if err != nil {
		return fmt.Errorf("error adding author last cursor: %s", err)
	}

	return PublishToDoc(ctx, rt.rds, rt.docID, op)
}

func (rt *Realtime) PublishDeleteDocCursor(ctx context.Context) error {
	docCursor := DocCursorOperation{
		Type:     "deleteCursor",
		UserID:   rt.userID,
		AuthorID: rt.authorID,
	}
	op, err := json.Marshal(docCursor)
	if err != nil {
		return fmt.Errorf("error marshalling doc cursor: %s", err)
	}

	return PublishToDoc(ctx, rt.rds, rt.docID, op)
}

func (rt *Realtime) PublishOp(docID string, op []byte) error {
	return rt.rds.Publish(context.Background(), fmt.Sprintf(constants.DocUpdateChanFormat, docID), op).Err()
}

func (rt *Realtime) PublishCursorUpdate(ctx context.Context, cursorUpdate []byte) error {
	cursor := DocCursorOperation{}
	err := json.Unmarshal(cursorUpdate, &cursor)
	if err != nil {
		log.Errorf("error unmarshalling cursor: %s %s %s", rt.userID, rt.authorID, err)
		return err
	}

	cursor.UserID = rt.userID
	cursor.AuthorID = rt.authorID
	cursor.Name = rt.name
	cursor.Color = rt.color

	_, err = rt.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		cursorUpdate, err = json.Marshal(cursor)
		if err != nil {
			return fmt.Errorf("error marshalling cursor: %s", err)
		}

		err := AddAuthorLastCursor(ctx, pipe, rt.docID, rt.userID, rt.authorID, cursorUpdate)
		if err != nil {
			return fmt.Errorf("error adding author last cursor: %s", err)
		}

		err = pipe.Publish(context.Background(), fmt.Sprintf(constants.DocUpdateChanFormat, rt.docID), cursorUpdate).Err()
		if err != nil {
			return fmt.Errorf("error publishing cursor update: %s", err)
		}

		return nil
	})

	return err
}
