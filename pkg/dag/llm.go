package dag

import (
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/stackerr"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMProvider string

const (
	OpenAI       LLMProvider = "openai"
	Groq         LLMProvider = "groq"
	Anthropic    LLMProvider = "anthropic"
	TestProvider LLMProvider = "test"
	VertexAI     LLMProvider = "vertexai"
)

const (
	DefaultModel       = "claude-sonnet-4-5"
	DefaultClaudeModel = "claude-sonnet-4-5"
	DefaultOpenAIModel = "gpt-4o-2024-11-20"
)

type LLMAdapter interface {
	GetModel() string
	GetProvider() LLMProvider
	GetLLM(context.Context) (llms.Model, error)
	SetLLM(llms.Model)
	GetLLMProvider(context.Context, string, string) (llms.Model, error)

	GenerateContent(context.Context, []llms.MessageContent, ...llms.CallOption) (*llms.ContentResponse, error)
	GenerateFromSinglePrompt(context.Context, string, ...llms.CallOption) (string, error)
	GenerateContentForStoredPrompt(
		context.Context,
		string, string,
		[]llms.MessageContent,
		...llms.CallOption,
	) (*llms.ContentResponse, error)

	GenerateStoredPrompt(
		context.Context,
		string, string,
		map[string]string,
		...llms.CallOption,
	) (*llms.ContentResponse, error)
}

type DefaultLLMAdapter struct {
	Provider         LLMProvider
	Model            string
	AllowLLMOverride bool

	llm llms.Model
}

func (n *DefaultLLMAdapter) GetModel() string {
	if n.Model == "" {
		return "claude-3-5-sonnet-20240620"
	}
	return n.Model
}

func (n *DefaultLLMAdapter) GetProvider() LLMProvider {
	if n.Provider == "" {
		return Anthropic
	}
	return n.Provider
}

func (n *DefaultLLMAdapter) GetLLM(ctx context.Context) (llms.Model, error) {
	if n.llm != nil {
		return n.llm, nil
	}

	switch n.GetProvider() {
	case OpenAI:
		return n.getOpenAI(ctx, n.GetModel())
	case Groq:
		return n.getGroq(ctx, n.GetModel())
	case Anthropic:
		return n.getClaude(ctx, n.GetModel())
	case VertexAI:
		return n.getVertexAI(ctx, n.GetModel())
	default:
		return n.getOpenAI(ctx, n.GetModel())
	}
}

func (n *DefaultLLMAdapter) SetLLM(llm llms.Model) {
	n.llm = llm
}

func (n *DefaultLLMAdapter) GetLLMProvider(ctx context.Context, provider, model string) (llms.Model, error) {
	if n.llm != nil {
		return n.llm, nil
	}
	log := env.SLog(ctx)
	log.Info("getting llm provider", "provider", provider, "model", model)

	switch provider {
	case "openai":
		return n.getOpenAI(ctx, model)
	case "groq":
		return n.getGroq(ctx, model)
	case "anthropic":
		return n.getClaude(ctx, model)
	case "test":
		return n.getTest(ctx, model)
	case "vertex":
		return n.getVertexAI(ctx, model)
	default:
		return n.getOpenAI(ctx, model)
	}
}

func (n *DefaultLLMAdapter) getOpenAI(ctx context.Context, model string) (*openai.LLM, error) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		return nil, stackerr.New(fmt.Errorf("OPENAI_API_KEY is not set"))
	}

	return openai.New(
		openai.WithModel(model),
		openai.WithBaseURL("https://api.openai.com/v1"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithHTTPClient(n.GetHTTPClient(ctx)),
	)
}

func (n *DefaultLLMAdapter) getGroq(ctx context.Context, model string) (*openai.LLM, error) {
	if os.Getenv("GROQ_API_KEY") == "" {
		return nil, stackerr.New(fmt.Errorf("GROQ_API_KEY is not set"))
	}

	return openai.New(
		openai.WithModel(model),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
		openai.WithToken(os.Getenv("GROQ_API_KEY")),
		openai.WithHTTPClient(n.GetHTTPClient(ctx)),
	)
}

func (n *DefaultLLMAdapter) getClaude(ctx context.Context, model string) (*anthropic.LLM, error) {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		return nil, stackerr.New(fmt.Errorf("ANTHROPIC_API_KEY is not set"))
	}

	log := env.SLog(ctx)
	log.Info("[dag] creating Claude client", "model", model)

	llm, err := anthropic.New(
		anthropic.WithModel(model),
		anthropic.WithToken(os.Getenv("ANTHROPIC_API_KEY")),
		anthropic.WithHTTPClient(n.GetHTTPClient(ctx)),
	)
	if err != nil {
		log.Error("[dag] failed to create Claude client", "error", err, "model", model)
		return nil, stackerr.Wrap(err)
	}

	return llm, nil
}

