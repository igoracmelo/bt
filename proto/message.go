package proto

import "fmt"

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

func MessageFrom(b []byte) (*Message, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("invalid message with length 0")
	}

	p := []byte{}
	if len(b) > 1 {
		p = b[1:]
	}

	msg := &Message{
		Id:      MsgId(b[0]),
		Payload: p,
	}

	return msg, nil
}

func (m *Message) Bytes() []byte {
	b := make([]byte, len(m.Payload)+1)
	b[0] = byte(m.Id)
	copy(b[1:], m.Payload)

	return b
}
