package fuel

import "testing"

func TestFuel(t *testing.T) {
	testcases := []struct {
		mass     int
		wantFuel int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tc := range testcases {
		if got := Fuel(tc.mass); got != tc.wantFuel {
			t.Errorf("Fuel(%d): want %d, got %d", tc.mass, tc.wantFuel, got)
		}
	}
}
