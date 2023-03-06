package main

import (
	"encoding/hex"
	"os"
	"testing"
	"time"
)

func Test_From(t *testing.T) {
	f, err := os.Open("./samples/debian.torrent")
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewMetaInfo(f)
	if err != nil {
		t.Fatal(err)
	}

	want := MetaInfo{
		Announce:     "http://bttracker.debian.org:6969/announce",
		Comment:      `"Debian CD from cdimage.debian.org"`,
		CreationDate: time.Unix(1671279452, 0),
		Info: Info{
			PieceLength: 262144,
			Name:        "debian-11.6.0-amd64-DVD-1.iso",
			Length:      3909091328,
			Pieces:      nil,
		},
		InfoHash: [20]byte{},
	}

	if want.Announce != got.Announce {
		t.Errorf("want announce: %s, got: %s", want.Announce, got.Announce)
	}

	if want.Comment != got.Comment {
		t.Errorf("want comment: %s, got: %s", want.Comment, got.Comment)
	}

	if want.CreationDate.Compare(got.CreationDate) != 0 {
		t.Errorf("want creation date: %v, got: %v", want.CreationDate, got.CreationDate)
	}

	if want.Info.PieceLength != got.Info.PieceLength {
		t.Errorf("want piece length: %d, got: %d", want.Info.PieceLength, got.Info.PieceLength)
	}

	if want.Info.Length != got.Info.Length {
		t.Errorf("want length: %d, got: %d", want.Info.Length, got.Info.Length)
	}

	if want.Info.Name != got.Info.Name {
		t.Errorf("want name: %s, got: %s", want.Info.Name, got.Info.Name)
	}

	wantHash := "0763b757341956b9bf5994e4a992c3646f9f19bd"
	gotHash := hex.EncodeToString(got.InfoHash[:])

	if wantHash != gotHash {
		t.Errorf("want InfoHash: %s, got: %s", wantHash, gotHash)
	}
}
