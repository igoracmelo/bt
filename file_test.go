package bt_test

import (
	"os"
	bt "shit"
	"shit/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test(t *testing.T) {
	pieceLength := 21
	piece := make([]byte, pieceLength)

	file, err := os.CreateTemp("", "output")
	if err != nil {
		t.Fatal(err)
	}

	peers := mocks.NewPeers(t)
	torrentFile := bt.NewFile(peers, 100, int64(pieceLength), file)

	peers.
		On("GetPiece", mock.Anything, mock.Anything).
		Return(func(piece []byte, i int64) error {
			for i := range piece {
				piece[i] = byte(i)
			}
			return nil
		})

	for i := 0; i < 5; i++ {
		err = torrentFile.ReadPiece(piece, 0)
		assert.NoError(t, err)
	}

	for i := 0; i < 5; i++ {
		n, err := file.Read(piece)
		assert.NoError(t, err)

		for _, b := range piece[:n] {
			assert.Equal(t, byte(i), b)
		}
	}
}
