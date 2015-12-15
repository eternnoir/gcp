package gcp

type Processor interface {
	Process(context interface{}) (interface{}, error)
}
