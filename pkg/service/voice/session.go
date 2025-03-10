package voice

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/service/messaging"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

const REALTIME_URL = "wss://api.openai.com/v1/realtime?model=gpt-4o-realtime-preview-2024-10-01"

var ignoreLogsForTypes = map[string]bool{
	"input_audio_buffer.append":              true,
	"response.audio_transcript.delta":        true,
	"response.audio.delta":                   true,
	"response.function_call_arguments.delta": true,
	"session.updated":                        true,
}

type Session struct {
	documentID string
	threadID   string
	authorID   string
	userID     string

	conn      *websocket.Conn
	connMutex sync.Mutex
	out       *websocket.Conn
	outMutex  sync.Mutex

	realtime       *rogue.Realtime
	docStore       *rogue.DocStore
	contentAddress *v3.ContentAddress
	lastCursor     *rogue.DocCursorOperation
	connected      bool

	messagesMap map[string]*dynamo.Message

	tools map[string]func(
		ctx context.Context,
		event ResponseOutputItemDone,
		session *Session,
	) error
	sessionDetails SessionDetails
}

func New(ctx context.Context, out *websocket.Conn, docID, threadID, userID, authorID string) (*Session, error) {
	s3 := env.S3(ctx)
	redis := env.Redis(ctx)
	query := env.Query(ctx)

	details, err := DefaultSessionDetails()
	if err != nil {
		return nil, err
	}

	s := &Session{
		authorID:       authorID,
		documentID:     docID,
		threadID:       threadID,
		userID:         userID,
		out:            out,
		messagesMap:    make(map[string]*dynamo.Message),
		sessionDetails: details,
		tools:          make(map[string]func(ctx context.Context, event ResponseOutputItemDone, session *Session) error),
	}

	s.realtime = rogue.NewRealtime(
		redis,
		query,
		s.documentID,
		s.userID,
		"",
		"",
	)

	s.docStore = rogue.NewDocStore(s3, query, redis)

	err = s.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Session) Connect(ctx context.Context) error {
	log := env.SLog(ctx)

	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer " + os.Getenv("OPENAI_API_KEY")}
	headers["OpenAI-Beta"] = []string{"realtime=v1"}

	var err error
	s.conn, _, err = websocket.DefaultDialer.Dial(REALTIME_URL, headers)
	if err != nil {
		log.Error("Failed to connect to server", "error", err)
		return err
	}

	log.Info("Connected to server", "url", REALTIME_URL, "connected", s.conn != nil)

	err = s.updateSessionDetails(ctx)
	if err != nil {
		return err
	}

	err = s.SendToUser(ctx, Connected{
		MessageBase: MessageBase{
			Type: "connected",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Listen(ctx context.Context) {
	log := env.SLog(ctx)

	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", "error", err)
			break
		}

		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			if strings.Contains(err.Error(), "closed network connection") {
				break
			}
			log.Error("Failed to unmarshal message", "error", err)
			continue
		}

		if _, ok := ignoreLogsForTypes[event.GetType()]; !ok {
			log.Info("🎙️ Received message", "type", event.Type, "id", event.EventID, "message", string(message))
		}
		switch event.Type {
		case "response.created":
			err := s.ResponseCreated(ctx, event)
			if err != nil {
				log.Error("Failed to append response created", "error", err)
				continue
			}
		case "response.audio.delta":
			err := s.AppendAudio(ctx, event)
			if err != nil {
				log.Error("Failed to append audio", "error", err)
				continue
			}
		case "response.audio_transcript.delta":
			err := s.AppendAudioTranscript(ctx, event)
			if err != nil {
				log.Error("Failed to append audio transcript", "error", err)
				continue
			}
		case "response.audio_transcript.done":
			err := s.AppendAudioTranscript(ctx, event)
			if err != nil {
				log.Error("Failed to append audio transcript", "error", err)
				continue
			}
		case "response.output_item.done":
			err := s.OutputItem(ctx, event)
			if err != nil {
				log.Error("Failed to append output item", "error", err)
				continue
			}
		case "input_audio_buffer.speech_started":
			err := s.SpeechStarted(ctx, event)
			if err != nil {
				log.Error("Failed to append input audio", "error", err)
				continue
			}
		case "conversation.item.input_audio_transcription.completed":
			err := s.AudioTranscriptComplete(ctx, event)
			if err != nil {
				log.Error("Failed to append input audio", "error", err)
				continue
			}
		case "response.done":
			err := s.ResponseDone(ctx, event)
			if err != nil {
				log.Error("Failed to append response done", "error", err)
				continue
			}
		case "error":
			err := Error{}
			err2 := event.Unwrap(&err)
			if err2 != nil {
				log.Error("Failed to unmarshal error message", "error", err2)
				continue
			}
			log.Error("Error from OpenAI", "error", err)
			if !s.connected {
				return
			}
		default:
			if _, ok := ignoreLogsForTypes[event.GetType()]; !ok {
				log.Info(fmt.Sprintf("🎧 Unhandled message: %s req_id=%s", event.Type, event.EventID))
			}
		}
	}
}

func (s *Session) CurrentContentAddress(ctx context.Context) (*v3.ContentAddress, error) {
	if s.contentAddress != nil {
		return s.contentAddress, nil
	}

	_, rog, err := s.docStore.GetCurrentDoc(ctx, s.documentID)
	if err != nil {
		return nil, err
	}

	ca, err := rog.GetFullAddress()
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (s *Session) CurrentMarshalledContentAddress(ctx context.Context) (string, error) {
	ca, err := s.CurrentContentAddress(ctx)
	if err != nil {
		return "", err
	}

	bts, err := json.Marshal(ca)
	if err != nil {
		return "", err
	}

	return string(bts), nil
}

func (s *Session) UpdateInstruction(ctx context.Context) error {
	log := env.SLog(ctx)
	_, rog, err := s.docStore.GetCurrentDoc(ctx, s.documentID)
	if err != nil {
		return err
	}

	ca, err := rog.GetFullAddress()
	if err != nil {
		log.Error("[voice.Session.UpdateInstruction] error getting content address", "error", err)
		return err
	}

	cursor, err := s.realtime.CursorForUser(ctx, s.userID, s.authorID)
	if err != nil {
		log.Error("[voice.Session.UpdateInstruction] error getting cursor", "error", err)
		return err
	}

	// If the content address hasn't changed, and the cursor hasn't changed, we don't need to update
	// the instruction
	// This allows OpenAI to cache the instruction, lowing cost by up to 80%
	if s.contentAddress == nil || reflect.DeepEqual(*s.contentAddress, *ca) {
		if (s.lastCursor == nil && cursor == nil) || (s.lastCursor != nil && reflect.DeepEqual(s.lastCursor.Range, cursor.Range)) {
			log.Info("[voice.Session.UpdateInstruction] content address unchanged")
			return nil
		}
	}

	s.contentAddress = ca
	s.lastCursor = cursor

	mkdoc, err := rog.GetFullMarkdown()
	if err != nil {
		log.Error("[voice.Session.UpdateInstruction] error getting markdown", "error", err)
		return err
	}

	input := map[string]string{"CurrentDocument": mkdoc}

	if cursor != nil {
		var startID *v3.ID
		var endID *v3.ID

		for i, id := range cursor.Range {
			if i == 0 {
				startID = &id
			}

			endID = &id
		}

		if startID != nil && endID != nil && startID != endID {
			selection, err := rog.GetMarkdown(*startID, *endID)
			if err == nil {
				input["CurrentSelection"] = selection
			}
		}
	}

	var instructions strings.Builder
	err = instructionTmpl.Execute(&instructions, input)
	if err != nil {
		log.Error("[voice.Session.UpdateInstruction] error executing template", "error", err)
		return err
	}

	log.Info("📚 Updated instruction", "instruction", instructions.String())

	s.sessionDetails.Instructions = instructions.String()
	return s.updateSessionDetails(ctx)
}

func (s *Session) SpeechStarted(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data InputAudioBufferSpeechStarted
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	err = s.SendToUser(ctx, SpeakingStarted{
		MessageBase: MessageBase{
			Type: "speaking_started",
		},
		IsSpeaking: true,
	})
	if err != nil {
		return err
	}

	log.Info("Speech started", "item", data.ItemID)
	
	_, err = s.GetRequestMessage(ctx, data.ItemID, nil)
	if err != nil {
		return err
	}

	return s.UpdateInstruction(ctx)
}

func (s *Session) AudioTranscriptComplete(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ConversationItemInputAudioTranscriptCompleted
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	log.Info("Speech complete", "item", data.ItemID, "transcript", data.Transcript)

	msg, err := s.GetRequestMessage(ctx, data.ItemID, nil)
	if err != nil {
		return err
	}

	msg.Content = data.Transcript
	msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) OutputItem(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseOutputItemDone
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	log.Info("Output item done", "data", data)

	if data.Item.Type == "function_call" {
		fn, ok := s.tools[data.Item.Name]
		if ok {
			go func() {
				err := fn(ctx, data, s)
				if err != nil {
					log.Error("[voice.Session.OutputItem] Error running function", "error", err)
				}
			}()

			return nil
		}

		return fmt.Errorf("function %s not found", data.Item.Name)
	}

	return nil
}

func (s *Session) FunctionCall(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseFunctionCallArgumentsDone
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	log.Info("Function call done", "data", data)

	return nil
}

func (s *Session) AppendAudio(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseAudioDelta
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	err = s.SendToUser(ctx, ResponseAudio{
		Delta:  data.Delta,
		ItemID: data.ItemID,
		MessageBase: MessageBase{
			Type: "response.audio.delta",
		},
	})
	if err != nil {
		log.Error("Failed to send audio buffer", "error", err)
		return err
	}

	return nil
}

func (s *Session) ResponseCreated(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseCreated
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
		return err
	}

	responseID := data.Response.ID
	if outputKey, ok := data.Response.Metadata["outputKey"]; ok {
		outputKeyStr, isString := outputKey.(string)
		if !isString {
			log.Error("outputKey is not a string", "outputKey", outputKey)
			return nil
		}

		existingMsg, exists := s.messagesMap[outputKeyStr]
		if exists {
			s.messagesMap[responseID] = existingMsg
			// we've mapped the message we can remove the key now
			delete(s.messagesMap, outputKeyStr)
		} else {
			log.Debug("No existing message found for outputKey", "outputKey", outputKeyStr)
		}
	}

	return nil
}

func (s *Session) AppendAudioTranscript(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseAudioTranscriptDelta
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}
	
	msg, err := s.GetResponseMessage(ctx, data.ResponseID, nil)
	if err != nil {
		log.Error("error getting response message", "error", err)
		return err
	}

	msg.Content = fmt.Sprintf("%s%s", msg.Content, data.Delta)
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		log.Error("error updating message", "error", err)
		return err
	}

	return nil
}

