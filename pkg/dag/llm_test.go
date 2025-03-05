package dag_test

import (
	"context"
	"strings"

	"github.com/tmc/langchaingo/llms"

	"github.com/teamreviso/code/pkg/dag"
	"github.com/teamreviso/code/pkg/stackerr"
)

type TestLLMAdapter struct {
	Chunks [][]byte
}

func NewTestLLMAdapterWithResponse(response string) *TestLLMAdapter {
	// breaking up the response into chunks
	chunks := make([][]byte, 0)
	for i := 0; i < len(response); i += 5 {
		e := min(len(response), i+5)
		chunks = append(chunks, []byte(response[i:e]))
	}

	return &TestLLMAdapter{
		Chunks: chunks,
	}
}

func (ta TestLLMAdapter) GenerateContentForStoredPrompt(ctx context.Context, _ string, _ string, _ []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	callOptions := llms.CallOptions{}
	for _, opt := range options {
		opt(&callOptions)
	}

	var content strings.Builder

	for _, chunk := range ta.Chunks {
		if callOptions.StreamingFunc != nil {
			err := callOptions.StreamingFunc(ctx, chunk)
			if err != nil {
				return nil, stackerr.Wrap(err)
			}

		}
		content.WriteString(string(chunk))
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{
				Content: content.String(),
			},
		},
	}, nil
}

func (ta TestLLMAdapter) GenerateStoredPrompt(ctx context.Context, _ string, _ string, _ map[string]string, options ...llms.CallOption) (_ *llms.ContentResponse, _ error) {
	return ta.GenerateContentForStoredPrompt(ctx, "", "", nil, options...)
}

func (TestLLMAdapter) GetModel() (_ string) {
	panic("not implemented GetModel") // TODO: Implement
}
func (TestLLMAdapter) GetProvider() (_ dag.LLMProvider) {
	panic("not implemented GetProvider") // TODO: Implement
}
func (TestLLMAdapter) GetLLM(_ context.Context) (_ llms.Model, _ error) {
	panic("not implemented GetLLM") // TODO: Implement
}
func (TestLLMAdapter) SetLLM(_ llms.Model) {
	panic("not implemented SetLLM") // TODO: Implement
}
func (TestLLMAdapter) GetLLMProvider(_ context.Context, _ string, _ string) (_ llms.Model, _ error) {
	panic("not implemented GetLLMProvider") // TODO: Implement
}
func (TestLLMAdapter) GenerateContent(_ context.Context, _ []llms.MessageContent, _ ...llms.CallOption) (_ *llms.ContentResponse, _ error) {
	panic("not implemented GenerateContent") // TODO: Implement
}
func (TestLLMAdapter) GenerateFromSinglePrompt(_ context.Context, _ string, _ ...llms.CallOption) (_ string, _ error) {
	panic("not implemented GenerateFromSinglePrompt") // TODO: Implement
}
