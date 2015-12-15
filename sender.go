package gcp

type Sender interface {
	Start() error
	Stop() error
	Send(context interface{}, timeout int) (interface{}, error)
}
