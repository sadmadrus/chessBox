package validation

import "testing"

func TestMoveBishop(t *testing.T) {
	tests := []struct {
		name string
		from square
		to   square
		isOk bool
	}{
		{"up-right", newSquare(9), newSquare(18), true},
		{"far", newSquare(0), newSquare(63), true},
		{"up-left", newSquare(10), newSquare(17), true},
		{"down-right", newSquare(9), newSquare(2), true},
		{"down-left", newSquare(10), newSquare(1), true},
		{"horizontal", newSquare(9), newSquare(12), false},
		{"vertical", newSquare(9), newSquare(25), false},
		{"knight", newSquare(9), newSquare(24), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := MoveBishop(tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}
