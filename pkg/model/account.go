package model

import "time"

type Account struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name" validate:"required"`
	Document   string    `json:"document" db:"document" validate:"required"`
	Balance    int64     `json:"balance" db:"balance" validate:"required,min=1"`
	Secret     string    `json:"-" db:"secret" validate:"required"`
	SecretSalt string    `json:"-" db:"secret_salt"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
