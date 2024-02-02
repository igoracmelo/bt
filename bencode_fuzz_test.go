package bt_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/igoracmelo/bt"
	"github.com/stretchr/testify/assert"
)

func Fuzz_Bencode_ReadString(f *testing.F) {
	f.Add("")
	f.Fuzz(func(t *testing.T, got string) {
		fstr := fmt.Sprintf("%d:%s", len(got), got)
		want, err := bt.ReadString(bufio.NewReader(strings.NewReader(fstr)))
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func Test_Bencode_ReadInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		num  int64
	}{
		{"zero", 0},
		{"positive", 200},
		{"negative", -123213},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fstr := fmt.Sprintf("i%de", tt.num)
			n, err := bt.ReadInt64(bufio.NewReader(strings.NewReader(fstr)))
			assert.NoError(t, err)
			assert.Equal(t, tt.num, n)
		})
	}
}

func Fuzz_Bencode_ReadInt64(f *testing.F) {
	f.Add(int64(1))
	f.Fuzz(func(t *testing.T, i int64) {
		fstr := fmt.Sprintf("i%de", i)
		n, err := bt.ReadInt64(bufio.NewReader(strings.NewReader(fstr)))
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	})
}
