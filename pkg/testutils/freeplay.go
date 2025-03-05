package testutils

import (
	"fmt"
	"time"

	"github.com/teamreviso/freeplay"
	"github.com/tmc/langchaingo/llms"
)

type MockFreeplay struct {
	Prompts map[string]*freeplay.Prompt
}

func NewMockFreeplay() *MockFreeplay {
	return &MockFreeplay{
		Prompts: make(map[string]*freeplay.Prompt),
	}
}

func (c *MockFreeplay) Get(key string) (*freeplay.Prompt, error) {
	prpt, ok := c.Prompts[key]
	if !ok {
		return nil, fmt.Errorf("prompt not found")
	}
	return prpt, nil
}

func (c *MockFreeplay) All() []*freeplay.Prompt {
	all := make([]*freeplay.Prompt, 0, len(c.Prompts))

	for _, prpt := range c.Prompts {
		all = append(all, prpt)
	}

	return all
}

func (c *MockFreeplay) RecordCompletion(id string, prompt *freeplay.Prompt, messages []llms.MessageContent, response *llms.ContentResponse, start, end time.Time) error {
	return nil
}

func (c *MockFreeplay) RecordTrace(id string, message *llms.MessageContent, response *llms.ContentResponse) error {
	return nil
}
