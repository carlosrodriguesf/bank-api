package model

import "time"

type Transfer struct {
	ID              string    `json:"id" db:"id"`
	OriginAccountID string    `json:"origin_account_id" db:"origin_account_id"`
	TargetAccountID string    `json:"target_account_id" db:"target_account_id"`
	Amount          int64     `json:"amount" db:"amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
