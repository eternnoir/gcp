package gcp

import "context"

type Processor interface {
	Process(ctx context.Context, payloda interface{}) (interface{}, error)
}
