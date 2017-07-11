package gcp

import (
	"context"
	"time"
)

type Sender interface {
	Start() error
	Stop() error
	Send(ctx context.Context, paylodad interface{}, timeout time.Duration) (interface{}, error)
}
