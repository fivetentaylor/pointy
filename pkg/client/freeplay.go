package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/teamreviso/freeplay"
	"github.com/tmc/langchaingo/llms"
)

type FreeplayInterface interface {
	Get(key string) (*freeplay.Prompt, error)
	All() ([]freeplay.Prompt, error)
	RecordCompletion(id string, prompt *freeplay.Prompt, messages []llms.MessageContent, response *llms.ContentResponse, start, end time.Time) error
	RecordTrace(id string, message *llms.MessageContent, response *llms.ContentResponse) error
}

type Freeplay struct {
	Client *freeplay.Client
	Logger *slog.Logger

	ctx       context.Context
	projectID string
	env       string
}

func NewFreeplayClientFromEnv() (*Freeplay, error) {
	url := os.Getenv("FREEPLAY_URL")
	projectID := os.Getenv("FREEPLAY_PROJECT_ID")
	env := os.Getenv("FREEPLAY_ENV")

	if url == "" || projectID == "" || env == "" {
		return nil, fmt.Errorf("FREEPLAY_URL, FREEPLAY_PROJECT_ID, and FREEPLAY_ENV must be set")
	}

	return NewFreeplayClient(url, projectID, env, 5*time.Minute)
}

func NewFreeplayClient(
	url string,
	projectID string,
	env string,
	duration time.Duration,
) (*Freeplay, error) {
	c, err := freeplay.NewClient(url)
	if err != nil {
		return nil, err
	}

	// logger := slog.Default()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	f := &Freeplay{
		Client: c,
		Logger: logger.WithGroup("freeplay"),

		projectID: projectID,
		env:       env,
		ctx:       context.Background(),
	}

	return f, nil
}

func (f *Freeplay) Get(key string) (*freeplay.Prompt, error) {
	prompt, err := f.get(key)
	if err != nil {
		f.Logger.Error("failed to get prompt", "key", key, "err", err)
		return nil, err
	}
	return prompt, nil
}

func (f *Freeplay) All() ([]freeplay.Prompt, error) {
	prompts, err := f.Client.GetAllPrompts(f.projectID)
	if err != nil {
		f.Logger.Error("failed to get prompts", "err", err)
		return nil, err
	}
	return prompts, nil
}

func (f *Freeplay) RecordCompletion(id string, provider, model, version string, messages []llms.MessageContent, data map[string]string, response *llms.ContentResponse, start, end time.Time) error {
	f.Logger.Info("[freeplay] recording completion", "response", response)

	msgs := []freeplay.Message{}
	for _, m := range messages {
		msgs = append(msgs, LlmMessageToFreeplayMessage(m))
	}

	for _, c := range response.Choices {
		msgs = append(msgs, freeplay.Message{
			Role:    string(llms.ChatMessageTypeAI),
			Content: c.Content,
		})
	}

	payload := freeplay.CompletionPayload{
		Messages: msgs,
		Inputs:   data,
		PromptInfo: freeplay.PromptInfo{
			PromptTemplateVersionID: version,
			Environment:             f.env,
		},
		CallInfo: &freeplay.CallInfo{
			StartTime:     float64(start.Unix()),
			EndTime:       float64(end.Unix()),
			Model:         model,
			Provider:      provider,
			ProviderInfo:  map[string]string{},
			LlmParameters: map[string]string{},
		},
	}

	f.Logger.Info("[freeplay] recording completion", "id", id, "payload", payload)

	comp, err := f.Client.RecordCompletion(
		f.projectID,
		id,
		&payload,
	)
	if err != nil {
		return err
	}

	f.Logger.Info("[freeplay] recorded completion", "completion", comp.CompletionID)

	return nil
}

func (f *Freeplay) RecordTrace(id string, message *llms.MessageContent, response *llms.ContentResponse) error {
	f.Logger.Info("[freeplay] recording trace", "msg", message, "response", response)
	if message == nil {
		f.Logger.Error("[freeplay] did not send trace, message is nil")
		return nil
	}
	inputMsg := LlmMessageToFreeplayMessage(*message)
	input := inputMsg.Content

	output := ""
	if response != nil {
		for _, c := range response.Choices {
			output += c.Content
		}
	}

	payload := freeplay.TracePayload{
		Input:  input,
		Output: output,
	}

	err := f.Client.RecordTrace(
		f.projectID,
		id,
		uuid.New().String(),
		&payload,
	)
	if err != nil {
		return err
	}

	f.Logger.Info("[freeplay] recorded trace")

	return nil
}

func (f *Freeplay) get(key string) (*freeplay.Prompt, error) {
	return f.Client.GetPrompt(
		f.projectID, key, false, f.env, map[string]string{},
	)
}

func (f *Freeplay) Close() {
	f.ctx.Done()
}

func LlmMessageToFreeplayMessage(m llms.MessageContent) freeplay.Message {
	content := strings.Builder{}
	for _, p := range m.Parts {
		switch v := p.(type) {
		case llms.TextContent:
			content.WriteString(v.Text)
		}
	}

	var roleStr string
	switch m.Role {
	case llms.ChatMessageTypeSystem:
		roleStr = "system"
	case llms.ChatMessageTypeHuman:
		roleStr = "user"
	case llms.ChatMessageTypeAI:
		roleStr = "assistant"
	default:
		roleStr = "user"
	}

	return freeplay.Message{
		Role:    roleStr,
		Content: content.String(),
	}
}

func strPrt(s string) *string {
	return &s
}
