package bt

import (
	"fmt"
)

type Bitfield struct {
	len  int64
	data []byte
}

func BitfieldFromLength(length int64) Bitfield {
	return Bitfield{
		len:  length,
		data: make([]byte, (length/8)+1),
	}
}

func BitfieldFromBytes(b []byte, length int64) Bitfield {
	if length == 0 {
		length = int64(len(b) * 8)
	}
	return Bitfield{
		len:  length,
		data: b,
	}
}

func (bf Bitfield) Has(i int64) bool {
	if i < 0 || i >= bf.len {
		return false
	}

	bi := int(i / 8)
	b := bf.data[bi]

	bit := (b >> (7 - (i % 8))) & 1

	return bit == 1
}

func (bf Bitfield) Set(i int64, set bool) bool {
	if i < 0 || i >= bf.len {
		return false
	}

	bit := byte(0)
	if set {
		bit = 1
	}

	bi := int(i / 8)
	b := bf.data[bi]
	b &= bit << (7 - (i % 8))

	return true
}

func (bf Bitfield) Complete() bool {
	for _, has := range bf.Pieces() {
		if !has {
			return false
		}
	}
	return true
}

func (bf Bitfield) Pieces() []bool {
	pieces := make([]bool, bf.len)

	for i := int64(0); i < bf.len; i++ {
		pieces[i] = bf.Has(i)
	}

	return pieces
}

func (bf Bitfield) Bytes() []byte {
	return bf.data
}

func (bf Bitfield) String() string {
	s := "[ "
	for _, b := range bf.data {
		s += fmt.Sprintf("%08b ", b)
	}
	s += "]"
	return s
}

func (bf Bitfield) DebugString() string {
	s := ""
	for i, b := range bf.data {
		s += fmt.Sprintf("%08b - byte %d - bit %d\n", b, i+1, (i+1)*8)
	}
	return s
}
