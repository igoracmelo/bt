package bt

import (
	"encoding/binary"
	"io"
)

type MsgId uint8

const (
	MsgChoke         MsgId = 0
	MsgUnchoke       MsgId = 1
	MsgInterested    MsgId = 2
	MsgNotInterested MsgId = 3
	MsgHave          MsgId = 4
	MsgBitfield      MsgId = 5
	MsgRequest       MsgId = 6
	MsgPiece         MsgId = 7
	MsgCancel        MsgId = 8
	MsgKeepAlive     MsgId = 22
)

type Message struct {
	Id      MsgId
	Payload []byte
}

func MessageFrom(r io.Reader) (*Message, error) {
	lenbuf := make([]byte, 4)
	_, err := io.ReadFull(r, lenbuf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenbuf)
	if length == 0 {
		return &Message{
			Id:      MsgKeepAlive,
			Payload: []byte{},
		}, nil
	}

	msgbuf := make([]byte, length)
	_, err = io.ReadFull(r, msgbuf)
	if err != nil {
		return nil, err
	}

	return &Message{
		Id:      MsgId(msgbuf[0]),
		Payload: msgbuf[1:],
	}, nil
}

func (m *Message) Bytes() []byte {
	if m.Id == MsgKeepAlive {
		return []byte{0, 0, 0, 0}
	}

	// 4 bytes: length, 1 byte: message type, +payload length
	b := make([]byte, 0, 4+1+len(m.Payload))

	b = binary.BigEndian.AppendUint32(b, 1+uint32(len(m.Payload)))
	b = append(b, byte(m.Id))
	b = append(b, m.Payload...)

	return b
}
