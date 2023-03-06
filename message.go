package main

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
)

type Message struct {
	Id      MsgId
	Payload []byte
}

func MessageFrom(r io.Reader) (*Message, error) {
	l := make([]byte, 4)
	_, err := io.ReadFull(r, l)
	if err != nil {
		return nil, err
	}

	msgbuf := make([]byte, binary.BigEndian.Uint32(l))
	_, err = io.ReadFull(r, msgbuf)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Id:      MsgId(msgbuf[0]),
		Payload: msgbuf[1:],
	}

	return msg, nil
}

func (m *Message) Bytes() []byte {
	// 4 bytes: length, 1 byte: message type, +payload length
	b := make([]byte, 0, 4+1+len(m.Payload))

	b = binary.BigEndian.AppendUint32(b, 1+uint32(len(m.Payload)))
	b = append(b, byte(m.Id))
	b = append(b, m.Payload...)

	return b
}
