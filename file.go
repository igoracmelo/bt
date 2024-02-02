package bt

import (
	"context"
	"os"
)

type File struct {
	ctx         context.Context
	cancel      context.CancelFunc
	pieceLength int64
	dst         *os.File
}
