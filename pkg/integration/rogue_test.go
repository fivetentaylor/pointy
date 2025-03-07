package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"

	"github.com/fivetentaylor/pointy/pkg/background/jobs"
	"github.com/fivetentaylor/pointy/pkg/background/wire"
	"github.com/fivetentaylor/pointy/pkg/config"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/server"
	"github.com/fivetentaylor/pointy/pkg/server/auth"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	"github.com/fivetentaylor/pointy/pkg/utils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

var NoMoreMessages = errors.New("no more messages")

type RogueTestSuite struct {
	suite.Suite

	server       *server.Server
	jwt          *auth.JWTEngine
	addr         string
	cancelWorker context.CancelFunc
}

func (suite *RogueTestSuite) SetupTest() {
	testutils.EnsureStorage()
	suite.cancelWorker = testutils.RunWorker(suite.T())
	gormdb := testutils.TestGormDb(suite.T())

	jwtsecret := "testtesttesttest"

	cfg := config.Server{
		Addr:            ":0",
		JWTSecret:       jwtsecret,
		OpenAIKey:       "sk-idontwork",
		AllowedOrigins:  []string{"https://*", "http://*"},
		WorkerRedisAddr: os.Getenv("REDIS_URL"),
		GoogleOauth: &config.GoogleOauth{
			ClientID:     "fake-client-id",
			ClientSecret: "fake-client-secret",
			RedirectURL:  "fake-redirect-url",
		},
	}

	s, err := server.NewServer(cfg, gormdb)
	if err != nil {
		suite.T().Fatal(fmt.Errorf("error creating server: %w", err))
	}

	// Start the server in a goroutine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			suite.T().Fatalf("Failed to start server: %v", err)
		}
	}()

	suite.jwt = auth.NewJWT(jwtsecret)

	time.Sleep(500 * time.Millisecond)

	suite.server = s
	suite.addr = s.Listener.Addr().String() // get the addr of the server
}

func (suite *RogueTestSuite) TearDownTest() {
	suite.cancelWorker()
}

func (suite *RogueTestSuite) TestSimpleDocumentInteraction() {
	t := suite.T()
	ctx := testutils.TestContext()
	user := testutils.CreateUser(t, ctx)
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "")
	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)

	// Generate auth token
	token, err := suite.jwt.GenerateUserToken(user)
	suite.Require().NoError(err)

	// Connect to WS
	u := url.URL{Scheme: "ws", Host: suite.addr, Path: fmt.Sprintf("/api/v1/documents/%s/rogue/ws", docID)}
	header := http.Header{}
	header.Add("Cookie", constants.CookieName+"="+token)
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	suite.Require().NoError(err)
	defer ws.Close()

	// Read messages into a channel
	messages := make(chan []byte)
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			log.Info("Read message", "message", string(msg))
			if err != nil {
				if cerr, ok := err.(*websocket.CloseError); ok {
					if cerr.Code == websocket.CloseAbnormalClosure {
						// test is completed, websocket is closed
						return
					}
				}

				log.Error("FAILED TO READ MESSAGE (might be caused by ws closing early)", "error", err)
			}
			messages <- msg
		}
	}()

	// Send subscribe
	sendMessage(t, ws, &rogue.Subscribe{
		Type:  "subscribe",
		DocID: docID,
	})
	// Receive auth
	msg, err := next(messages)
	suite.Require().NoError(err)
	authPayload := rogue.AuthEvent{}
	suite.Require().NoError(json.Unmarshal(msg, &authPayload))
	suite.Equal("1", authPayload.AuthorID) // should always be the first author
	// Receive snapshot
	msg, err = next(messages)
	suite.Require().NoError(err)
	snapshotPayload := v3.SnapshotOp{}
	suite.Require().NoError(json.Unmarshal(msg, &snapshotPayload))
	rogueInstance := v3.NewRogueForQuill(authPayload.AuthorID)
	rogueInstance.MergeOp(snapshotPayload)
	currentHtml, err := rogueInstance.GetFullHtml(true, false)
	suite.Require().NoError(err)
	suite.Equal("<p data-rid=\"q_1\"></p>", currentHtml)
	// Receive Loaded
	msg, err = next(messages)
	suite.Require().NoError(err)
	loadedPayload := map[string]string{}
	suite.Require().NoError(json.Unmarshal(msg, &loadedPayload))
	suite.Equal("loaded", loadedPayload["event"])
	// Receive newCursor
	msg, err = next(messages)
	suite.Require().NoError(err)
	newCursorPayload := map[string]interface{}{}
	suite.Require().NoError(json.Unmarshal(msg, &newCursorPayload))
	suite.Equal("newCursor", newCursorPayload["type"])
	suite.Equal("1", newCursorPayload["authorID"])
	// Check there are no more messages
	msg, err = next(messages)
	suite.Require().ErrorIs(err, NoMoreMessages, "expected no more messages but got %q", string(msg))

	text := "hello, world!"

	for i, c := range text {
		op, err := rogueInstance.Insert(i, string(c))
		suite.Require().NoError(err)

		opBytes, err := json.Marshal(op)
		suite.Require().NoError(err)

		sendMessage(t, ws, &rogue.Operation{
			Type: "op",
			Op:   string(opBytes),
		})
	}

	var count int
	for msg, err := next(messages); err != NoMoreMessages; msg, err = next(messages) {
		suite.Require().NoError(err)
		fmt.Printf("WS Received: %q\n", string(msg))
		count++
	}
	suite.Equal(13, count)

	// Check the HTML is correct
	html, err := suite.getDocHTML(docID, token)
	suite.Require().NoError(err)
	suite.Contains(html, "hello, world!")

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	size, err := ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(int64(13), size)

	// Run a Snapshot manually
	err = jobs.SnapshotRogueJob(ctx, &wire.SnapshotRogue{
		DocId: docID,
	})
	suite.Require().NoError(err, "SnapshotRogue failed for doc %q: %s", docID, err)

	// Keep checking the HTML while the job is running
	// break out after 100 tries
	for i := 0; i < 100; i++ {
		html, err := suite.getDocHTML(docID, token)
		suite.Require().NoError(err)
		suite.Contains(html, "hello, world!")
		time.Sleep(100 * time.Millisecond)

		size, err := ds.SizeDeltaLog(ctx, docID)
		suite.Require().NoError(err)
		if size == 0 {
			break
		}
	}

	size, err = ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(
		int64(0),
		size,
		"delta log should be empty after snapshot, got %d. This either because the snapshot failed or it took longer than 10s", size)
}

