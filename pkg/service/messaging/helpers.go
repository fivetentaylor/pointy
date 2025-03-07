package messaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/graph/loaders"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/storage/dynamo"
)

func MessageAuthor(ctx context.Context, msg *dynamo.Message) (*models.User, error) {
	q := env.Query(ctx)

	return q.User.Where(q.User.ID.Eq(msg.AuthorID)).First()
}

func HydrateMentions(ctx context.Context, msg *dynamo.Message) (*dynamo.Message, error) {
	log := env.Log(ctx)
	log.Debug("hydrating mentions", "msg", msg)
	if msg.MentionedUserIds == nil {
		return msg, nil
	}

	_, err := env.UserClaim(ctx)
	if err != nil {
		log.Errorf("error getting current user: %s", err)
		return nil, fmt.Errorf("please login")
	}

	mentionedUsers, err := loaders.GetUsers(ctx, msg.MentionedUserIds)
	if err != nil {
		log.Errorf("error loading mentioned users: %s", err)
		return nil, fmt.Errorf("error loading mentioned users: %s", err)
	}

	//go through msg.Content, find all mentions in format @:user:userID@ and add user's dysplay name like: @:user:userID:displayName@
	log.Debug("hydrating mentions", "mentionedUsers", mentionedUsers)
	if len(mentionedUsers) > 0 {
		for _, user := range mentionedUsers {
			msg.Content = strings.ReplaceAll(msg.Content, fmt.Sprintf("@:user:%s@", user.ID), fmt.Sprintf("@:user:%s:%s@", user.ID, user.DisplayName))
		}
	}
	log.Info("updated content", "content", msg.Content)
	return msg, nil
}
