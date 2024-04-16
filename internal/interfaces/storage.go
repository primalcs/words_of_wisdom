package interfaces

import (
	"context"
	"time"
)

type Cache interface {
	InsertClientToken(context.Context, string, string, time.Duration) error
	GetClientToken(context.Context, string) (string, error)
	Delete(context.Context, string) error
	Close(context.Context) error
}
