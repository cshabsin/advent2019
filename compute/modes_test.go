package compute

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestMakeModes(t *testing.T) {
	testcases := []struct {
		opcode int64
		want   modes
	}{
		{
			opcode: 1,
			want:   modes{},
		},
		{
			opcode: 101,
			want:   modes{1},
		},
		{
			opcode: 201,
			want:   modes{2},
		},
		{
			opcode: 2001,
			want:   modes{0, 2},
		},
		{
			opcode: 12001,
			want:   modes{0, 2, 1},
		},
	}
	for _, tc := range testcases {
		m := makeModes(tc.opcode)
		if !cmp.Equal(m, tc.want) {
			t.Errorf("makeModes(%d): want %v, got %v", tc.opcode, tc.want, m)
		}
	}
}
