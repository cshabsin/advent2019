package maze

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetVal(t *testing.T) {
	testcases := []struct {
		x, y, val int
		want      Board
	}{
		{
			x:   0,
			y:   0,
			val: 1,
			want: Board{
				board: [][]int{{1}},
			},
		},
		{
			x:   1,
			y:   1,
			val: 1,
			want: Board{
				board: [][]int{nil, {0, 1}},
			},
		},
		{
			x:   -1,
			y:   -1,
			val: 1,
			want: Board{
				offsetX: 1,
				offsetY: 1,
				board:   [][]int{{1}},
			},
		},
	}
	for _, tc := range testcases {
		b := NewBoard()
		b.SetVal(tc.x, tc.y, tc.val)

		if diff := cmp.Diff(tc.want, *b, cmp.AllowUnexported(Board{})); diff != "" {
			t.Errorf("unexpected board:\n%s", diff)
		}
	}
}
