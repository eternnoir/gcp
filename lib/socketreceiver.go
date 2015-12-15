package lib

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gcp"
	"net"
)

type SocketReceiver struct {
	Name           string
	basegcp        *gcp.Gcp
	Type           string
	Host           string
	Port           string
	ProcressorList []gcp.Processor
	listener       net.Listener
	logEntry       *log.Entry
}

func InitSocketReceiver(name, host, port string, bgcp *gcp.Gcp) (gcp.Receiver, error) {
	ret := SocketReceiver{}
	ret.Name = name
	ret.basegcp = bgcp
	ret.Host = host
	ret.Port = port
	ret.Type = "tcp"
	ret.ProcressorList = []gcp.Processor{}
	ret.logEntry = bgcp.Logger.WithFields(log.Fields{
		"module": ret.Name,
	})
	return &ret, nil
}

func (sr *SocketReceiver) Start() error {
	var err error
	sr.listener, err = net.Listen(sr.Type, sr.Host+":"+sr.Port)
	if err != nil {
		return err
	}
	defer sr.listener.Close()
	for {
		conn, err := sr.listener.Accept()
		if err != nil {
			sr.logEntry.Errorln(err)
			continue
		}
		go sr.handleIncome(conn)
	}
	return nil
}

func (sr *SocketReceiver) handleIncome(conn net.Conn) {
	for _, proc := range sr.ProcressorList {
		proc.Process(conn)
	}
}

func (sr *SocketReceiver) AddProcessor(processor gcp.Processor) []gcp.Processor {
	sr.ProcressorList = append(sr.ProcressorList, processor)
	return sr.ProcressorList
}
