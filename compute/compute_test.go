package compute

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestCompute(t *testing.T) {
	testcases := []struct {
		desc string
		in   []int
		want []int
	}{
		{
			desc: "first",
			in:   []int{1, 0, 0, 0, 99},
			want: []int{2, 0, 0, 0, 99},
		},
	}
	for _, tc := range testcases {
		got, err := Run(tc.in)
		if err != nil {
			t.Errorf("Run(%q): unexpected error: %v", tc.desc, err)
			continue
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Run(%q): incorrect result, -want +got:\n%s", tc.desc, diff)
		}
	}
}
