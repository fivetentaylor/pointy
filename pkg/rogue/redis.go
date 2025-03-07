package rogue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/fivetentaylor/pointy/pkg/constants"
)

func SubscribeToDoc(ctx context.Context, client redis.UniversalClient, docID string) *redis.PubSub {
	return client.Subscribe(ctx, fmt.Sprintf(constants.DocUpdateChanFormat, docID))
}

func PublishToDoc(ctx context.Context, client redis.PubSubCmdable, docID string, op []byte) error {
	return client.Publish(context.Background(), fmt.Sprintf(constants.DocUpdateChanFormat, docID), op).Err()
}

func AddAuthorToActiveConnections(ctx context.Context, client redis.SetCmdable, docID, userID, authorID string) error {
	return client.SAdd(ctx,
		fmt.Sprintf(constants.DocActiveConnectionsKey, docID),
		fmt.Sprintf("%s:%s", userID, authorID),
	).Err()
}

func CurrentActiveConnections(ctx context.Context, client redis.SetCmdable, docID string) ([]string, error) {
	return client.SMembers(context.Background(), fmt.Sprintf(constants.DocActiveConnectionsKey, docID)).Result()
}

func RemoveAuthorFromActiveConnections(ctx context.Context, client redis.SetCmdable, docID, userID, authorID string) error {
	return client.SRem(ctx,
		fmt.Sprintf(constants.DocActiveConnectionsKey, docID),
		fmt.Sprintf("%s:%s", userID, authorID),
	).Err()
}

func AddAuthorLastCursor(ctx context.Context, client redis.StringCmdable, docID, userID, authorID string, cursor []byte) error {
	return client.Set(
		ctx,
		fmt.Sprintf(constants.DocUserConnectionKey, docID, userID, authorID),
		cursor,
		PresenceCheckInterval+1*time.Second,
	).Err()
}

func ExtendAuthorLastCursor(ctx context.Context, client redis.GenericCmdable, docID, userID, authorID string) error {
	return client.Expire(
		ctx,
		fmt.Sprintf(constants.DocUserConnectionKey, docID, userID, authorID),
		PresenceCheckInterval+1*time.Second,
	).Err()
}

func GetAuthorLastCursor(ctx context.Context, client redis.StringCmdable, docID, userID, authorID string) ([]byte, error) {
	return client.Get(
		ctx,
		fmt.Sprintf(constants.DocUserConnectionKey, docID, userID, authorID),
	).Bytes()
}

func RemoveAuthorLastCursor(ctx context.Context, client redis.GenericCmdable, docID, userID, authorID string) error {
	return client.Del(
		ctx,
		fmt.Sprintf(constants.DocUserConnectionKey, docID, userID, authorID),
	).Err()
}
