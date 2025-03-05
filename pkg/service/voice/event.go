package voice

import (
	"encoding/json"
	"fmt"
	"strings"
)

func DefaultSessionDetails() (SessionDetails, error) {
	var instructions strings.Builder
	err := instructionTmpl.Execute(&instructions, map[string]string{})
	if err != nil {
		return SessionDetails{}, err
	}

	return SessionDetails{
		Modalities:        []string{"text", "audio"},
		Instructions:      instructions.String(),
		Voice:             "alloy",
		InputAudioFormat:  "pcm16",
		OutputAudioFormat: "pcm16",
		InputAudioTranscription: InputAudioTranscription{
			Model: "whisper-1",
		},
		TurnDetection: TurnDetection{
			Type:              "server_vad",
			Threshold:         0.5,
			PrefixPaddingMs:   300,
			SilenceDurationMs: 500,
		},
		Tools:                   []ToolDefinition{},
		ToolChoice:              "auto",
		Temperature:             0.8,
		MaxResponseOutputTokens: "inf",
	}, nil
}

type Event struct {
	Raw []byte
	Base
}

type EventInterface interface {
	GetEventID() string
	GetType() string
}

type Base struct {
	EventID string `json:"event_id"`
	Type    string `json:"type"`
}

func (e Base) GetEventID() string {
	return e.EventID
}

func (e Base) GetType() string {
	return e.Type
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.EventID = aux.EventID
	e.Type = aux.Type
	e.Raw = data

	return nil
}

func (e *Event) Unwrap(p interface{}) error {
	return json.Unmarshal(e.Raw, p)
}

// EVENT TYPES

type InputBufferAppend struct {
	Audio string `json:"audio"`
	Base
}

type InputAudioBufferSpeechStarted struct {
	AudioStartMs int    `json:"audio_start_ms"`
	ItemID       string `json:"item_id"`
	Base
}

type ConversationItemInputAudioTranscriptCompleted struct {
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
	ItemID       string `json:"item_id"`
	Base
}

type Error struct {
	Error InnerError `json:"error"`
	Base
}

type InnerError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Param   string `json:"param"`
	EventID string `json:"event_id"`
}

type ResponseAudioDelta struct {
	ResponseID   string `json:"response_id"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`

	Base
}

type ConversationItemCreate struct {
	PreviousItemID string      `json:"previous_item_id"`
	Item           ItemDetails `json:"item"`

	Base
}

type ResponseAudioTranscriptDelta struct {
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ResponseID   string `json:"response_id"`

	Base
}

type ResponseAudioTranscriptDone struct {
	ContentIndex int    `json:"content_index"`
	Transcript   string `json:"transcript"`
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ResponseID   string `json:"response_id"`

	Base
}

type ResponseDone struct {
	Response ResponseDoneDetails `json:"response"`
	Base
}

type ResponseDoneDetails struct {
	Object        string        `json:"object"`
	ID            string        `json:"id"`
	Status        string        `json:"status"`
	StatusDetails StatusDetails `json:"status_details"`
	Output        []ItemDetails `json:"output"`
	Usage         Usage         `json:"usage"`
}

