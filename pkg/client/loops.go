package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type LoopsInterface interface {
	SendEvent(ctx context.Context, eventName string, contactProperties LoopsContactProperties, eventProperties LoopsEventProperties) error
}

type Loops struct {
	Client  *http.Client
	Logger  *slog.Logger
	BaseURL string
	Token   string
}

type LoopsContactProperties struct {
	Email  string `json:"email"`
	UserId string `json:"userId"`
}

type LoopsEventProperties interface {
	json.Marshaler
}

func NewLoopsClientFromEnv() (*Loops, error) {
	token := os.Getenv("LOOPS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("LOOPS_TOKEN must be set")
	}

	return NewLoopsClient(token)
}

func NewLoopsClient(token string) (*Loops, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Loops{
		Client:  &http.Client{Timeout: 10 * time.Second},
		Logger:  logger.WithGroup("loops"),
		BaseURL: "https://app.loops.so/api/v1/events/send",
		Token:   token,
	}, nil
}

func (l *Loops) SendEvent(ctx context.Context, eventName string, contactProperties LoopsContactProperties, eventProperties LoopsEventProperties) error {
	l.Logger.Info("[loops] Sending event to loops", "eventName", eventName, "contactProperties", contactProperties, "eventProperties", eventProperties)

	eventPropsJSON, err := json.Marshal(eventProperties)
	if err != nil {
		return fmt.Errorf("failed to marshal event properties: %w", err)
	}

	payload := map[string]interface{}{
		"email":           contactProperties.Email,
		"userId":          contactProperties.UserId,
		"eventName":       eventName,
		"eventProperties": json.RawMessage(eventPropsJSON),
		"mailingLists":    map[string]interface{}{},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.BaseURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+l.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := l.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	l.Logger.Info("[loops] loops response", "response", string(body))

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
