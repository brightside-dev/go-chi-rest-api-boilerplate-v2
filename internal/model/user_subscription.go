package model

import "time"

type UserSubscription struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	SubscriptionType string    `json:"subscription_type"`
	IsActive         bool      `json:"is_active"`
	StartDate        string    `json:"start_date"`
	EndDate          string    `json:"end_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
