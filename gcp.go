package gcp

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

type Gcp struct {
	receivers map[string]Receiver
	senders   map[string]Sender
	Logger    *log.Logger
	logEntry  *log.Entry
}

func InitGcp(logger *log.Logger) (*Gcp, error) {
	ret := Gcp{}
	ret.receivers = map[string]Receiver{}
	ret.senders = map[string]Sender{}
	ret.Logger = logger
	ret.logEntry = logger.WithFields(log.Fields{
		"module": "gcp",
	})
	return &ret, nil
}

func (gcp *Gcp) Start() {
	gcp.startAllReceiver()

}

func (gcp *Gcp) AddReceiver(receiverName string, receiver Receiver) error {
	if _, ok := gcp.receivers[receiverName]; ok {
		return errors.New("Receiver Name has already exists.")
	}
	gcp.receivers[receiverName] = receiver
	return nil
}

func (gcp *Gcp) AddSender(senderName string, sender Sender) error {
	if _, ok := gcp.senders[senderName]; ok {
		return errors.New("Sender Name has already exists.")
	}
	gcp.senders[senderName] = sender
	return nil
}

func (gcp *Gcp) startAllReceiver() {
	for key := range gcp.receivers {
		gcp.logEntry.Info("Start Receiver:" + key)
		err := gcp.receivers[key].Start()
		if err != nil {
			gcp.logEntry.Errorf("Start Recevier %s %s. %s", key, " Fail", err)
			continue
		}
	}
}
