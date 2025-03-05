package models

import "time"

type MessagingLimitType string

const (
	MessagingLimitTypeFree      MessagingLimitType = "FREE"
	MessagingLimitTypePremium   MessagingLimitType = "PREMIUM"
	MessagingLimitTypeEducation MessagingLimitType = "EDUCATION"
)

type MessagingLimit struct {
	Type       MessagingLimitType `json:"type"`
	Used       int                `json:"used"`
	Total      int                `json:"total"`
	StartingAt time.Time          `json:"starting_at"`
	EndingAt   time.Time          `json:"ending_at"`
}

func (m *MessagingLimit) Open() bool {
	return m.Used >= 0 && m.Used < m.Total
}
