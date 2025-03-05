package ai

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"github.com/teamreviso/code/pkg/env"
)

func Stream(ctx context.Context, req openai.ChatCompletionRequest, onMessage func(m openai.ChatCompletionStreamResponse) error) (*openai.ChatCompletionStreamResponse, error) {
	openaiClient := env.OpenAi(ctx)

	stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error creating chat completion stream: %w", err)
	}
	defer stream.Close()

	response := &openai.ChatCompletionStreamResponse{}
	for {
		delta, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return response, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error receiving chat completion stream: %w", err)
		}

		err = mergeDeltas(delta, response)
		if err != nil {
			return nil, fmt.Errorf("error merging deltas: %w", err)
		}

		err = onMessage(*response)
		if err != nil {
			return nil, fmt.Errorf("error on message: %w", err)
		}
	}
}

func mergeDeltas(delta openai.ChatCompletionStreamResponse, response *openai.ChatCompletionStreamResponse) error {
	if delta.ID != "" {
		response.ID = delta.ID
	}

	if delta.Object != "" {
		response.Object = delta.Object
	}

	if delta.Created != 0 {
		response.Created = delta.Created
	}

	if len(delta.Choices) > 0 {
		if response.Choices == nil {
			response.Choices = []openai.ChatCompletionStreamChoice{}
		}
		if len(response.Choices) == 0 {
			response.Choices = append(response.Choices, delta.Choices[0])
		} else {
			cd := delta.Choices[0].Delta

			if cd.Content != "" {
				response.Choices[0].Delta.Content += cd.Content
			}

			if cd.Role == "" {
				response.Choices[0].Delta.Role = cd.Role
			}

			if len(cd.ToolCalls) > 0 {
				if response.Choices[0].Delta.ToolCalls == nil {
					response.Choices[0].Delta.ToolCalls = []openai.ToolCall{
						cd.ToolCalls[0],
					}
				}
				if len(response.Choices[0].Delta.ToolCalls) == 0 {
					response.Choices[0].Delta.ToolCalls[0] = cd.ToolCalls[0]
				} else {
					tc := cd.ToolCalls[0]

					if tc.ID != "" {
						response.Choices[0].Delta.ToolCalls[0].ID = tc.ID
					}

					if tc.Type != "" {
						response.Choices[0].Delta.ToolCalls[0].Type = tc.Type
					}

					if tc.Function.Name != "" {
						response.Choices[0].Delta.ToolCalls[0].Function.Name = tc.Function.Name
					}

					if tc.Function.Arguments != "" {
						response.Choices[0].Delta.ToolCalls[0].Function.Arguments += tc.Function.Arguments
					}
				}
			}
		}
	}

	return nil
}
