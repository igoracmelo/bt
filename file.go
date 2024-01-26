package bt

import (
	"context"
	"fmt"
	"io"
	"log"
)

type FileReader struct {
	file   File
	ctx    context.Context
	length int64
	cursor int64
}

var _ io.ReadSeekCloser = &FileReader{}

// Close implements io.Closer.
func (*FileReader) Close() error {
	panic("unimplemented")
}

// Read implements io.Reader.
func (f *FileReader) Read(p []byte) (n int, err error) {
	panic("unimplemented")

	// index := f.cursor / f.file.PieceLength()
	// start := f.cursor % f.file.PieceLength()

	// piece, err := f.file.ReadPiece(index)
}

// Seek implements io.Seeker.
func (f *FileReader) Seek(offset int64, whence int) (int64, error) {
	newCursor := int64(-1)

	switch whence {
	case io.SeekStart:
		newCursor = offset
	case io.SeekCurrent:
		newCursor = f.cursor + offset
	case io.SeekEnd:
		newCursor = f.length - offset
	}

	if newCursor < 0 {
		return 0, fmt.Errorf("invalid seek: offset %d, whence %d", offset, whence)
	}

	f.cursor = newCursor
	return newCursor, nil
}

type File struct {
	peers       Peers
	pieceLength int64
	fileLength  int64
	bf          Bitfield
	output      ReadWriteSeekCloser
}

func NewFile(peers Peers, fileLength int64, pieceLength int64, output ReadWriteSeekCloser) *File {
	return &File{
		peers:       peers,
		pieceLength: pieceLength,
		fileLength:  fileLength,
		bf:          BitfieldFromLength(fileLength),
		output:      output,
	}
}

func (f File) PieceLength() int64 {
	return f.pieceLength
}

func (f File) writePiece(piece []byte, index int64) (int, error) {
	if len(piece) != int(f.pieceLength) {
		return 0, fmt.Errorf("piece with unexpected length %dB, wanted %dB", len(piece), f.pieceLength)
	}

	pieceOffset := f.pieceLength * index

	_, err := f.output.Seek(int64(pieceOffset), io.SeekStart)
	if err != nil {
		return 0, err
	}

	return f.output.Write(piece)
}

func (f File) readPiece(piece []byte, index int64) (int, error) {
	if !f.bf.Has(index) {
		return 0, fmt.Errorf("piece not available")
	}

	pieceOffset := f.pieceLength * index

	_, err := f.output.Seek(int64(pieceOffset), io.SeekStart)
	if err != nil {
		return 0, err
	}

	return f.output.Read(piece)
}

func (f File) ReadPiece(piece []byte, index int64) error {
	if f.bf.Has(index) {
		_, err := f.readPiece(piece, index)
		return err
	}

	err := f.peers.GetPiece(piece, index)
	if err != nil {
		return err
	}

	_, err = f.writePiece(piece, index)
	if err != nil {
		return err
	}

	ok := f.bf.Set(index, true)
	if !ok {
		log.Fatalf("failed to set piece %d. bitfield: %s", index, f.bf.String())
	}

	return nil
}

// Peers interface abstract away the complexity of
// managing and communicating with multiple peers,
// potentially concurrently.
type Peers interface {
	GetPiece(piece []byte, index int64) error
}

type ReadWriteSeekCloser interface {
	io.ReadWriteSeeker
	io.Closer
}
