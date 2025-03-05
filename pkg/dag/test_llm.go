package dag

import (
	"context"

	"github.com/teamreviso/code/pkg/env"
	"github.com/tmc/langchaingo/llms"
)

type TestLLM struct {
	Response string
}

func (t *TestLLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	log := env.Log(ctx)
	opts := &llms.CallOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.StreamingFunc != nil {
		log.Info("TestLLM: streaming", "response", t.Response)
		opts.StreamingFunc(ctx, []byte(t.Response))
	}

	return &llms.ContentResponse{Choices: []*llms.ContentChoice{
		{Content: t.Response},
	}}, nil
}

func (t *TestLLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	return t.Response, nil
}
