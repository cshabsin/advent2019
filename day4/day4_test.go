package main

import (
	"testing"
)

func TestMatch(t *testing.T) {
	testcases := []struct{
		in int
		want bool
	}{
		{111111, true},
		{223450, false},
		{223457, true},
		{123789, false},
	}
	for _, tc := range testcases {
		if got, want := match(tc.in), tc.want; got != want {
			t.Errorf("match(%d): got %v, want %v", tc.in, got, want)
		}
	}
}