func (n *DefaultLLMAdapter) getTest(ctx context.Context, model string) (*TestLLM, error) {
	log := env.SLog(ctx)

	log.Info("Using test LLM", "response (model)", model)

	return &TestLLM{
		Response: model,
	}, nil
}

func (n *DefaultLLMAdapter) getVertexAI(ctx context.Context, model string) (*vertex.Vertex, error) {
	creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if creds == "" {
		return nil, stackerr.New(fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set"))
	}

	// Decode base64 credentials
	decodedCreds, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		return nil, stackerr.New(fmt.Errorf("failed to decode credentials: %w", err))
	}

	/*
		credsFiles := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
		if credsFiles == "" {
			return nil, stackerr.New(fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS_JSON is not set"))
		}
	*/

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, stackerr.New(fmt.Errorf("GOOGLE_CLOUD_PROJECT is not set"))
	}

	location := os.Getenv("GOOGLE_CLOUD_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	return vertex.New(ctx,
		googleai.WithCloudProject(projectID),
		googleai.WithCloudLocation(location),
		googleai.WithCredentialsJSON(decodedCreds),
		// googleai.WithCredentialsFile(credsFiles),
		googleai.WithDefaultModel(model),
	)
}

func (n *DefaultLLMAdapter) GetHTTPClient(ctx context.Context) *HTTPLogger {
	id := GetRunID(ctx)
	return &HTTPLogger{
		id:     id,
		Client: http.DefaultClient,
	}
}

func (n *DefaultLLMAdapter) GenerateFromSinglePrompt(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	llm, err := n.GetLLM(ctx)
	if err != nil {
		return "", err
	}
	start := time.Now()

	saveLogFile(ctx, fmt.Sprintf("llm_prompt-%d.txt", start.Unix()), prompt)

	out, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt, options...)
	if err != nil {
		saveLogFile(ctx, fmt.Sprintf("llm_error-%d.json", start.Unix()), err.Error())
		return "", err
	}

	saveLogFile(ctx, fmt.Sprintf("llm_output-%d.txt", start.Unix()), out)
	return out, nil
}

func (n *DefaultLLMAdapter) GenerateContentForStoredPrompt(
	ctx context.Context,
	id, promptName string,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	log := env.SLog(ctx)
	tbl := env.Query(ctx).Prompt

	prompt, err := tbl.Where(tbl.PromptName.Eq(promptName)).First()
	if err != nil {
		log.Error("[dag] error getting prompt", "error", err)
		return nil, fmt.Errorf("error getting prompt %s: %w", promptName, err)
	}

	msgs := make([]llms.MessageContent, 0, len(messages)+1)
	if prompt.SystemContent != nil {
		msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeSystem, *prompt.SystemContent))
	}
	msgs = append(msgs, messages...)

	llm, err := n.GetLLMProvider(ctx, prompt.Provider, prompt.ModelName)
	if err != nil {
		log.Error("[dag] error getting llm for freeplay", "error", err)
		return nil, fmt.Errorf("error getting llm for freeplay: %w", err)
	}

	for _, msg := range messages {
		log.Warn("freeplay", "message", msg)
	}

	// Note: Anthropic models don't allow both temperature and top_p to be set
	if prompt.Temperature > 0 {
		options = append(options, llms.WithTemperature(prompt.Temperature))
	} else if prompt.TopP > 0 {
		options = append(options, llms.WithTopP(prompt.TopP))
	}

	if prompt.MaxTokens > 0 {
		options = append(options, llms.WithMaxTokens(int(prompt.MaxTokens)))
	}

	log.Info("[freeplay] running prompt", "id", id, "promptName", promptName)

	start := time.Now()
	response, err := llm.GenerateContent(ctx, msgs, options...)
	if err != nil {
		log.Error("[dag] error generating content for freeplay prompt", "error", err)
		return nil, fmt.Errorf("error generating content for freeplay prompt: %w", err)
	}
	end := time.Now()
	log.Info("[freeplay] recorded trace", "id", id, "promptName", promptName, "took", end.Sub(start))

	return response, nil
}

