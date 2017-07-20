package gcp

import (
	"context"
	"net"
	"time"
)

type Sender interface {
	Start() error
	Stop() error
	Send(ctx context.Context, paylodad interface{}, timeout time.Duration) (interface{}, error)
	GetConn(ctx context.Context) (net.Conn, error)
}
