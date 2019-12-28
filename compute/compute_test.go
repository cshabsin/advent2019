package compute

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSizeUp(t *testing.T) {
	b := []int{}
	b = sizeUp(b, 2)
	if len(b) != 3 {
		t.Errorf("sizeUp: wanted len 3, got %d", len(b))
	}
}

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
		intcode := NewIntcode(tc.in, nil)
		got, _, err := intcode.Run()
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
			want:    nil,
			wantout: []int{37},
		},
		{
			desc:    "equal",
			inbuf:   []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:  []int{8},
			want:    nil,
			wantout: []int{1},
		},
		{
			desc:    "not equal",
			inbuf:   []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:  []int{7},
			want:    nil,
			wantout: []int{0},
		},
		{
			desc:    "position jump zero",
			inbuf:   []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:  []int{0},
			want:    nil,
			wantout: []int{0},
		},
		{
			desc:    "position jump nonzero",
			inbuf:   []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:  []int{5},
			want:    nil,
			wantout: []int{1},
		},
		{
			desc:    "immediate jump zero",
			inbuf:   []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:  []int{0},
			want:    nil,
			wantout: []int{0},
		},
		{
			desc:    "immediate jump non-zero",
			inbuf:   []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:  []int{4},
			want:    nil,
			wantout: []int{1},
		},
	}
	for _, tc := range testcases {
		intcode := NewIntcode(tc.inbuf, tc.inputs)
		got, gotout, err := intcode.Run()
		if err != nil {
			t.Errorf("Run(%q): unexpected error: %v", tc.desc, err)
			continue
		}
		if tc.want != nil {
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Run(%q): incorrect result, -want +got:\n%s", tc.desc, diff)
			}
		}
		if diff := cmp.Diff(tc.wantout, gotout); diff != "" {
			t.Errorf("Run(%q): incorrect output, -want +got:\n%s", tc.desc, diff)
		}
	}
}
