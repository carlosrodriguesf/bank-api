//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package cache

import (
	"context"
	"time"
)

type Options struct {
	Url      string
	Password string
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, d time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	GetUpdating(ctx context.Context, key string, value interface{}, d time.Duration) error
	Close() error
	IsErrCacheMissing(err error) bool
}
