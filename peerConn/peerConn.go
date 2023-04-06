package peerconn

import (
	"establishConn/conn"
	"fmt"
	"sync"

	"github.com/pion/webrtc/v3"
)

type WebRTCConn struct {
	peerConn          *webrtc.PeerConnection
	signal            conn.Conn
	pendingCandidates []*webrtc.ICECandidate
	candidateMtx      sync.Mutex
}

func NewWebRTCConn() *WebRTCConn {
	return &WebRTCConn{}
}

func (c *WebRTCConn) PendingCandidate(candidate *webrtc.ICECandidate) {
	c.candidateMtx.Lock()
	defer c.candidateMtx.Unlock()

	c.pendingCandidates = append(c.pendingCandidates, candidate)
}

func (c *WebRTCConn) GatheringCandidate(stunServer string) error {
	peerConn, err := webrtc.NewPeerConnection(webrtc.Configuration{})

	if err != nil {
		return err
	}

	c.peerConn = peerConn
	c.pendingCandidates = make([]*webrtc.ICECandidate, 0)

	peerConn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		c.PendingCandidate(candidate)
	})

	offer, err := peerConn.CreateOffer(nil)
	if err != nil {
		return err
	}

	err = peerConn.SetLocalDescription(offer)

	if err != nil {
		return err
	}

	promise := webrtc.GatheringCompletePromise(peerConn)

	<-promise

	// newCandidates := make([]*webrtc.Certificate, 0)
	for _, candidate := range c.pendingCandidates {
		fmt.Println(candidate.Address)
		fmt.Println(candidate.Port)
		fmt.Println(candidate.Protocol)

		fmt.Println(candidate.RelatedAddress)
		fmt.Println(candidate.RelatedPort)
	}

	return nil
}

func getCandidateByStun(candidate)