func (n *DefaultLLMAdapter) GenerateStoredPrompt(
	ctx context.Context,
	id, promptName string,
	data map[string]string,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	log := env.SLog(ctx)
	tbl := env.Query(ctx).Prompt
	start := time.Now()

	prompt, err := tbl.Where(tbl.PromptName.Eq(promptName)).First()
	if err != nil {
		log.Error("[dag] error getting prompt", "error", err)
		return nil, fmt.Errorf("error getting prompt %s: %w", promptName, err)
	}

	formattedMessages, err := prompt.RenderContent(data)
	if err != nil {
		log.Error("[dag] error rendering prompt", "error", err)
		return nil, fmt.Errorf("error rendering prompt: %w", err)
	}

	// Store formatted messages in a log file
	go func() {
		var builder strings.Builder
		for _, msg := range formattedMessages {
			builder.WriteString(fmt.Sprintf("[%s]:\n%s\n\n", msg.Role, html.UnescapeString(msg.Content)))
		}

		saveLogFile(ctx, fmt.Sprintf("llm_formatted-%d.txt", start.Unix()), builder.String())
	}()

	// Build llm messages
	msgs := make([]llms.MessageContent, 0, len(formattedMessages))
	for _, msg := range formattedMessages {
		switch msg.Role {
		case "system":
			msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeSystem, msg.Content))
		case "user":
			msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeHuman, msg.Content))
		case "assistant":
			msgs = append(msgs, llms.TextParts(llms.ChatMessageTypeAI, msg.Content))
		default:
			return nil, fmt.Errorf("unknown role %q", msg.Role)
		}
	}

	provider, modelName, err := n.determineProviderAndModel(ctx, prompt)
	if err != nil {
		return nil, err
	}

	log.Info("[dag] attempting to get LLM provider", "provider", provider, "model", modelName, "promptName", promptName)

	// Get the prompt's model provider
	llm, err := n.GetLLMProvider(ctx, provider, modelName)
	if err != nil {
		log.Error("[dag] error getting llm for freeplay", "error", err, "provider", provider, "model", modelName, "promptName", promptName)
		return nil, fmt.Errorf("error getting llm for freeplay (provider=%s, model=%s): %w", provider, modelName, err)
	}

	// Append the options configured in freeplay into the options for the LLM call
	// Note: Anthropic models don't allow both temperature and top_p to be set
	if prompt.Temperature > 0 {
		options = append(options, llms.WithTemperature(prompt.Temperature))
	} else if prompt.TopP > 0 {
		options = append(options, llms.WithTopP(prompt.TopP))
	}
	if prompt.MaxTokens > 0 {
		options = append(options, llms.WithMaxTokens(int(prompt.MaxTokens)))
	}

	log.Info("[dag] calling GenerateContent", "provider", provider, "model", modelName, "promptName", promptName, "numMessages", len(msgs))

	response, err := llm.GenerateContent(ctx, msgs, options...)
	if err != nil {
		log.Error("[dag] error generating content for stored prompt", "error", err, "provider", provider, "model", modelName, "promptName", promptName)
		return nil, fmt.Errorf("error generating content for stored prompt (provider=%s, model=%s, prompt=%s): %w", provider, modelName, promptName, err)
	}

	go func() {
		saveLogFile(ctx, fmt.Sprintf("llm_output-%d.json", start.Unix()), response)

		var builder strings.Builder
		for _, choice := range response.Choices {
			builder.WriteString(fmt.Sprintf("%s\n\n", html.UnescapeString(choice.Content)))
		}

		saveLogFile(ctx, fmt.Sprintf("llm_choice-%d.txt", start.Unix()), builder.String())
	}()

	return response, nil
}

func (n *DefaultLLMAdapter) determineProviderAndModel(ctx context.Context, prompt *models.Prompt) (string, string, error) {
	provider := prompt.Provider
	modelName := prompt.ModelName

	llmChoice, err := GetStateKey[models.LLM_CHOICE](ctx, "llmChoice")
	if err != nil {
		return "", "", fmt.Errorf("error getting llm choice: %w", err)
	}

	if n.AllowLLMOverride && llmChoice != models.LLM_CHOICE_LLM_CHOICE_UNSPECIFIED {
		switch llmChoice {
		case models.LLM_CHOICE_LLM_GPT4O:
			provider = string(OpenAI)
			modelName = DefaultOpenAIModel
		case models.LLM_CHOICE_LLM_CLAUDE:
			provider = string(Anthropic)
			modelName = DefaultClaudeModel
		}
	}

	return provider, modelName, nil
}

func (n *DefaultLLMAdapter) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	llm, err := n.GetLLM(ctx)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	saveLogFile(ctx, fmt.Sprintf("llm_msgs-%d.json", start.Unix()), messages)
	response, err := llm.GenerateContent(ctx, messages, options...)
	if err != nil {
		return nil, err
	}

	saveLogFile(ctx, fmt.Sprintf("llm_ouput-%d.json", start.Unix()), response)
	if response.Choices != nil && len(response.Choices) > 0 {
		saveLogFile(ctx, fmt.Sprintf("llm_ouput_text-%d.txt", start.Unix()), response.Choices[0].Content)
	}

	return response, err
}

func lastUserMessage(messages []llms.MessageContent) *llms.MessageContent {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == llms.ChatMessageTypeHuman {
			return &messages[i]
		}
	}
	return nil
}

func lastAiMessage(messages []llms.MessageContent) *llms.MessageContent {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == llms.ChatMessageTypeSystem {
			return &messages[i]
		}
	}
	return nil
}