func (s *Session) CompleteAudioTranscript(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseAudioTranscriptDone
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	err = s.CloseResponseMessage(ctx, data.ResponseID)

	return nil
}

func (s *Session) ResponseDone(ctx context.Context, event Event) error {
	log := env.SLog(ctx)

	var data ResponseDone
	err := event.Unwrap(&data)
	if err != nil {
		log.Error("Failed to unmarshal message", "error", err)
	}

	if data.Response.Status == "failed" {
		log.Error("OpenAI response failed", "response", data.Response)

		s.SendToUser(ctx, Failure{
			MessageBase: MessageBase{
				Type: "failure",
			},
			Reason: "Upstream provider failed",
		})
	}

	log.Info("Response done", "data", data)
	return nil
}

func (s *Session) GetRequestMessage(ctx context.Context, key string, userMsg *dynamo.Message) (*dynamo.Message, error) {
	log := env.SLog(ctx)

	msg, ok := s.messagesMap[key]
	if ok {
		return msg, nil
	}

	adr, err := s.CurrentMarshalledContentAddress(ctx)
	if err != nil {
		log.Error("[voice] error getting current content address", "error", err)
		return nil, err
	}

	if userMsg == nil {
		userMsg = &dynamo.Message{
			Content:        "",
			LifecycleStage: dynamo.MessageLifecycleStagePending,
		}
	}

	userMsg.ContainerID = fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, s.threadID)
	userMsg.ChannelID = s.threadID
	userMsg.DocID = s.documentID
	userMsg.AuthorID = s.authorID
	userMsg.UserID = s.userID
	if userMsg.LifecycleStage == dynamo.MessageLifecycleStageUnknown {
		userMsg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	}

	if userMsg.MessageMetadata == nil {
		userMsg.MessageMetadata = &models.MessageMetadata{}
	}
	userMsg.MessageMetadata.AllowDraftEdits = true
	userMsg.MessageMetadata.ContentAddress = adr

	err = messaging.CreateMessage(ctx, s.documentID, userMsg)
	if err != nil {
		log.Error("error creating message", "error", err)
		return nil, err
	}

	s.messagesMap[key] = userMsg

	voiceMsg := NewMessage{
		MessageBase: MessageBase{
			Type: "new_message",
		},
	}
	err = s.SendToUser(ctx, voiceMsg)
	if err != nil {
		log.Error("error sending voice message", "error", err)
	}

	return userMsg, nil
}

