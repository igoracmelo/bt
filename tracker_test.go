package main

import (
	"encoding/hex"
	"testing"
)

func Test_FindPeersHTTP(t *testing.T) {
	infoHash, err := hex.DecodeString("0763b757341956b9bf5994e4a992c3646f9f19bd")
	if err != nil {
		t.Fatal(err)
	}

	peers, err := FindPeersHTTP("http://bttracker.debian.org:6969/announce", FindPeersParams{
		InfoHash:   string(infoHash),
		PeerId:     "It's a me, Igor! Ok?",
		Port:       54321,
		Uploaded:   0,
		Downloaded: 0,
		Left:       0,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(peers) > 0 {
		t.Logf("peer Ip: %s, Port: %d", string(peers[0].Ip), peers[0].Port)
	}
}
