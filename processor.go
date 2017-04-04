package gcp

import "context"

type Processor interface {
	Process(ctx context.Context, payload interface{}) (interface{}, error)
}
