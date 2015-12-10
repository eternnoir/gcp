package lib

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gcp"
	"net"
)

type SocketReceiver struct {
	basegcp        *gcp.Gcp
	Type           string
	Host           string
	Port           string
	ProcressorList []gcp.Processor
	listener       net.Listener
	logEntry       *log.Entry
}

func InitSocketReceiver(host, port string, gcp *gcp.Gcp) (*SocketReceiver, error) {
	ret := SocketReceiver{}
	ret.basegcp = gcp
	ret.Host = host
	ret.Port = port
	ret.logEntry = gcp.Logger.WithFields(log.Fields{
		"module": "SocketReceiver",
	})
	return &ret, nil
}

func (sr SocketReceiver) Start() error {
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

func (sr SocketReceiver) handleIncome(conn net.Conn) {
	for _, proc := range sr.ProcressorList {
		go proc.Process(conn, sr.basegcp)
	}
}