func (suite *RogueTestSuite) connectToDoc(docID string, user *models.User) (*websocket.Conn, chan []byte) {
	token, err := suite.jwt.GenerateUserToken(user)
	suite.Require().NoError(err)

	u := url.URL{Scheme: "ws", Host: suite.addr, Path: fmt.Sprintf("/api/v1/documents/%s/rogue/ws", docID)}
	header := http.Header{}
	header.Add("Cookie", constants.CookieName+"="+token)
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	suite.Require().NoError(err)

	// Read messages into a channel
	messages := make(chan []byte)
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				if cerr, ok := err.(*websocket.CloseError); ok {
					if cerr.Code == websocket.CloseAbnormalClosure {
						// test is complete
						return
					}
				}

				log.Error("FAILED TO READ MESSAGE (might be caused by ws closing early)", "error", err)
			}
			messages <- msg
		}
	}()

	return ws, messages
}

func (suite *RogueTestSuite) connectAndSubscribeToDoc(docID string, user *models.User) (*websocket.Conn, chan []byte, *v3.Rogue) {
	ws, messages := suite.connectToDoc(docID, user)

	// Send subscribe
	sendMessage(suite.T(), ws, &rogue.Subscribe{
		Type:  "subscribe",
		DocID: docID,
	})
	// Receive auth
	msg, err := next(messages)
	suite.Require().NoError(err)
	authPayload := rogue.AuthEvent{}
	suite.Require().NoError(json.Unmarshal(msg, &authPayload))
	// Receive snapshot
	msg, err = next(messages)
	suite.Require().NoError(err)
	snapshotPayload := v3.SnapshotOp{}
	suite.Require().NoError(json.Unmarshal(msg, &snapshotPayload))
	rogueInstance := v3.NewRogueForQuill(authPayload.AuthorID)
	rogueInstance.MergeOp(snapshotPayload)

	// Receive Loaded
	msg, err = next(messages)
	suite.Require().NoError(err)
	loadedPayload := map[string]string{}
	suite.Require().NoError(json.Unmarshal(msg, &loadedPayload), "expected loaded but got %q", string(msg))
	suite.Equal("loaded", loadedPayload["event"])

	for _, err := next(messages); err != NoMoreMessages; _, err = next(messages) {
		if errors.Is(err, NoMoreMessages) {
			break
		}
		suite.Require().NoError(err)
	}

	return ws, messages, rogueInstance
}

