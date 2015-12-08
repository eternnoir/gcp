package gcp

type Sender interface {
	Start() error
	Stop() error
}
