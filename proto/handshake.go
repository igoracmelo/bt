package proto

import (
	"bytes"
	"fmt"
)

const ProtoLen = 19
const HandshakePrefix = "\x13BitTorrent protocol"

type Handshake struct {
	Extensions [8]byte
	InfoHash   [20]byte
	PeerId     [20]byte
}

func NewHandshake(b []byte) (*Handshake, error) {
	if len(b) != 68 {
		return nil, fmt.Errorf("handshake must have 68 bytes, but got %d", len(b))
	}

	if !bytes.HasPrefix(b, []byte(HandshakePrefix)) {
		return nil, fmt.Errorf("Invalid prefix: %s", string(b[:20]))
	}

	i := 20
	extensions := b[i : i+8]

	i += 8
	infoHash := b[i : i+20]

	i += 20
	peerId := b[i : i+20]

	hs := &Handshake{
		Extensions: [8]byte(extensions),
		InfoHash:   [20]byte(infoHash),
		PeerId:     [20]byte(peerId),
	}

	return hs, nil
}

func (hs *Handshake) Bytes() []byte {
	b := make([]byte, 0, 68)

	b = append(b, []byte(HandshakePrefix)...)
	b = append(b, hs.Extensions[:]...)
	b = append(b, hs.InfoHash[:]...)
	b = append(b, hs.PeerId[:]...)

	return b
}
