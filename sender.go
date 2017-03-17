package gcp

import "context"

type Sender interface {
	Start() error
	Stop() error
	Send(ctx context.Context, paylodad interface{}, timeout int) (interface{}, error)
}
