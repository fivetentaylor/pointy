package dag

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

type BaseLLMNode struct {
	Provider         LLMProvider
	Model            string
	AllowLLMOverride bool
	Adapter          LLMAdapter

	llm llms.Model

	Base
}

func (n *BaseLLMNode) GetAdapter() LLMAdapter {
	if n.Adapter == nil {
		n.Adapter = &DefaultLLMAdapter{
			Provider:         n.Provider,
			Model:            n.Model,
			AllowLLMOverride: n.AllowLLMOverride,
		}
	}
	return n.Adapter
}

func (n *BaseLLMNode) SetAdapter(adapter LLMAdapter) {
	n.Adapter = adapter
}

func (n *BaseLLMNode) GetModel() string {
	return n.GetAdapter().GetModel()
}

func (n *BaseLLMNode) GetProvider() LLMProvider {
	return n.GetAdapter().GetProvider()
}

func (n *BaseLLMNode) GetLLM(ctx context.Context) (llms.Model, error) {
	return n.GetAdapter().GetLLM(ctx)
}

func (n *BaseLLMNode) SetLLM(llm llms.Model) {
	n.GetAdapter().SetLLM(llm)
}

func (n *BaseLLMNode) GetLLMProvider(ctx context.Context, provider, model string) (llms.Model, error) {
	return n.GetAdapter().GetLLMProvider(ctx, provider, model)
}

func (n *BaseLLMNode) GenerateFromSinglePrompt(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return n.GetAdapter().GenerateFromSinglePrompt(ctx, prompt, options...)
}

func (n *BaseLLMNode) GenerateContentForStoredPrompt(
	ctx context.Context,
	id, promptName string,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return n.GetAdapter().GenerateContentForStoredPrompt(ctx, id, promptName, messages, options...)
}

func (n *BaseLLMNode) GenerateStoredPrompt(
	ctx context.Context,
	id, promptName string,
	data map[string]string,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return n.GetAdapter().GenerateStoredPrompt(ctx, id, promptName, data, options...)
}

func (n *BaseLLMNode) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	return n.GetAdapter().GenerateContent(ctx, messages, options...)
}
