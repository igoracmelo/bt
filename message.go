package bt

import (
	"bytes"
	"encoding/binary"
	"io"
)

type MessageID byte

func (id MessageID) String() string {
	switch id {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	case MsgKeepAlive:
		return "KeepAlive"
	}
	return "Unknown"
}

const (
	MsgChoke         MessageID = 0
	MsgUnchoke       MessageID = 1
	MsgInterested    MessageID = 2
	MsgNotInterested MessageID = 3
	MsgHave          MessageID = 4
	MsgBitfield      MessageID = 5
	MsgRequest       MessageID = 6
	MsgPiece         MessageID = 7
	MsgCancel        MessageID = 8

	// not an actual protocol message id
	MsgKeepAlive MessageID = 255
)

type Message struct {
	ID      MessageID
	Payload []byte
}

var _ io.WriterTo = (*Message)(nil)

// WriteTo implements io.WriterTo.
func (msg *Message) WriteTo(w io.Writer) (int64, error) {
	b := msg.Bytes()
	n, err := w.Write(b)
	return int64(n), err
}

func (msg *Message) Bytes() []byte {
	if msg.ID == MsgKeepAlive {
		return []byte{0, 0, 0, 0}
	}

	b := make([]byte, 4+1+len(msg.Payload))
	length := uint32(len(msg.Payload) + 1)

	binary.BigEndian.PutUint32(b, length)
	b[4] = byte(msg.ID)
	copy(b[5:], msg.Payload)

	return b
}

func (msg *Message) Read(r io.Reader) error {
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBytes)
	if err != nil {
		return err
	}

	length := binary.BigEndian.Uint32(lengthBytes)
	if length == 0 {
		msg.ID = MsgKeepAlive
		return nil
	}

	buf := make([]byte, length)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}

	msg.ID = MessageID(buf[0])
	if n < 2 {
		return nil
	}

	msg.Payload = buf[1:]
	return nil
}

func (msg *Message) FromBytes(b []byte) {
	msg.Read(bytes.NewReader(b))
}

func (msg *Message) KeepAlive() {
	msg.ID = MsgKeepAlive
	msg.Payload = nil
}

func (msg *Message) Choke() {
	msg.ID = MsgChoke
	msg.Payload = nil
}

func (msg *Message) Unchoke() {
	msg.ID = MsgUnchoke
	msg.Payload = nil
}

func (msg *Message) Interested() {
	msg.ID = MsgInterested
	msg.Payload = nil
}

func (msg *Message) NotInterested() {
	msg.ID = MsgNotInterested
	msg.Payload = nil
}

func (msg *Message) Have(index int) {
	msg.ID = MsgHave
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(index))
	msg.Payload = buf
}

func (msg *Message) Bitfield(bitfield Bitfield) {
	msg.ID = MsgBitfield
	msg.Payload = bitfield.Bytes()
}

func (msg *Message) Request(index, begin, length int) {
	msg.ID = MsgRequest
	buf := make([]byte, 4+4+4)
	binary.BigEndian.PutUint32(buf, uint32(index))
	binary.BigEndian.PutUint32(buf[4:], uint32(begin))
	binary.BigEndian.PutUint32(buf[8:], uint32(length))
	msg.Payload = buf
}

func (msg *Message) Piece(index, begin int, block []byte) {
	msg.ID = MsgPiece
	buf := make([]byte, 4+4+len(block))
	binary.BigEndian.PutUint32(buf, uint32(index))
	binary.BigEndian.PutUint32(buf[4:], uint32(begin))
	copy(buf[8:], block)
	msg.Payload = buf
}

func (msg *Message) Cancel(index, begin, length int) {
	msg.ID = MsgCancel
	buf := make([]byte, 4+4+4)
	binary.BigEndian.PutUint32(buf, uint32(index))
	binary.BigEndian.PutUint32(buf[4:], uint32(begin))
	binary.BigEndian.PutUint32(buf[8:], uint32(length))
	msg.Payload = buf
}
