package bt

import (
	"bytes"
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
	l, err := net.Listen("tcp", "localhost:54321")
	if err != nil {
		t.Fatal(l)
	}

	peer := Peer{net.IP("localhost"), 54321}
	clientId := [20]byte([]byte("It's a me, Igor! Ok?"))
	serverId := [20]byte([]byte("It's a me, Rogi! Ok?"))
	infoHash := [20]byte([]byte("It's a me, InfoHash!"))

	clientHs := &Handshake{
		PeerId:   clientId,
		InfoHash: infoHash,
	}

	serverHs := Handshake{
		InfoHash: infoHash,
		PeerId:   serverId,
	}

	var pc *PeerConn
	hsDone := make(chan struct{})

	go func() {
		// client connects and sends handshake
		var err error
		pc, err = NewPeerConn(peer, clientHs)
		if err != nil {
			panic(err)
		}
		close(hsDone)
	}()

	// server accepts connection
	server, err := l.Accept()
	if err != nil {
		panic(err)
	}

	// server reads handshake
	buf := make([]byte, 1024)
	n, err := server.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	// checking server received correct handshake
	gotHs, err := HandshakeFrom(buf[:n])
	if err != nil {
		t.Fatal(err)
	}
	if gotHs.InfoHash != clientHs.InfoHash {
		t.Fatalf("want infohash: %x, got: %x", clientHs.InfoHash, gotHs.InfoHash)
	}

	// server sends its handshake
	_, err = server.Write(serverHs.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	<-hsDone
	defer pc.Close()

	// client sends message
	err = pc.WriteMessage(Message{
		Id:      MsgChoke,
		Payload: []byte("hey there my friend"),
	})
	if err != nil {
		t.Fatal(err)
	}

	// server receive message
	msg, err := MessageFrom(server)
	if err != nil {
		t.Fatal(err)
	}

	// server sends back message
	_, err = server.Write(msg.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// client reads message
	msg2, err := pc.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	// validating message
	if msg2.Id != msg.Id {
		t.Errorf("want Id: %d, got: %d", msg.Id, msg2.Id)
	}
	if !bytes.Equal(msg2.Payload, msg.Payload) {
		t.Errorf("want payload: %x, got: %x", msg.Payload, msg2.Payload)
	}
}
