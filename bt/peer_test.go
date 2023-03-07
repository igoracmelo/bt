package bt

import (
	"log"
	"net"
	"testing"
)

func Test_Peer_Address(t *testing.T) {
	p := Peer{
		Ip:   net.IP("123.456.789.0"),
		Port: 54321,
	}

	want := "123.456.789.0:54321"
	got := p.Address()
	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func Test_PeerConn_New(t *testing.T) {
	infoHash := [20]byte([]byte("It's a me, InfoHash!"))

	yourHs := Handshake{
		InfoHash: infoHash,
		PeerId:   [20]byte([]byte("It's a me, Rogi! Ok?")),
	}

	ready := make(chan struct{})

	// dummy servers that only accepts a connection and sends the expected infohash
	go func() {
		l, err := net.Listen("tcp", "localhost:54321")
		if err != nil {
			panic(err)
		}
		close(ready)

		urPeer, err := l.Accept()
		if err != nil {
			panic(err)
		}
		urPeer.Write(yourHs.Bytes())
	}()

	<-ready

	peer := Peer{net.IP("localhost"), 54321}
	myId := [20]byte([]byte("It's a me, Igor! Ok?"))

	pc, err := NewPeerConn(peer, myId, infoHash)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()
}
