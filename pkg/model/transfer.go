package model

import "time"

type (
	Transfer struct {
		ID              string    `json:"id" db:"id"`
		OriginAccountID string    `json:"origin_account_id" db:"origin_account_id"`
		TargetAccountID string    `json:"target_account_id" db:"target_account_id"`
		Amount          int64     `json:"amount" db:"amount"`
		CreatedAt       time.Time `json:"created_at" db:"created_at"`
	}
	TransferDetailed struct {
		Transfer
		Sent              bool   `json:"sent" db:"sent"`
		OriginAccountName string `json:"origin_account_name" db:"origin_account_name"`
		TargetAccountName string `json:"target_account_name" db:"target_account_name"`
	}
)
