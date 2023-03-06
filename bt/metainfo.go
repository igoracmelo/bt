package bt

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"github.com/igoracmelo/bt/logger"
	bencode "github.com/jackpal/bencode-go"
)

type MetaInfo struct {
	Announce     string
	Comment      string
	CreationDate time.Time
	Info         Info
	InfoHash     [20]byte
}

type Info struct {
	PieceLength uint64
	Length      uint64
	Name        string
	Pieces      [][20]byte
}

func NewMetaInfo(r io.Reader) (*MetaInfo, error) {
	dto := struct {
		Announce     string `bencode:"announce"`
		Comment      string `bencode:"comment"`
		CreationDate uint64 `bencode:"creation date"`
		Info         struct {
			Length      uint64 `bencode:"length"`
			Name        string `bencode:"name"`
			PieceLength uint64 `bencode:"piece length"`
			Pieces      string `bencode:"pieces"`
		}
	}{}

	err := bencode.Unmarshal(r, &dto)
	if err != nil {
		return nil, err
	}

	infobuf := bytes.Buffer{}
	err = bencode.Marshal(&infobuf, dto.Info)
	if err != nil {
		return nil, err
	}

	x := infobuf.Bytes()
	_ = x
	infoHash := sha1.Sum(infobuf.Bytes())

	pieces := []byte(dto.Info.Pieces)
	if len(pieces)%sha1.Size != 0 {
		return nil, fmt.Errorf("Malformed torrent pieces with length %d", len(pieces))
	}

	numHashes := len(pieces) / sha1.Size
	if uint64(numHashes)*dto.Info.PieceLength != dto.Info.Length {
		return nil, fmt.Errorf("Malformed torrent length != #pieces * pieceLength (%d != %d * %d)", dto.Info.Length, numHashes, dto.Info.PieceLength)
	}

	torr := &MetaInfo{
		Announce:     dto.Announce,
		Comment:      dto.Comment,
		CreationDate: time.Unix(int64(dto.CreationDate), 0),
		Info: Info{
			Name:        dto.Info.Name,
			Length:      dto.Info.Length,
			Pieces:      make([][20]byte, numHashes),
			PieceLength: dto.Info.PieceLength,
		},
		InfoHash: infoHash,
	}

	logger.Debugf("InfoHash: %x", infoHash)

	for i := 0; i < numHashes; i++ {
		ini := i * sha1.Size
		end := (i + 1) * sha1.Size
		copy(torr.Info.Pieces[i][:], pieces[ini:end])
	}

	return torr, nil
}
