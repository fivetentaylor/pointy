package pubsub

import (
	"context"
)

type PubSubInterface interface {
	Subscribe(docID string, onMsg func(msg []byte)) (error, context.CancelFunc)
	Publish(docID string, op []byte) error
}
