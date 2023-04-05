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

	dataChannel, err := peerConn.CreateDataChannel("data", nil)

	if err != nil {
		log.DefaultLogger().Error("CreateDataChannel", 1, err.Error())
		return
	}

	dataChannel.OnOpen(func() {
		log.DefaultLogger().Info("OnOpen", 0, "")
		dataChannel.SendText("hello world")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.DefaultLogger().Info("OnMessage", 0, string(msg.Data))
		time.Sleep(time.Second * 2)
		dataChannel.SendText("hello world")
	})

	peerConn.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		log.DefaultLogger().Info("OnConnectionStateChange", 0, pcs.String())
	})

	peerConn.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.DefaultLogger().Info("OnDataChannel", 0, dc.Label())
	})

	c, err := conn.NewConn(typ)
	if err != nil {
		log.DefaultLogger().Error("NewConn", 1, err.Error())
		return
	}

	err = c.Bind(*clientAddr)

	if err != nil {
		log.DefaultLogger().Error("Bind", 1, err.Error())
		return
	}

	err = c.Connect(*serverAddr)

	if err != nil {
		log.DefaultLogger().Error("Connect", 1, err.Error())
		return
	}

	loop := loop.NewLoop(c)

	loop.OnCandidate(func(candidate string) {
		err := peerConn.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})
		if err != nil {
			log.DefaultLogger().Error("AddICECandidate", 1, err.Error())
		}
	})

	loop.OnSdp(func(sdp webrtc.SessionDescription) {
		err := peerConn.SetRemoteDescription(sdp)
		if err != nil {
			log.DefaultLogger().Error("SetRemoteDescription", 1, err.Error())
		}
	})

	peerConn.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidate := i.ToJSON().Candidate
		err := c.SendMessage(conn.Message{Typ: conn.Candidate, Msg: candidate}, time.Second)
		if err != nil {
			log.DefaultLogger().Error("SendMessage", 1, err.Error())
		}
	})

	go loop.Start()

	offer, err := peerConn.CreateOffer(nil)

	if err != nil {
		log.DefaultLogger().Error("CreateOffer", 1, err.Error())
		return
	}

	err = peerConn.SetLocalDescription(offer)
	if err != nil {
		log.DefaultLogger().Error("SetLocalDescription", 1, err.Error())
		return
	}

	err = c.SendMessage(conn.Message{Typ: conn.MsgType(offer.Type), Msg: offer.SDP}, time.Second)
	if err != nil {
		log.DefaultLogger().Error("SendMessage", 1, err.Error())
		return
	}

	select {}
}
