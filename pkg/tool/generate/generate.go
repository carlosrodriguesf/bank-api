//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package generate

import (
	"github.com/google/uuid"
	"time"
)

type (
	Generate interface {
		UUID() string
		CurrentTime() time.Time
	}
	generate struct {
	}
)

func New() Generate {
	return &generate{}
}

func (g generate) UUID() string {
	return uuid.NewString()
}

func (g generate) CurrentTime() time.Time {
	return time.Now()
}
