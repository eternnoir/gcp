package gcp

type Processor interface {
	Process(context interface{}, basegcp *Gcp) error
}
