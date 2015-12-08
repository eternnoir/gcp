package gcp

import (
	"errors"
)

type Gcp struct {
	receivers map[string]*Receiver
	senders   map[string]*Sender
}

func InitGcp() (*Gcp, error) {
	ret := Gcp{}
	return &ret, nil
}

func (gcp Gcp) AddReceiver(receiverName string, receiver *Receiver) error {
	if _, ok := gcp.receivers[receiverName]; ok {
		return errors.New("Receiver Name has already exists.")
	}
	return nil
}

func (gcp Gcp) AddSender(senderName string, sender *Sender) error {
	if _, ok := gcp.senders[senderName]; ok {
		return errors.New("Sender Name has already exists.")
	}
	return nil
}
