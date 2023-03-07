package bt

import (
	"fmt"
	"net"
)

type Peer struct {
	Ip   net.IP
	Port uint16
}

func (p *Peer) Address() string {
	return net.JoinHostPort(string(p.Ip), fmt.Sprint(p.Port))
}

type PeerConn struct {
	buf  []byte
	bf   Bitfield
	peer Peer
	conn net.Conn
	busy atomic.Bool
}

func NewPeerConn(peer Peer, myId [20]byte, infoHash [20]byte) (*PeerConn, error) {
	conn, err := net.DialTimeout("tcp", peer.Address(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	defer conn.SetDeadline(time.Time{})

	myHs := &Handshake{
		InfoHash: infoHash,
		PeerId:   myId,
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write(myHs.Bytes())
	if err != nil {
		return nil, err
	}

	hsBytes := make([]byte, 68)
	_, err = conn.Read(hsBytes)
	if err != nil {
		return nil, err
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	yourHs, err := HandshakeFrom(hsBytes)
	if err != nil {
		return nil, err
	}

	if yourHs.InfoHash != myHs.InfoHash {
		return nil, fmt.Errorf("infoHash do not match. want: %x, got: %x", myHs.InfoHash, yourHs.InfoHash)
	}

	return &PeerConn{
		buf:  make([]byte, 4096),
		peer: peer,
		conn: conn,
		busy: atomic.Bool{},
	}, nil
}
