package chess

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
			m, err := ParseUCI(tc.move)
			if err != nil {
				t.Fatal(err)
			}
			from := board.Sq(tc.from)
			to := board.Sq(tc.to)
			if from != m.FromSquare() || to != m.ToSquare() {
				t.Fatalf("want from %v to %v, got from %v to %v", from, to, m.FromSquare(), m.ToSquare())
			}
		})
	}
}

func TestParseUCIPromotion(t *testing.T) {
	tests := []struct {
		move  string
		piece board.Piece
	}{
		{"h7h8n", board.WhiteKnight},
		{"h2h1n", board.BlackKnight},
		{"e7e8q", board.WhiteQueen},
		{"d2d1q", board.BlackQueen},
		{"b7b8r", board.WhiteRook},
		{"a2a1r", board.BlackRook},
		{"a7a8b", board.WhiteBishop},
		{"b2b1b", board.BlackBishop},
	}

	for _, tc := range tests {
		t.Run(tc.move, func(t *testing.T) {
			m, err := ParseUCI(tc.move)
			if err != nil {
				t.Fatal(err)
			}

			prom, ok := m.(promotion)
			if !ok {
				t.Fatalf("want promotion, not some other kind of move")
			}
			if prom.promoteTo != tc.piece {
				t.Fatalf("want %s, got %s", tc.piece.String(), prom.promoteTo.String())
			}
		})
	}
}
