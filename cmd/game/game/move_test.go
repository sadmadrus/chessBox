package game

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestParseUCI(t *testing.T) {
	tests := []struct {
		move string
		from string
		to   string
	}{
		{"e2e4", "e2", "e4"},
	}

	for _, tc := range tests {
		t.Run(tc.move, func(t *testing.T) {
			m, err := parseUCI(tc.move)
			if err != nil {
				t.Fatal(err)
			}
			from := board.Sq(tc.from)
			to := board.Sq(tc.to)
			if from != m.fromSquare() || to != m.toSquare() {
				t.Fatalf("want from %v to %v, got from %v to %v", from, to, m.fromSquare(), m.toSquare())
			}
		})
	}
}
