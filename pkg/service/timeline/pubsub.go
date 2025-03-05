package timeline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/storage/dynamo"
)

// EventType represents the type of timeline event
type EventType string

const (
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
	EventTypeInsert EventType = "insert"
)

// getChannelKey returns the Redis channel key for a given event type and document ID
func getChannelKey(eventType EventType, docID string) (string, error) {
	switch eventType {
	case EventTypeUpdate:
		return fmt.Sprintf(constants.ChannelTimelineEventUpdateFormat, docID), nil
	case EventTypeDelete:
		return fmt.Sprintf(constants.ChannelTimelineEventDeleteFormat, docID), nil
	case EventTypeInsert:
		return fmt.Sprintf(constants.ChannelTimelineEventInsertFormat, docID), nil
	}
	// If no valid event type is provided, return an error
	return "", fmt.Errorf("invalid event type: %s", eventType)
}

func ListenForTimelineEvents(ctx context.Context, ch chan *dynamo.TimelineEvent, docID string, eventTypes ...EventType) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)

	channels := make([]string, 0, len(eventTypes))
	for _, eventType := range eventTypes {
		key, err := getChannelKey(eventType, docID)
		if err != nil {
			log.Error("error getting channel key", "error", err)
			continue
		}
		channels = append(channels, key)
	}

	pubsub := rc.Subscribe(ctx, channels...)
	incoming := pubsub.Channel()

	defer func() {
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		close(ch)
	}()

	for {
		log.Debug("waiting for timeline events")

		select {
		case rm := <-incoming:
			payload := rm.Payload
			event := &dynamo.TimelineEvent{}

			err := json.Unmarshal([]byte(payload), event)
			if err != nil {
				log.Errorf("error listener unable to unmarshalling timeline event %s: %s", payload, err)
				continue
			}

			log.Info("timeline event payload", "payload", payload)
			log.Debug("sending event", "payload", payload)
			ch <- event
		case <-ctx.Done():
			return
		}
	}
}

func PublishTimelineEvent(ctx context.Context, event *dynamo.TimelineEvent, eventType EventType) error {
	log := env.Log(ctx)
	redis := env.Redis(ctx)
	bts, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key, err := getChannelKey(eventType, event.DocID)
	if err != nil {
		log.Error("error getting channel key", "error", err)
		return err
	}

	log.Info("publishing timeline event", "key", key, "eventID", event.EventID)
	return redis.Publish(ctx, key, string(bts)).Err()
}
