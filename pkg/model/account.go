package model

import "time"

type Account struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name" validate:"required"`
	Document   string    `json:"document" db:"document" validate:"required"`
	Secret     string    `json:"-" db:"secret" validate:"required"`
	SecretSalt string    `json:"-" db:"secret_salt"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
