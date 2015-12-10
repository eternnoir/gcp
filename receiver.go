package gcp

type Receiver interface {
	Start() error
	AddProcessor(processor *Processor) error
}
