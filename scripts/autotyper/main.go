package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/server/auth"
	"github.com/teamreviso/code/pkg/testutils"
	v3 "github.com/teamreviso/code/rogue/v3"
)

const editInterval = 2 * time.Second
const typingDelay = 50 * time.Millisecond

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run script.go <doc_id> <user@email.com>")
		os.Exit(1)
	}
	docID := os.Args[1]
	email := os.Args[2]

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	autotyper := NewAutoTyper(docID, email)
	autotyper.Start()
}

type AutoTyper struct {
	DocID string
	Email string
	Doc   *v3.Rogue

	authorID string
	ctx      context.Context
	conn     *websocket.Conn
}

func NewAutoTyper(docID, email string) *AutoTyper {
	ctx := testutils.TestContext()

	return &AutoTyper{
		DocID: docID,
		Email: email,
		ctx:   ctx,
	}
}

func (at *AutoTyper) Start() {
	fmt.Printf("Connecting to WebSocket...\ndocID: %s\nemail: %s\n", at.DocID, at.Email)

	urlStr := fmt.Sprintf("wss://app.reviso.dev:9090/api/v1/documents/%s/rogue/ws", at.DocID)

	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}

	jwt := auth.NewJWT(os.Getenv("JWT_SECRET"))

	userTbl := env.Query(at.ctx).User
	user, err := userTbl.Where(userTbl.Email.Eq(at.Email)).First()
	if err != nil {
		log.Fatal("Error getting reviso user:", err)
	}

	token, err := jwt.GenerateUserToken(user)
	if err != nil {
		log.Fatal("Error generating JWT token:", err)
	}

	header := http.Header{}
	cookie := &http.Cookie{
		Name:  constants.CookieName,
		Value: token,
	}
	header.Add("Cookie", cookie.String())

	dialer := websocket.DefaultDialer

	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	log.Println("Connected to WebSocket!")
	at.conn = conn

	go at.listen()
	go at.edit()

	message := map[string]string{
		"type":  "subscribe",
		"docID": at.DocID,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
		log.Fatal("Error sending message:", err)
	}

	// wait for signal
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
}

func (at *AutoTyper) listen() {
	for {
		messageType, message, err := at.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break // Exit loop on error
		}

		switch messageType {
		case websocket.TextMessage:
			at.onMessage(message)
		default:
			log.Fatalf("Unexpected message type: %d", messageType)
		}
	}
}

func (at *AutoTyper) edit() {
	for {
		time.Sleep(editInterval)

		markdown, err := at.Doc.GetMarkdownBeforeAfter(v3.RootID, v3.LastID)
		if err != nil {
			log.Println("Error getting markdown:", err)
			continue
		}

		client := env.OpenAi(at.ctx)
		model := "gpt-4-turbo-2024-04-09"
		req := openai.ChatCompletionRequest{
			Model:       model,
			Temperature: 0.7,
			TopP:        0.8,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Repeat the markdown document that the user sent, but update the spelling and grammatical errors. Keep the exact formatting of the document.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: markdown,
				},
			},
		}

		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			log.Println("Error creating chat completion:", err)
			continue
		}

		updatedocument := resp.Choices[0].Message.Content

		mop, _, err := at.Doc.ApplyMarkdownDiff(at.Doc.Author, updatedocument, v3.RootID, v3.LastID)
		if err != nil {
			log.Println("Error applying markdown diff:", err)
			continue
		}

		opsJSON, err := json.Marshal(mop)
		if err != nil {

		}
		msg := map[string]interface{}{
			"type": "op",
			"op":   string(opsJSON),
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			continue
		}

		at.conn.WriteMessage(websocket.TextMessage, msgJSON)
		time.Sleep(typingDelay)
	}
}

func (at *AutoTyper) onMessage(message []byte) {
	// if the first byte is a { then it's an event
	if message[0] == '{' {
		var event map[string]interface{}
		err := json.Unmarshal(message, &event)
		if err != nil {
			log.Fatal("Error unmarshalling JSON:", err)
		}

		eventType, ok := event["type"]
		if !ok {
			log.Fatal("Error getting event type")
		}

		switch eventType {
		case "auth":
			at.authorID = event["authorID"].(string)
			at.Doc = v3.NewRogueForQuill(at.authorID)

			fmt.Printf("authorID: %s\n", at.authorID)
		default:
			fmt.Printf("Ignoring event: %s\n", eventType)
		}
		return
	}

	// otherwise it's a op
	msg := v3.Message{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

	_, err = at.Doc.MergeOp(msg.Op)
	if err != nil {
		log.Fatal("Error merging op:", err)
	}
}
