package main

import (
	"establishConn/conn"
	"establishConn/log"
	"establishConn/loop"
	"flag"
	"time"

	"github.com/pion/webrtc/v3"
)

func main() {
	clientAddr := flag.String("client", "127.0.0.1:59990", "")
	serverAddr := flag.String("server", "127.0.0.1:59991", "")
	connType := flag.String("connType", "udp", "")

	flag.Parse()

	typ := conn.UDP

	switch *connType {
	case "UDP":
		typ = conn.UDP
	case "ZMQ":
		typ = conn.ZMQ
	case "HTTP":
		typ = conn.HTTP
	default:
	}

	peerConn, err := webrtc.NewPeerConnection(webrtc.Configuration{})

	if err != nil {
		log.DefaultLogger().Error("NewPeerConnection", 1, err.Error())
		return
	}

	peerConn.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		log.DefaultLogger().Info("OnConnectionStateChange", 0, pcs.String())
	})

	peerConn.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.DefaultLogger().Info("OnDataChannel", 0, dc.Label())

		dc.OnOpen(func() {
			log.DefaultLogger().Info("OnOpen", 0, "")
		})

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.DefaultLogger().Info("OnMessage", 0, string(msg.Data))
			dc.SendText("hello world")
		})
	})

	c, err := conn.NewConn(typ)

	if err != nil {
		log.DefaultLogger().Error("NewConn", 1, err.Error())
		return
	}

	err = c.Bind(*serverAddr)
	if err != nil {
		log.DefaultLogger().Error("BindUDPConn", 1, err.Error())
		return
	}

	err = c.Connect(*clientAddr)

	if err != nil {
		log.DefaultLogger().Error("ConnectUDPConn", 1, err.Error())
		return
	}

	peerConn.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidate := i.ToJSON().Candidate
		err := c.SendMessage(conn.Message{
			Typ: conn.Candidate,
			Msg: candidate,
		}, time.Second)

		if err != nil {
			log.DefaultLogger().Error("SendMessage", 1, err.Error())
		}
	})

	loop := loop.NewLoop(c)

	loop.OnSdp(func(sdp webrtc.SessionDescription) {
		err := peerConn.SetRemoteDescription(sdp)
		if err != nil {
			log.DefaultLogger().Error("SetRemoteDescription", 1, err.Error())
			return
		}

		answer, err := peerConn.CreateAnswer(nil)

		if err != nil {
			log.DefaultLogger().Error("CreateAnswer", 1, err.Error())
			return
		}

		err = peerConn.SetLocalDescription(answer)

		if err != nil {
			log.DefaultLogger().Error("SetLocalDescription", 1, err.Error())
			return
		}

		err = c.SendMessage(conn.Message{
			Typ: conn.MsgType(answer.Type),
			Msg: answer.SDP,
		}, time.Second)

		if err != nil {
			log.DefaultLogger().Error("SendMessage", 1, err.Error())
			return
		}
	})

	loop.OnCandidate(func(candidate string) {
		c := webrtc.ICECandidateInit{Candidate: candidate}
		err := peerConn.AddICECandidate(c)
		if err != nil {
			log.DefaultLogger().Error("AddICECandidate", 1, err.Error())
		}
	})

	loop.Start()
}
