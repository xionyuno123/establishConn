package main

import (
	peerconn "establishConn/peerConn"
	"fmt"
)

func main() {
	peerConn := peerconn.NewWebRTCConn()

	err := peerConn.GatheringCandidate("stun.l.google.com:19302")
	if err != nil {
		fmt.Println(err.Error())
	}
}