func (s *Session) GetResponseMessage(ctx context.Context, key string, aiMsg *dynamo.Message) (*dynamo.Message, error) {
	log := env.SLog(ctx)

	msg, ok := s.messagesMap[key]
	if ok {
		return msg, nil
	}

	adr, err := s.CurrentMarshalledContentAddress(ctx)
	if err != nil {
		log.Error("[voice] error getting current content address", "error", err)
		return nil, err
	}

	if aiMsg == nil {
		aiMsg = &dynamo.Message{
			Content:        "",
			LifecycleStage: dynamo.MessageLifecycleStagePending,
		}
	}
	aiMsg.ContainerID = fmt.Sprintf("%s%s", dynamo.AiThreadPrefix, s.threadID)
	aiMsg.ChannelID = s.threadID
	aiMsg.DocID = s.documentID
	aiMsg.AuthorID = fmt.Sprintf("!%s", s.authorID)
	aiMsg.UserID = constants.RevisoUserID
	if aiMsg.LifecycleStage == dynamo.MessageLifecycleStageUnknown {
		aiMsg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	}

	if aiMsg.MessageMetadata == nil {
		aiMsg.MessageMetadata = &models.MessageMetadata{}
	}
	aiMsg.MessageMetadata.AllowDraftEdits = true
	aiMsg.MessageMetadata.ContentAddress = adr

	err = messaging.CreateMessage(ctx, s.documentID, aiMsg)
	if err != nil {
		log.Error("error creating message", "error", err)
		return nil, err
	}

	s.messagesMap[key] = aiMsg

	voiceMsg := NewMessage{
		MessageBase: MessageBase{
			Type: "new_message",
		},
	}
	err = s.SendToUser(ctx, voiceMsg)
	if err != nil {
		log.Error("error sending voice message", "error", err)
	}

	return aiMsg, nil
}