func (suite *RogueTestSuite) TestMultipleAuthors() {
	t := suite.T()
	ctx := testutils.TestContext()
	user1 := testutils.CreateUser(t, ctx)
	user2 := testutils.CreateUser(t, ctx)
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "")
	testutils.AddOwnerToDocument(t, ctx, docID, user1.ID)
	testutils.AddOwnerToDocument(t, ctx, docID, user2.ID)
	ws1, messages1, rogueInstance1 := suite.connectAndSubscribeToDoc(docID, user1)
	ws2, messages2, rogueInstance2 := suite.connectAndSubscribeToDoc(docID, user2)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		text := "hello, world!"
		for i, c := range text {
			op, err := rogueInstance1.Insert(i, string(c))
			suite.Require().NoError(err)

			opBytes, err := json.Marshal(op)
			suite.Require().NoError(err)

			sendMessage(t, ws1, &rogue.Operation{
				Type: "op",
				Op:   string(opBytes),
			})
		}

		// Receive all messages
		for msg, err := next(messages1); err != NoMoreMessages; msg, err = next(messages1) {
			suite.Require().NoError(err)
			fmt.Printf("WS Received: %q\n", string(msg))
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		text := "goodbye, moon!"
		for i, c := range text {
			op, err := rogueInstance2.Insert(i, string(c))
			suite.Require().NoError(err)

			opBytes, err := json.Marshal(op)
			suite.Require().NoError(err)

			sendMessage(t, ws2, &rogue.Operation{
				Type: "op",
				Op:   string(opBytes),
			})
		}

		// Receive all messages
		for msg, err := next(messages2); err != NoMoreMessages; msg, err = next(messages2) {
			suite.Require().NoError(err)
			fmt.Printf("WS Received: %q\n", string(msg))
		}
		wg.Done()
	}()
	wg.Wait()

	// Check the HTML is correct
	token, err := suite.jwt.GenerateUserToken(user1)
	suite.Require().NoError(err)
	html, err := suite.getDocHTML(docID, token)
	suite.Require().NoError(err)
	suite.Contains(html, "hello, world!")
	suite.Contains(html, "goodbye, moon!")

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	size, err := ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(int64(27), size)

	// Run a Snapshot manually
	err = jobs.SnapshotRogueJob(ctx, &wire.SnapshotRogue{
		DocId: docID,
	})
	suite.Require().NoError(err, "SnapshotRogue failed for doc %q: %s", docID, err)

	// Keep checking the HTML while the job is running
	// break out after 100 tries
	for i := 0; i < 100; i++ {
		html, err := suite.getDocHTML(docID, token)
		suite.Require().NoError(err)
		suite.Contains(html, "hello, world!")
		time.Sleep(100 * time.Millisecond)

		size, err := ds.SizeDeltaLog(ctx, docID)
		suite.Require().NoError(err)
		if size == 0 {
			break
		}
	}

	size, err = ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(
		int64(0),
		size,
		"delta log should be empty after snapshot, got %d. This either because the snapshot failed or it took longer than 10s", size)
}

func (suite *RogueTestSuite) TestLargeDocument() {
	t := suite.T()
	t.Skip("Skipping large document test")
	ctx := testutils.TestContext()
	user := testutils.CreateUser(t, ctx)
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "")
	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)
	ws, messages, rogueInstance := suite.connectAndSubscribeToDoc(docID, user)

	for i := 0; i < 10000; i++ {
		if i%10 == 0 {
			op, err := rogueInstance.Insert(i, "\n")
			suite.Require().NoError(err)
			opBytes, err := json.Marshal(op)
			suite.Require().NoError(err)

			sendMessage(t, ws, &rogue.Operation{
				Type: "op",
				Op:   string(opBytes),
			})
			continue
		}
		s := utils.RandomSafeString(10)
		op, err := rogueInstance.Insert(i, s)
		suite.Require().NoError(err)
		opBytes, err := json.Marshal(op)
		suite.Require().NoError(err)

		sendMessage(t, ws, &rogue.Operation{
			Type: "op",
			Op:   string(opBytes),
		})
	}

	// Receive all messages
	for msg, err := next(messages); err != NoMoreMessages; msg, err = next(messages) {
		suite.Require().NoError(err)
		fmt.Printf("WS Received: %q\n", string(msg))
	}

	// Check the HTML is correct
	token, err := suite.jwt.GenerateUserToken(user)
	suite.Require().NoError(err)
	_, err = suite.getDocHTML(docID, token)
	suite.Require().NoError(err)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	_, err = ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)

	// Run a Snapshot manually
	err = jobs.SnapshotRogueJob(ctx, &wire.SnapshotRogue{
		DocId: docID,
	})
	suite.Require().NoError(err, "SnapshotRogue failed for doc %q: %s", docID, err)

	size, err := ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(
		int64(0),
		size,
		"delta log should be empty after snapshot, got %d. This either because the snapshot failed or it took longer than 10s", size)
}

