package lib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gcp"
	"github.com/eternnoir/pool"
	"net"
	"reflect"
	"time"
)

var TimeOutError = errors.New("Send data timeout.")

type SocketSender struct {
	logEntry *log.Entry
	basegcp  *gcp.Gcp
	connpool pool.Pool
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
	factory := func() (net.Conn, error) { return net.Dial("tcp", host+":"+port) }
	p, err := pool.NewChannelPool(minconn, maxconn, factory)
	if err != nil {
		return nil, err
	}
	sender.connpool = p
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
	c := make(chan net.Conn, 1)
	go func() {
		conn, err := sender.connpool.Get(timeout / 2)
		if err != nil {
			sender.logEntry.Debugf("SocketSender Get connection error.%s", err)
			return
		}
		c <- conn
	}()

	defer func() {
		if len(c) > 0 {
			conn := <-c
			conn.Close()
		}
	}()
	select {
	case conn := <-c:
		return sender.fireRequest(data, conn)
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, TimeOutError
	}
}

func (sender *SocketSender) fireRequest(data []byte, conn net.Conn) ([]byte, error) {
	defer conn.Close()
	_, err := conn.Write(data)
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

func decodeLen(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}
