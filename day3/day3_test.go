package main

import (
	"testing"
)

func TestIntersect(t *testing.T) {
	testcases := []struct {
		p1, p2   string
		wantDist int
	}{
		{
			p1:       "R8,U5,L5,D3",
			p2:       "U7,R6,D4,L4",
			wantDist: 6,
		},
		{
			p1:       "R75,D30,R83,U83,L12,D49,R71,U7,L72",
			p2:       "U62,R66,U55,R34,D71,R55,D58,R83",
			wantDist: 159,
		},
		{
			p1:       "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
			p2:       "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
			wantDist: 135,
		},
	}
	for _, tc := range testcases {
		p1, err := makePath(tc.p1)
		if err != nil {
			t.Errorf("makePath(%q): %v", tc.p1, err)
			continue
		}
		res, err := p1.intersect(tc.p2)
		if err != nil {
			t.Errorf("intersect(%q): %v", tc.p2, err)
			continue
		}
		if res != tc.wantDist {
			t.Errorf("intersect(%q): want %d, got %d", tc.p2, tc.wantDist, res)
		}

	}
}