func (s *Session) CloseResponseMessage(ctx context.Context, key string) error {
	log := env.SLog(ctx)

	msg, ok := s.messagesMap[key]
	if !ok {
		return fmt.Errorf("[CloseResponseMessage] message not found")
	}

	// refresh message
	var err error
	msg, err = env.Dynamo(ctx).GetAiThreadMessage(s.threadID, msg.MessageID)
	if err != nil {
		log.Error("error getting message", "error", err)
		return err
	}

	msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
	err = messaging.UpdateMessage(ctx, msg)
	if err != nil {
		log.Error("error updating message", "error", err)
		return err
	}

	delete(s.messagesMap, key)

	return nil
}

func (s *Session) RemoveMessage(ctx context.Context, key string) {
	delete(s.messagesMap, key)
}

func (s *Session) Close(ctx context.Context) {
	log := env.SLog(ctx)

	err := s.conn.Close()
	if err != nil {
		log.Error("[realtime.Session] error closing websocket", "error", err)
	}

	for _, msg := range s.messagesMap {
		msg.LifecycleStage = dynamo.MessageLifecycleStageCompleted
		err = messaging.UpdateMessage(ctx, msg)
		if err != nil {
			log.Error("[realtime.Session] error updating message", "error", err)
		}
	}
}

