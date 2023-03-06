package main

type Bitfield []byte

func (bf Bitfield) Get(n int) bool {
	b := bf[n/8]
	bit := b >> (7 - n%8)
	bit &= 1
	return bit == 1
}

func (bf Bitfield) Set(n int, val bool) {
	b := bf[n/8]
	if val {
		b |= 1 << (7 - n%8)
	} else {
		b &= ^(1 << (7 - n%8))
	}
	bf[n/8] = b
}
