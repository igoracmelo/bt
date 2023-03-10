package bt

import (
	"bytes"
	"testing"
)

func Test_MessageShouldResultInSameValue(t *testing.T) {
	m1 := &Message{
		Id:      MsgBitfield,
		Payload: []byte("hey there"),
	}

	r := bytes.NewReader(m1.Bytes())

	m2, err := MessageFrom(r)
	if err != nil {
		t.Fatal(err)
	}

	if m1.Id != m2.Id || !bytes.Equal(m1.Payload, m2.Payload) {
		t.Log("expected to be equal:")
		t.Logf("msg 1: %x", m1)
		t.Logf("msg 2: %x", m2)
		t.FailNow()
	}
}
