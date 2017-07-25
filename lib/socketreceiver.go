package lib

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/eternnoir/gcp"
	log "github.com/sirupsen/logrus"
)

type SocketReceiver struct {
	Name           string
	basegcp        *gcp.Gcp
	Type           string
	ListenAddr     string
	ProcressorList []gcp.Processor
	listener       net.Listener
	tlsconfig      *tls.Config
	logEntry       *log.Entry
}

func InitSocketReceiver(bgcp *gcp.Gcp, name, listenAddr string, tlscofnig *tls.Config) (gcp.Receiver, error) {
	ret := SocketReceiver{}
	ret.Name = name
	ret.basegcp = bgcp
	ret.ListenAddr = listenAddr
	ret.Type = "tcp"
	ret.ProcressorList = []gcp.Processor{}
	ret.logEntry = bgcp.Logger.WithFields(log.Fields{
		"module": ret.Name,
	})
	ret.tlsconfig = tlscofnig
	return &ret, nil
}

func (sr *SocketReceiver) Start() error {
	if sr.tlsconfig == nil {
		return sr.startSocket()
	} else {
		return sr.startTls()
	}
}

func (sr *SocketReceiver) startTls() error {
	var err error
	sr.listener, err = tls.Listen(sr.Type, sr.ListenAddr, sr.tlsconfig)
	sr.logEntry.Infof("Start tls %s server. %s", sr.Type, sr.ListenAddr)
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

}

func (sr *SocketReceiver) startSocket() error {
	var err error
	sr.listener, err = net.Listen(sr.Type, sr.ListenAddr)
	sr.logEntry.Infof("Start %s server. %s", sr.Type, sr.ListenAddr)
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
}

func (sr *SocketReceiver) handleIncome(conn net.Conn) {
	for _, proc := range sr.ProcressorList {
		ctx := context.Background()
		proc.Process(ctx, conn)
	}
}

func (sr *SocketReceiver) AddProcessor(processor gcp.Processor) []gcp.Processor {
	sr.ProcressorList = append(sr.ProcressorList, processor)
	return sr.ProcressorList
}
