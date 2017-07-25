package gcp

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Gcp struct {
	Receivers map[string]Receiver
	Senders   map[string]Sender
	Logger    *log.Entry
}

func InitGcp(logger *log.Entry) (*Gcp, error) {
	ret := Gcp{}
	ret.Receivers = map[string]Receiver{}
	ret.Senders = map[string]Sender{}
	ret.Logger = logger
	return &ret, nil
}

func (gcp *Gcp) Start() {
	gcp.startAllReceiver()
}

func (gcp *Gcp) AddReceiver(receiverName string, receiver Receiver) error {
	if _, ok := gcp.Receivers[receiverName]; ok {
		return errors.New("Receiver Name has already exists.")
	}
	gcp.Receivers[receiverName] = receiver
	return nil
}

func (gcp *Gcp) AddSender(senderName string, sender Sender) error {
	if _, ok := gcp.Senders[senderName]; ok {
		return errors.New("Sender Name has already exists.")
	}
	gcp.Senders[senderName] = sender
	return nil
}

func (gcp *Gcp) startAllReceiver() {
	for key := range gcp.Receivers {
		gcp.Logger.Info("Start Receiver:" + key)
		err := gcp.Receivers[key].Start()
		if err != nil {
			gcp.Logger.Errorf("Start Recevier %s %s. %s", key, " Fail", err)
			continue
		}
	}
}
