package model

import (
	"time"
)

type GeneratedData struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
