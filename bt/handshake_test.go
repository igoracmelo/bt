package bt

import (
	"strings"
	"testing"
)

func Test_HandshakeShouldResultInSameValue(t *testing.T) {
	hash := [20]byte([]byte("this is the infohash"))
	peerId := [20]byte([]byte("It's a me, Igor! Ok?"))

	hs1 := &Handshake{
		InfoHash: hash,
		PeerId:   peerId,
	}

	hs2, err := HandshakeFrom(hs1.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	if *hs1 != *hs2 {
		t.Fatalf("Handshakes do not match:\n%v\n%v", string(hs1.Bytes()), string(hs2.Bytes()))
	}
}

func Test_InvalidHandshakeBytes(t *testing.T) {
	_, err := HandshakeFrom([]byte("incorrect size"))
	if !strings.Contains(err.Error(), "68 bytes") {
		t.Errorf("error should contain '68 bytes' but got %v", err)
	}

	b := []byte("\x10Correct size but incorrect data being passed over here man 68 char.")
	_, err = HandshakeFrom(b)
	if !strings.Contains(err.Error(), "prefix") {
		t.Errorf("error should contain 'prefix' but got %v", err)
	}
}
