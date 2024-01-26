package bt_test

import (
	bt "shit"
	"slices"
	"testing"
)

func Test_Bitfield(t *testing.T) {
	t.Run("Has", func(t *testing.T) {
		bf := bt.BitfieldFromBytes([]byte{0b10101111}, 7)

		tests := []struct {
			index int64
			want  bool
		}{
			{-1, false},

			{0, true},
			{1, false},
			{2, true},
			{3, false},
			{4, true},
			{5, true},
			{6, true},

			// last bit is outside range 0-7
			{7, false},

			{15, false},
		}

		for _, tt := range tests {
			got := bf.Has(tt.index)
			if tt.want != got {
				t.Errorf("piece at %d - want: %v, got: %v", tt.index, tt.want, got)
			}
		}
	})

	t.Run("Set_setSameValues", func(t *testing.T) {
		bf := bt.BitfieldFromBytes([]byte{0b10101111, 0b11110000}, 0)

		before := bf.Pieces()
		after := bf.Pieces()

		for i := range before {
			i := int64(i)
			if !bf.Set(i, before[i]) {
				t.Errorf("failed to set %d to %v", i, before[i])
			}
			after[i] = bf.Has(i)
		}

		if !slices.Equal(before, after) {
			t.Errorf("before:\n%v\nafter:\n%v", before, after)
		}
	})

	t.Run("Set_FlipValuesTwice", func(t *testing.T) {
		bf := bt.BitfieldFromBytes([]byte{0b10101111, 0b11110000}, 0)

		before := bf.Pieces()
		after := bf.Pieces()

		for i := range before {
			i := int64(i)
			if !bf.Set(i, before[i]) {
				t.Errorf("failed to set %d to %v", i, before[i])
			}
			if !bf.Set(i, !bf.Has(i)) {
				t.Errorf("failed to set %d to %v", i, !bf.Has(i))
			}
			after[i] = bf.Has(i)
		}

		if !slices.Equal(before, after) {
			t.Errorf("before:\n%v\nafter:\n%v", before, after)
		}
	})

	t.Run("Set_EdgeCases", func(t *testing.T) {
		length := int64(35)
		bf := bt.BitfieldFromLength(length)

		tests := []struct {
			index int64
			want  bool
		}{
			{-1, false},
			{length, false},
			{length + 1, false},
		}

		for _, tt := range tests {
			got := bf.Set(tt.index, false)
			if got != tt.want {
				t.Errorf("set %d to %v - want: %v, got: %v", tt.index, false, tt.want, got)
			}
		}
	})
}
