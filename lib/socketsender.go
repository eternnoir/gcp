package lib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gcp"
	pool "github.com/eternnoir/gncp"
	"net"
	"reflect"
	"time"
)

var TimeOutError = errors.New("Send data timeout.")

type SocketSender struct {
	logEntry *log.Entry
	basegcp  *gcp.Gcp
	connpool pool.ConnPool
	Host     string
	Port     string
	MinConn  int
	MaxConn  int
}

func InitSocketSender(bgcp *gcp.Gcp, host, port string, minconn, maxconn int) (*SocketSender, error) {
	sender := new(SocketSender)
	sender.Host = host
	sender.Port = port
	sender.MaxConn = maxconn
	sender.MinConn = minconn
	sender.basegcp = bgcp
	sender.logEntry = bgcp.Logger.WithFields(log.Fields{
		"module": "SocketSender",
	})
	err := sender.initConnectionPool()
	if err != nil {
		return nil, err
	}
	return sender, nil
}

func (sender *SocketSender) Start() error {
	return nil
}

func (sender *SocketSender) Stop() error {
	return nil
}

func (sender *SocketSender) Send(context interface{}, timeout int) (interface{}, error) {
	binaryary, ok := context.([]byte)
	if !ok {
		return nil, errors.New(fmt.Sprintf("SocketSender cast process context error. It should be []byte but %s", reflect.TypeOf(context)))
	}
	return sender.processSendRequest(binaryary, timeout)
}

func (sender *SocketSender) processSendRequest(data []byte, timeout int) ([]byte, error) {
	result := make(chan []byte, 1)
	errc := make(chan error, 1)
	go func() {
		resultba, err := sender.fireRequest(data, timeout)
		if err != nil {
			sender.logEntry.Debugf("SocketSender Get connection error.%s", err)
			errc <- err
			return
		}
		result <- resultba
	}()
	select {
	case err := <-errc:
		return nil, err
	case res := <-result:
		return res, nil
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, TimeOutError
	}
}

func (sender *SocketSender) getConn(timeout int) (net.Conn, error) {
	return sender.connpool.GetWithTimeout(time.Duration(timeout) * time.Millisecond)
}

func (sender *SocketSender) fireRequest(data []byte, timeout int) ([]byte, error) {
	if sender.connpool == nil {
		return nil, errors.New("Connection Pool error.")
	}
	conn, err := sender.getConn(timeout / 2)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		sender.logEntry.Errorf("SocketSender Send data error. %s", err)
		return nil, err
	}

	lenbuffer := make([]byte, 2)
	_, err = conn.Read(lenbuffer)

	if err != nil {
		sender.logEntry.Error(conn.RemoteAddr().String(), " connection error: ", err)
		return nil, err
	}

	bytelen := decodeLen(lenbuffer)

	buffer := make([]byte, bytelen)
	_, err2 := conn.Read(buffer)
	sender.logEntry.Debugf("Socket Sender Receive data : %s", hex.EncodeToString(buffer))

	if err2 != nil {
		sender.logEntry.Error(conn.RemoteAddr().String(), " connection error: ", err)
		return nil, err
	}
	return buffer, nil
}

func (sender *SocketSender) initConnectionPool() error {
	factory := func() (net.Conn, error) { return net.Dial("tcp", sender.Host+":"+sender.Port) }
	p, err := pool.NewPool(sender.MinConn, sender.MaxConn, factory)
	if err != nil {
		return err
	}
	sender.connpool = p
	return nil
}

func decodeLen(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}
