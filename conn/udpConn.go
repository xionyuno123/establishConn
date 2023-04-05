package conn

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

type UDPConn struct {
	conn       *net.UDPConn
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
	buf        []byte
}

func NewUDPConn() *UDPConn {
	return &UDPConn{}
}

func (c *UDPConn) Connect(addr string) error {
	udpAddr, err := ParseUDPAddr(addr)
	if err != nil {
		return err
	}
	c.remoteAddr = udpAddr
	return nil
}

func (c *UDPConn) Bind(addr string) error {
	udpAddr, err := ParseUDPAddr(addr)
	if err != nil {
		return err
	}

	c.localAddr = udpAddr
	conn, err := net.ListenUDP("udp", c.localAddr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.buf = make([]byte, MTU)
	return nil
}

func (c *UDPConn) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *UDPConn) SendMessage(msg Message, ms time.Duration) error {
	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	c.conn.SetWriteDeadline(time.Now().Add(ms))
	_, err = c.conn.WriteTo(data, c.remoteAddr)
	return err
}

func (c *UDPConn) RecvMessage(ms time.Duration) (Message, error) {
	c.conn.SetReadDeadline(time.Now().Add(ms))
	rlen, addr, err := c.conn.ReadFromUDP(c.buf)
	if err != nil {
		return Message{}, err
	}

	if addr.IP.String() != c.remoteAddr.IP.String() ||
		addr.Port != c.remoteAddr.Port {
		return Message{}, errors.New("not a message from remote endpoint")
	}

	msg := Message{}

	err = json.Unmarshal(c.buf[:rlen], &msg)
	if err != nil {
		return Message{}, err
	}

	return msg, nil
}

func ParseUDPAddr(addr string) (*net.UDPAddr, error) {
	s := strings.Split(addr, ":")
	if len(s) != 2 {
		return nil, errors.New("invalid addr")
	}

	var ip net.IP

	if len(s[0]) == 0 {
		ip = net.ParseIP("127.0.0.1")
	} else {
		ip = net.ParseIP(s[0])
	}

	if ip == nil {
		return nil, errors.New("invalid ip addr")
	}

	port, err := strconv.ParseUint(s[1], 10, 16)

	if err != nil {
		return nil, errors.New("invalid port")
	}

	return &net.UDPAddr{Port: int(port), IP: ip}, nil
}
