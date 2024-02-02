package bt_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/igoracmelo/bt"
	"github.com/stretchr/testify/assert"
)

func Test_Bencode_ReadString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
	}{
		{"empty string", ""},
		{"one character", "a"},
		{"multiple characters", "a longer string"},
		{"unicode", "é"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.str, func(t *testing.T) {
			fstr := fmt.Sprintf("%d:%s", len(tt.str), tt.str)
			s, err := bt.ReadString(bufio.NewReader(strings.NewReader(fstr)))
			assert.NoError(t, err)
			assert.Equal(t, tt.str, s)
		})
	}
}
