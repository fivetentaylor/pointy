package messaging

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/env"
	"github.com/teamreviso/code/pkg/models"
)

var FreeTierMessageLimit = 10
var PremiumTierMessageLimit = 999

// How long we keep usage data for
var TwoMonths = time.Hour * 24 * 60

func MessageLimit(ctx context.Context, userID string) (*models.MessagingLimit, error) {
	sub := messageLimitSub(ctx, userID)
	key := messageLimitKey(userID, sub)

	val, err := env.Redis(ctx).Get(ctx, key).Int()

	// ignore not found errors
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}

	limit := &models.MessagingLimit{
		Used:       val,
		Total:      FreeTierMessageLimit,
		Type:       models.MessagingLimitTypeFree,
		StartingAt: sub.CurrentPeriodStart,
		EndingAt:   sub.CurrentPeriodEnd,
	}
	if sub.SubscriptionPlanID != "" {
		limit.Total = PremiumTierMessageLimit
		limit.Type = models.MessagingLimitTypePremium
	}

	return limit, err
}

func IncrMessageUsage(ctx context.Context, userID string) error {
	rdb := env.Redis(ctx)
	key := messageLimitKey(userID, messageLimitSub(ctx, userID))
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	err = rdb.Expire(ctx, key, TwoMonths).Err()
	if err != nil {
		return err
	}

	return nil
}

func ResetMessageLimit(ctx context.Context, userID string) error {
	key := messageLimitKey(userID, messageLimitSub(ctx, userID))

	err := env.Redis(ctx).Del(ctx, key).Err()
	if err != nil {
		return err
	}

	_, err = MessageLimit(ctx, userID)

	return err
}

// messageLimitKeyAndSub returns the key for the message limit for the given user and if the user has a subscription
func messageLimitSub(ctx context.Context, userID string) *models.UserSubscription {
	q := env.Query(ctx)

	subTbl := q.UserSubscription

	userSub, err := subTbl.Where(subTbl.UserID.Eq(userID)).First()
	if err != nil || userSub == nil {
		start, end := currentFreePeriod()
		return &models.UserSubscription{
			SubscriptionPlanID: "",
			CurrentPeriodStart: start,
			CurrentPeriodEnd:   end,
		}
	}

	return userSub
}

func messageLimitKey(userID string, userSub *models.UserSubscription) string {
	period := fmt.Sprintf(
		"%s_%s",
		userSub.CurrentPeriodStart.Format("2006-01-02"),
		userSub.CurrentPeriodEnd.Format("2006-01-02"),
	)

	return fmt.Sprintf(constants.MessagingUsageKeyFormat, userID, period)
}

func currentFreePeriod() (time.Time, time.Time) {
	now := time.Now()

	// maybe make this weekly?
	// offset := (int(now.Weekday()) + 6) % 7 // Adjust for Go's Weekday() (Sunday=0)
	// firstDay := now.AddDate(0, 0, -offset).Truncate(24 * time.Hour)
	//
	// // Calculate the end of the current week (Sunday)
	// lastDay := weekStart.AddDate(0, 0, 6)

	// Get the first day of the current month
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// Get the last day of the current month
	lastDay := firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}
