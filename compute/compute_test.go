package compute

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestComputeNoIO(t *testing.T) {
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
		{
			desc: "immediate",
			in:   []int{101, 4, 0, 0, 99},
			want: []int{105, 4, 0, 0, 99},
		},
	}
	for _, tc := range testcases {
		got, _, err := Run(tc.in, nil)
		if err != nil {
			t.Errorf("Run(%q): unexpected error: %v", tc.desc, err)
			continue
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Run(%q): incorrect result, -want +got:\n%s", tc.desc, diff)
		}
	}
}

func TestCompute(t *testing.T) {
	testcases := []struct {
		desc    string
		inbuf   []int
		inputs  []int
		want    []int
		wantout []int
	}{
		{
			desc:    "input",
			inbuf:   []int{3, 3, 99, 0},
			inputs:  []int{1},
			want:    []int{3, 3, 99, 1},
			wantout: []int{},
		},
		{
			desc:    "output",
			inbuf:   []int{104, 37, 99},
			inputs:  []int{},
			want:    []int{104, 37, 99},
			wantout: []int{37},
		},
	}
	for _, tc := range testcases {
		got, gotout, err := Run(tc.inbuf, tc.inputs)
		if err != nil {
			t.Errorf("Run(%q): unexpected error: %v", tc.desc, err)
			continue
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Run(%q): incorrect result, -want +got:\n%s", tc.desc, diff)
		}
		if diff := cmp.Diff(tc.wantout, gotout); diff != "" {
			t.Errorf("Run(%q): incorrect output, -want +got:\n%s", tc.desc, diff)
		}
	}
}
