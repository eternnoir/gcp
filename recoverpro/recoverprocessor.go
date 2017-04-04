package recoverpro

import (
	"context"
	"fmt"

	"github.com/eternnoir/gcp"
)

func New(gp *gcp.Gcp, np gcp.Processor) (*RecoverProcessor, error) {
	return &RecoverProcessor{gcp: gp, nextProcessor: np}, nil
}

type RecoverProcessor struct {
	gcp           *gcp.Gcp
	nextProcessor gcp.Processor
}

func (rp *RecoverProcessor) Process(ctx context.Context, payload interface{}) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	return rp.nextProcessor.Process(ctx, payload)
}