func (s *Session) AddTool(ctx context.Context, tool Tool) error {
	definition := tool.Definition
	s.tools[definition.Name] = tool.Function
	s.sessionDetails.Tools = append(s.sessionDetails.Tools, definition)

	return s.updateSessionDetails(ctx)
}

func (s *Session) AppendInputAudio(ctx context.Context, buffer []byte) error {
	encodedAudio := ArrayBufferToBase64(buffer)

	event := InputBufferAppend{
		Audio: encodedAudio,
		Base: Base{
			EventID: GenerateID("evt_", 21),
			Type:    "input_audio_buffer.append",
		},
	}

	return s.SendToRealtime(ctx, event)
}

func (s *Session) updateSessionDetails(ctx context.Context) error {
	event := SessionUpdate{
		Session: s.sessionDetails,
		Base: Base{
			EventID: GenerateID("evt_", 21),
			Type:    "session.update",
		},
	}

	return s.SendToRealtime(ctx, event)
}

// addThreadContext is and attempt to add the context of the thread to the session
// Unforunately, if you add `input_text` types the model starts to respond back in
// text. If you turn off the text modality, it starts to respond back in audio, but
// then the audio transcriptions of the user don't come back...
//
// Alternatively, we could add the conversation to the instructions, but that would
// increase cost and honestly i don't know how important this is.
func (s *Session) addThreadContext(ctx context.Context) error {
	messages, err := env.Dynamo(ctx).GetMessagesForThread(s.threadID)
	if err != nil {
		return err
	}

	for _, message := range messages {
		role := "assistant"
		inputType := "text"
		if message.AuthorID == s.authorID {
			role = "user"
			inputType = "input_text"
		}

		event := ConversationItemCreate{
			Item: ItemDetails{
				ID:   GenerateID("msg", 21),
				Type: "message",
				Role: role,
				Content: []ItemContent{
					{
						Type: inputType,
						Text: message.Content,
					},
				},
			},
			Base: Base{
				EventID: GenerateID("evt_", 21),
				Type:    "conversation.item.create",
			},
		}

		err = s.SendToRealtime(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) SendToRealtime(ctx context.Context, event EventInterface) error {
	log := env.SLog(ctx)
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	if _, ok := ignoreLogsForTypes[event.GetType()]; !ok {
		log.Info("🎧 Sending realtime event", "type", event.GetType(), "id", event.GetEventID(), "event", event)
	}

	if s.conn == nil {
		log.Error("Connection is nil")
		return fmt.Errorf("connection is nil")
	}

	err := s.conn.WriteJSON(event)
	if err != nil {
		log.Error("Failed to send realtime event", "error", err)
		return err
	}

	return err
}

func (s *Session) SendToUser(ctx context.Context, message MessageInterface) error {
	log := env.SLog(ctx)
	s.outMutex.Lock()
	defer s.outMutex.Unlock()

	if _, ok := ignoreLogsForTypes[message.GetType()]; !ok {
		log.Info("👤 Sending message to user", "message", message)
	}

	if s.out == nil {
		log.Error("Output connection is nil")
		return fmt.Errorf("output connection is nil")
	}

	err := s.out.WriteJSON(message)
	if err != nil {
		log.Error("Failed to send user message", "error", err)
		return err
	}

	return err
}

func GenerateID(prefix string, length int) string {
	if length == 0 {
		length = 21
	}
	const chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length-len(prefix); i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}

	return prefix + sb.String()
}

func ArrayBufferToBase64(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

func Base64ToBuffer(buffer string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(buffer)
}
