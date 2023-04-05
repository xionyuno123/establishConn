package loop

import (
	"establishConn/conn"
	"establishConn/log"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

type Loop struct {
	conn         conn.Conn
	mtx          sync.Mutex
	running      bool
	onSdpF       func(sdp webrtc.SessionDescription)
	onCandidateF func(candidate string)
}

func NewLoop(conn conn.Conn) *Loop {
	return &Loop{
		conn:    conn,
		running: false,
	}
}

func (c *Loop) OnSdp(f func(sdp webrtc.SessionDescription)) {
	c.onSdpF = f
}

func (c *Loop) OnCandidate(f func(candidate string)) {
	c.onCandidateF = f
}

func (c *Loop) Start() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.running = true

	for {
		if !c.running {
			break
		}

		msg, err := c.conn.RecvMessage(time.Microsecond * 10)

		if err != nil {
			continue
		}

		if msg.Typ > conn.SDPTypeStart && msg.Typ < conn.SDPTypeEnd {
			sdp := webrtc.SessionDescription{}
			sdp.SDP = msg.Msg
			sdp.Type = webrtc.SDPType(msg.Typ)

			if c.onSdpF != nil {
				c.onSdpF(sdp)
			}
		} else if msg.Typ == conn.Candidate {
			if c.onCandidateF != nil {
				c.onCandidateF(msg.Msg)
			}
		} else {
			log.DefaultLogger().Error("udpConn recvMessage", 1, "recv unkown type message")
		}
	}
}

func (c *Loop) Stop() {
	c.running = true
}
