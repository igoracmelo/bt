package bt_test

import (
	"bytes"
	"testing"

	"github.com/igoracmelo/bt"
	"github.com/stretchr/testify/assert"
)

func Test_Message_Serialize_And_Unserialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		id    bt.MessageID
		setup func(msg *bt.Message)
	}{
		{
			id: bt.MsgKeepAlive,
			setup: func(msg *bt.Message) {
				msg.KeepAlive()
			},
		},
		{
			id: bt.MsgChoke,
			setup: func(msg *bt.Message) {
				msg.Choke()
			},
		},
		{
			id: bt.MsgUnchoke,
			setup: func(msg *bt.Message) {
				msg.Unchoke()
			},
		},
		{
			id: bt.MsgInterested,
			setup: func(msg *bt.Message) {
				msg.Interested()
			},
		},
		{
			id: bt.MsgNotInterested,
			setup: func(msg *bt.Message) {
				msg.NotInterested()
			},
		},
		{
			id: bt.MsgHave,
			setup: func(msg *bt.Message) {
				msg.Have(0)
			},
		},
		{
			id: bt.MsgBitfield,
			setup: func(msg *bt.Message) {
				bf := bt.BitfieldFromBytes([]byte{0xFF}, 0)
				msg.Bitfield(bf)
			},
		},
		{
			id: bt.MsgRequest,
			setup: func(msg *bt.Message) {
				msg.Request(0, 0, 0)
			},
		},
		{
			id: bt.MsgPiece,
			setup: func(msg *bt.Message) {
				msg.Piece(0, 0, []byte{0xFF, 0xAC, 0xFF})
			},
		},
		{
			id: bt.MsgCancel,
			setup: func(msg *bt.Message) {
				msg.Cancel(0, 0, 0)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id.String(), func(t *testing.T) {
			t.Parallel()

			msg := &bt.Message{}
			tt.setup(msg)
			b := msg.Bytes()

			msg = &bt.Message{}
			msg.FromBytes(b)

			assert.Equal(t, tt.id, msg.ID)
			assert.Equal(t, b, msg.Bytes())
		})
	}
}

func Test_Message_WriteTo(t *testing.T) {
	t.Parallel()

	msg := &bt.Message{}
	msg.KeepAlive()

	buf := &bytes.Buffer{}
	msg.WriteTo(buf)

	assert.Equal(t, msg.Bytes(), buf.Bytes())
}
