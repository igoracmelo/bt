package bt

import "testing"

func Test_BitfieldGet(t *testing.T) {
	bf := Bitfield([]byte{0b10010001})

	tests := []struct {
		i    int
		want bool
	}{
		{0, true},
		{1, false},
		{2, false},
		{3, true},
		{4, false},
		{5, false},
		{6, false},
		{7, true},
	}

	for _, tt := range tests {
		got := bf.Get(tt.i)
		if got != tt.want {
			t.Errorf("position %d: want %v, got %v", tt.i, tt.want, got)
		}
	}
}

func Test_BitfieldSet(t *testing.T) {
	bf := Bitfield([]byte{0b11110000})
	bf.Set(0, false)
	bf.Set(1, true)
	bf.Set(2, false)
	bf.Set(3, true)
	bf.Set(4, false)
	bf.Set(5, true)
	bf.Set(6, false)
	bf.Set(7, true)

	tests := []struct {
		i    int
		want bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
	}

	for _, tt := range tests {
		got := bf.Get(tt.i)
		if got != tt.want {
			t.Errorf("position %d: want %v, got %v", tt.i, tt.want, got)
		}
	}
}
