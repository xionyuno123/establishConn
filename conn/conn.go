package conn

import (
	"errors"
	"time"
)

type Conn interface {
	Connect(addr string) error
	Bind(addr string) error
	Close()
	SendMessage(msg Message, ms time.Duration) error
	RecvMessage(ms time.Duration) (Message, error)
}

const MTU int = 1518

type MsgType int

const (
	SDPTypeStart MsgType = iota
	SDPTypeOffer
	SDPTypePranswer
	SDPTypeAnswer
	SDPTypeRollback
	SDPTypeEnd

	Candidate
)

type Message struct {
	Typ MsgType `json:"msgType"`
	Msg string  `json:"msg"`
}

type ConnType int

const (
	UDP ConnType = iota
	ZMQ
	HTTP
)

func NewConn(typ ConnType) (Conn, error) {
	switch typ {
	case UDP:
		return NewUDPConn(), nil
	case ZMQ:
		return nil, errors.New("not implemented")
	case HTTP:
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("unkown conn type")
	}
}
