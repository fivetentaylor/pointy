package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

func ListenForMessages(ctx context.Context, ch chan *dynamo.Message, docID, channelID string) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)
	pubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.MsgUpsertChanFormat, channelID))
	incoming := pubsub.Channel()

	defer func() {
		log.Info("closing message listener", "channelID", channelID)
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		close(ch)
	}()

	for {
		log.Debug("waiting for messages")

		select {
		case rm := <-incoming: // Message from redis
			payload := rm.Payload
			msg := &dynamo.Message{}

			err := json.Unmarshal([]byte(payload), msg)
			if err != nil {
				log.Errorf("error listener unable to unmarshalling message %s: %s", payload, err)
				continue
			}

			log.Debug("sending message", "payload", payload)
			ch <- msg
		case <-ctx.Done():
			return
		}
	}
}

func PublishMessage(ctx context.Context, message *dynamo.Message) error {
	log := env.Log(ctx)
	redis := env.Redis(ctx)
	bts, err := json.Marshal(message)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.MsgUpsertChanFormat, message.ChannelID)
	log.Debug("publishing message", "messageID", message.MessageID, "key", key, "content", message.Content)
	return redis.Publish(ctx, key, string(bts)).Err()
}

func ListenForChannels(ctx context.Context, ch chan *dynamo.Channel, docID string, userID string) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)
	pubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.ChannelUpsertChanFormat, docID))
	incoming := pubsub.Channel()

	unreadPubsub := rc.Subscribe(ctx, fmt.Sprintf(constants.UnreadChannelUpdateChanFormat, docID, userID))
	unreadIncoming := unreadPubsub.Channel()

	defer func() {
		log.Info("closing channel listener", "docID", docID)
		unreadPubsub.Unsubscribe(ctx)
		unreadPubsub.Close()
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		close(ch)
	}()

	for {
		log.Debug("waiting for channel updates")

		select {
		case rm := <-incoming: // Message from redis
			payload := rm.Payload
			msg := &dynamo.Channel{}

			err := json.Unmarshal([]byte(payload), msg)
			if err != nil {
				log.Errorf("error listener unable to unmarshalling channel %s: %s", payload, err)
				continue
			}

			if utils.Contains(msg.UserIDs, userID) {
				log.Info("sending channel", "payload", payload)
				ch <- msg
			}
		case rm := <-unreadIncoming:
			payload := rm.Payload
			msg := &dynamo.Channel{}

			err := json.Unmarshal([]byte(payload), msg)
			if err != nil {
				log.Errorf("error listener unable to unmarshalling channel %s: %s", payload, err)
				continue
			}

			ch <- msg
		case <-ctx.Done():
			return
		}
	}
}

func PublishChannel(ctx context.Context, channel *dynamo.Channel) error {
	log := env.Log(ctx)
	redis := env.Redis(ctx)
	bts, err := json.Marshal(channel)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.ChannelUpsertChanFormat, channel.DocID)
	log.Info("publishing channel update", "channelID", channel.ChannelID, "key", key)
	return redis.Publish(ctx, key, string(bts)).Err()
}

func PublishChannelUnread(ctx context.Context, channel *dynamo.Channel, rr *dynamo.ReadReceipt) error {
	log := env.Log(ctx)
	redis := env.Redis(ctx)
	bts, err := json.Marshal(channel)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.UnreadChannelUpdateChanFormat, rr.DocID, rr.UserID)
	log.Info("publishing channel unread", "channelID", channel.ChannelID, "key", key)
	return redis.Publish(ctx, key, string(bts)).Err()
}

func PublishMessageUnread(ctx context.Context, msg *dynamo.Message, rr *dynamo.ReadReceipt) error {
	// if message is a reply, update the parent message
	if msg.ParentContainerID != nil && strings.HasPrefix(msg.ContainerID, dynamo.MsgPrefix) {
		redis := env.Redis(ctx)
		parentMsgID := strings.TrimPrefix(msg.ContainerID, dynamo.MsgPrefix)
		msg, err := env.Dynamo(ctx).GetMessage(*msg.ParentContainerID, parentMsgID)
		if err != nil {
			return err
		}
		bts, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		key := fmt.Sprintf(constants.UnreadMessageUpdateChanFormat, rr.DocID, rr.UserID)
		return redis.Publish(ctx, key, string(bts)).Err()
	}

	return nil
}

func ListenForThreads(ctx context.Context, ch chan *dynamo.Thread, docID, userID string) {
	log := env.Log(ctx)
	rc := env.Redis(ctx)
	key := fmt.Sprintf(constants.ThreadUpdateChanFormat, docID, userID)
	pubsub := rc.Subscribe(ctx, key)
	incoming := pubsub.Channel()

	defer func() {
		log.Info("closing thread listener", "docID", docID)
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		close(ch)
	}()

	log.Debug("👂 listening for threads", "key", key)

	for {
		log.Debug("waiting for threads")

		select {
		case rm := <-incoming: // Thread from redis
			payload := rm.Payload
			msg := &dynamo.Thread{}

			err := json.Unmarshal([]byte(payload), msg)
			if err != nil {
				log.Errorf("error listener unable to unmarshalling thread %s: %s", payload, err)
				continue
			}

			log.Debug("sending thread", "payload", payload)
			ch <- msg
		case <-ctx.Done():
			return
		}
	}
}

func PublishThread(ctx context.Context, thread *dynamo.Thread) error {
	log := env.Log(ctx)
	redis := env.Redis(ctx)
	bts, err := json.Marshal(thread)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constants.ThreadUpdateChanFormat, thread.DocID, thread.UserID)
	log.Debug("publishing thread", "threadID", thread.ThreadID, "key", key)
	return redis.Publish(ctx, key, string(bts)).Err()
}
