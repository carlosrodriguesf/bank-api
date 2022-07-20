package model

import (
	"regexp"
	"time"
)

var DocumentRegex = regexp.MustCompile("\\D")

type GeneratedData struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