func (suite *RogueTestSuite) TestConcurrentWebsocketWrites() {
	t := suite.T()
	t.Skip("Skipping large document test")
	ctx := testutils.TestContext()
	user := testutils.CreateUser(t, ctx)
	docID := uuid.NewString()
	testutils.CreateTestDocument(t, ctx, docID, "")
	testutils.AddOwnerToDocument(t, ctx, docID, user.ID)
	ws, messages, rogueInstance := suite.connectAndSubscribeToDoc(docID, user)

	opCount := 5000
	ops := make([]v3.Op, 0, opCount)

	for i := 0; i < opCount; i++ {
		s := utils.RandomSafeString(10)
		op, err := rogueInstance.Insert(i, s)
		suite.Require().NoError(err)
		ops = append(ops, op)
	}

	for _, op := range ops {
		go func(op v3.Op) {
			opBytes, err := json.Marshal(op)
			suite.Require().NoError(err)

			sendMessage(t, ws, &rogue.Operation{
				Type: "op",
				Op:   string(opBytes),
			})
		}(op)
	}

	// Receive all messages
	for msg, err := next(messages); err != NoMoreMessages; msg, err = next(messages) {
		suite.Require().NoError(err)
		fmt.Printf("WS Received: %q\n", string(msg))
	}

	// Check the HTML is correct
	token, err := suite.jwt.GenerateUserToken(user)
	suite.Require().NoError(err)
	html, err := suite.getDocHTML(docID, token)
	suite.Require().NoError(err)
	iHtml, err := rogueInstance.GetFullHtml(true, false)
	suite.Require().NoError(err)
	suite.Equal(iHtml, html)

	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	ds := rogue.NewDocStore(s3, query, redis)
	_, err = ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)

	// Run a Snapshot manually
	err = jobs.SnapshotRogueJob(ctx, &wire.SnapshotRogue{
		DocId: docID,
	})
	suite.Require().NoError(err, "SnapshotRogue failed for doc %q: %s", docID, err)

	size, err := ds.SizeDeltaLog(ctx, docID)
	suite.Require().NoError(err)
	suite.Equal(
		int64(0),
		size,
		"delta log should be empty after snapshot, got %d. This either because the snapshot failed or it took longer than 10s", size)
}

func (suite *RogueTestSuite) getDocHTML(docID, token string) (string, error) {
	u := url.URL{Scheme: "http", Host: suite.addr, Path: fmt.Sprintf("/api/v1/documents/%s/doc.html", docID)}
	var body []byte
	var err error

	// Retry 10 times
	for i := 0; i < 10; i++ {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Add("Cookie", constants.CookieName+"="+token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if i == 9 {
				return "", fmt.Errorf("failed to send request after 10 retries: %w", err)
			}
			time.Sleep(time.Second) // Sleep between retries
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			if i == 9 {
				return "", fmt.Errorf("unexpected status code after 10 retries: %d", res.StatusCode)
			}
			time.Sleep(time.Second) // Sleep between retries
			continue
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			if i == 9 {
				return "", fmt.Errorf("failed to read response body after 10 retries: %w", err)
			}
			time.Sleep(time.Second) // Sleep between retries
			continue
		}

		return string(body), nil
	}

	// If we reach here, something went wrong in all retries
	return "", err
}

func TestRogueTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RogueTestSuite))
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg interface{}) {
	if err := ws.WriteJSON(msg); err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}
}

func next(messages <-chan []byte) ([]byte, error) {
	select {
	case msg := <-messages:
		return msg, nil
	case <-time.After(100 * time.Millisecond):
		return nil, NoMoreMessages
	}
}