type Usage struct {
	TotalTokens       int `json:"total_tokens"`
	InputTokens       int `json:"input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	InputTokenDetails struct {
		TextTokens          int `json:"text_tokens"`
		AudioTokens         int `json:"audio_tokens"`
		CachedTokens        int `json:"cached_tokens"`
		CachedTokensDetails struct {
			TextTokens  int `json:"text_tokens"`
			AudioTokens int `json:"audio_tokens"`
		} `json:"cached_tokens_details"`
	} `json:"input_token_details"`
	OutputTokenDetails struct {
		TextTokens  int `json:"text_tokens"`
		AudioTokens int `json:"audio_tokens"`
	} `json:"output_token_details"`
}

type StatusDetails struct {
	Type   string     `json:"type"`
	Error  InnerError `json:"error"`
	Reason string     `json:"reason"`
}

type ResponseOutputItemAdded struct {
	OutputIndex int         `json:"output_index"`
	ResponseID  string      `json:"response_id"`
	Item        ItemDetails `json:"item"`
	Base
}

type ResponseOutputItemDone struct {
	ResponseID  string      `json:"response_id"`
	OutputIndex int         `json:"output_index"`
	Item        ItemDetails `json:"item"`
	Base
}

func GetArgument[T any](r ResponseOutputItemDone, key string) (T, error) {
	var args map[string]any
	err := json.Unmarshal([]byte(r.Item.Arguments), &args)
	if err != nil {
		return *new(T), err
	}

	val, ok := args[key]
	if !ok {
		return *new(T), fmt.Errorf("key %s not found in arguments", key)
	}

	tval, ok := val.(T)
	if !ok {
		return *new(T), fmt.Errorf("key %s is not of type %T", key, tval)
	}

	return tval, nil
}

type ItemDetails struct {
	ID        string        `json:"id,omitempty"`
	Object    string        `json:"object,omitempty"`
	Type      string        `json:"type,omitempty"`
	Status    string        `json:"status,omitempty"`
	Role      string        `json:"role,omitempty"`
	Content   []ItemContent `json:"content,omitempty"`
	CallID    string        `json:"call_id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Arguments string        `json:"arguments,omitempty"`
	Output    string        `json:"output,omitempty"`
}

type ItemContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type ResponseFunctionCallArgumentsDone struct {
	Arguments   string `json:"arguments"`
	ResponseID  string `json:"response_id"`
	ItemID      string `json:"item_id"`
	OutputIndex int    `json:"output_index"`
	CallID      string `json:"call_id"`
	Base
}

type SessionUpdate struct {
	Session SessionDetails `json:"session"`
	Base
}

type SessionDetails struct {
	Modalities              []string                `json:"modalities"`
	Instructions            string                  `json:"instructions"`
	Voice                   string                  `json:"voice"`
	InputAudioFormat        string                  `json:"input_audio_format"`
	OutputAudioFormat       string                  `json:"output_audio_format"`
	InputAudioTranscription InputAudioTranscription `json:"input_audio_transcription"`
	TurnDetection           TurnDetection           `json:"turn_detection"`
	Tools                   []ToolDefinition        `json:"tools"`
	ToolChoice              string                  `json:"tool_choice"`
	Temperature             float64                 `json:"temperature"`
	MaxResponseOutputTokens string                  `json:"max_response_output_tokens"`
}

type InputAudioTranscription struct {
	Model string `json:"model"`
}

type TurnDetection struct {
	Type              string  `json:"type"`
	Threshold         float64 `json:"threshold"`
	PrefixPaddingMs   int     `json:"prefix_padding_ms"`
	SilenceDurationMs int     `json:"silence_duration_ms"`
}

type ToolDefinition struct {
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters,omitempty"`
}

type ToolParameters struct {
	Type       string                  `json:"type"`
	Properties map[string]ToolProperty `json:"properties"`
	Required   []string                `json:"required"`
}

type ToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ResponseCreate struct {
	Response ResponseDetails `json:"response"`
	Base
}

type ResponseCreated struct {
	Response ResponseCreatedDetails `json:"response"`
	Base
}

type ResponseCreatedDetails struct {
	ID string `json:"id"`
	Metadata map[string]any `json:"metadata"`
}

type ResponseDetails struct {
	Modalities              []string                 `json:"modalities,omitempty"`
	Instructions            string                   `json:"instructions,omitempty"`
	Voice                   string                   `json:"voice,omitempty"`
	InputAudioFormat        string                   `json:"input_audio_format,omitempty"`
	OutputAudioFormat       string                   `json:"output_audio_format,omitempty"`
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	Tools                   []Tool                   `json:"tools,omitempty"`
	ToolChoice              string                   `json:"tool_choice,omitempty"`
	Temperature             float64                  `json:"temperature,omitempty"`
	MaxResponseOutputTokens int                      `json:"max_response_output_tokens,omitempty"`
	Metadata                map[string]any           `json:"metadata,omitempty"`
}
