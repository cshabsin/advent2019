package compute

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFile(t *testing.T) {
	testcases := []struct {
		content string
		want    []int64
	}{
		{"1", []int64{1}},
		{"1,2", []int64{1, 2}},
		{"104,1125899906842624,99", []int64{104, 1125899906842624, 99}},
	}
	for _, tc := range testcases {
		got, err := ParseFile([]byte(tc.content))
		if err != nil {
			t.Errorf("ParseFile(%q): %v", tc.content, err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("ParseFile(%q): got %v, want %v", tc.content, got, tc.want)
		}
		strs := []string{}
		for _, v := range got {
			strs = append(strs, fmt.Sprintf("%d", v))
		}
		if gotstr := strings.Join(strs, ","); gotstr != tc.content {
			t.Errorf("ParseFile(%q) to string: got %v, want %v", tc.content, gotstr, tc.content)
		}
	}
}

func TestSizeUp(t *testing.T) {
	b := []int64{}
	b = sizeUp(b, 2)
	if len(b) != 3 {
		t.Errorf("sizeUp: wanted len 3, got %d", len(b))
	}
}

func TestComputeNoIO(t *testing.T) {
	testcases := []struct {
		desc string
		in   []int64
		want []int64
	}{
		{
			desc: "first",
			in:   []int64{1, 0, 0, 0, 99},
			want: []int64{2, 0, 0, 0, 99},
		},
		{
			desc: "immediate",
			in:   []int64{101, 4, 0, 0, 99},
			want: []int64{105, 4, 0, 0, 99},
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
		inbuf   []int64
		inputs  []int64
		want    []int64
		wantout []int64
	}{
		{
			desc:    "input",
			inbuf:   []int64{3, 3, 99, 0},
			inputs:  []int64{1},
			want:    []int64{3, 3, 99, 1},
			wantout: []int64{},
		},
		{
			desc:    "output",
			inbuf:   []int64{104, 37, 99},
			inputs:  []int64{},
			want:    nil,
			wantout: []int64{37},
		},
		{
			desc:    "equal",
			inbuf:   []int64{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:  []int64{8},
			want:    nil,
			wantout: []int64{1},
		},
		{
			desc:    "not equal",
			inbuf:   []int64{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8},
			inputs:  []int64{7},
			want:    nil,
			wantout: []int64{0},
		},
		{
			desc:    "position jump zero",
			inbuf:   []int64{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:  []int64{0},
			want:    nil,
			wantout: []int64{0},
		},
		{
			desc:    "position jump nonzero",
			inbuf:   []int64{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9},
			inputs:  []int64{5},
			wantout: []int64{1},
		},
		{
			desc:    "immediate jump zero",
			inbuf:   []int64{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:  []int64{0},
			wantout: []int64{0},
		},
		{
			desc:    "immediate jump non-zero",
			inbuf:   []int64{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1},
			inputs:  []int64{4},
			wantout: []int64{1},
		},
		{
			desc:    "relative copy self",
			inbuf:   []int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
			wantout: []int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
		{
			desc:    "calculated large number",
			inbuf:   []int64{1102, 34915192, 34915192, 7, 4, 7, 99, 0},
			wantout: []int64{1219070632396864},
		},
		{
			desc:    "large number",
			inbuf:   []int64{104, 1125899906842624, 99},
			wantout: []int64{1125899906842624},
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
