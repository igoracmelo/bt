package bt

import (
	"fmt"
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

func Test_PeerConn_HandshakeSendMessageAndReceiveMessage(t *testing.T) {
	infoHash := [20]byte([]byte("It's a me, InfoHash!"))

	yourHs := Handshake{
		InfoHash: infoHash,
		PeerId:   [20]byte([]byte("It's a me, Rogi! Ok?")),
	}

	ready := make(chan struct{})

	// dummy server
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

		buf := make([]byte, 1024)

		// read handshake
		n, err := urPeer.Read(buf)
		if err != nil {
			panic(err)
		}
		if n != 68 {
			panic("invalid message handshake length")
		}

		// send handshake
		urPeer.Write(yourHs.Bytes())

		// read message
		n, err = urPeer.Read(buf)
		if err != nil {
			panic(err)
		}

		if n != 4+1+19 {
			panic(fmt.Sprintf("expected message to be 4 + 1 + 19 bytes, but got %d", n))
		}

		msg := Message{
			Id:      MsgInterested,
			Payload: []byte("hey my man!"),
		}

		_, err = urPeer.Write(msg.Bytes())
		if err != nil {
			panic(err)
		}
	}()

	<-ready

	peer := Peer{net.IP("localhost"), 54321}
	myId := [20]byte([]byte("It's a me, Igor! Ok?"))

	pc, err := NewPeerConn(peer, myId, infoHash)
	if err != nil {
		t.Fatal(err)
	}

	err = pc.WriteMessage(Message{
		Id:      MsgChoke,
		Payload: []byte("hey there my friend"),
	})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := pc.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	if msg.Id != MsgInterested {
		t.Errorf("want: MsgInterested (%d), got: %d", MsgInterested, msg.Id)
	}

	if string(msg.Payload) != "hey my man!" {
		t.Errorf("want: hey my man!, got: %s", string(msg.Payload))
	}
}
