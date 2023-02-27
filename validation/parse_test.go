package validation

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestParsePiece(t *testing.T) {
	tests := []struct {
		name        string
		pieceString string
		piece       board.Piece
		isErr       bool
	}{
		{"P", "P", board.WhitePawn, false},
		{"p", "p", board.BlackPawn, false},
		{"N", "N", board.WhiteKnight, false},
		{"n", "n", board.BlackKnight, false},
		{"B", "B", board.WhiteBishop, false},
		{"b", "b", board.BlackBishop, false},
		{"R", "R", board.WhiteRook, false},
		{"r", "r", board.BlackRook, false},
		{"Q", "Q", board.WhiteQueen, false},
		{"q", "q", board.BlackQueen, false},
		{"K", "K", board.WhiteKing, false},
		{"k", "k", board.BlackKing, false},
		{"abcd", "abcd", 0, true},
		{"", "", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parsePiece(tc.pieceString)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err: %v", err)
			}

			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}

			if res != tc.piece {
				t.Fatalf("want %s, got %s", tc.piece, res)
			}
		})
	}
}

func TestParseNewpiece(t *testing.T) {
	tests := []struct {
		name           string
		newpieceString string
		newpiece       board.Piece
		isErr          bool
	}{
		{"P", "P", 0, true},
		{"p", "p", 0, true},
		{"N", "N", board.WhiteKnight, false},
		{"n", "n", board.BlackKnight, false},
		{"B", "B", board.WhiteBishop, false},
		{"b", "b", board.BlackBishop, false},
		{"R", "R", board.WhiteRook, false},
		{"r", "r", board.BlackRook, false},
		{"Q", "Q", board.WhiteQueen, false},
		{"q", "q", board.BlackQueen, false},
		{"K", "K", 0, true},
		{"k", "k", 0, true},
		{"abcd", "abcd", 0, true},
		{"", "", 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseNewpiece(tc.newpieceString)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err: %v", err)
			}

			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}

			if res != tc.newpiece {
				t.Fatalf("want %s, got %s", tc.newpiece, res)
			}
		})
	}
}

func TestParseSquare(t *testing.T) {
	tests := []struct {
		name         string
		squareString string
		sq           square
		isErr        bool
	}{
		{"a1", "a1", newSquare(0), false},
		{"g8", "g8", newSquare(62), false},
		{"19", "19", newSquare(19), false},
		{"50", "50", newSquare(50), false},

		{"A1", "A1", newSquare(-1), true},
		{"", "", newSquare(-1), true},
		{"64", "64", newSquare(-1), true},
		{"-1", "-1", newSquare(-1), true},
		{"ab", "ab", newSquare(-1), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseSquare(tc.squareString)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err: %v", err)
			}

			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}

			if res != tc.sq {
				t.Fatalf("want %v, got %v", tc.sq, res)
			}
		})
	}
}
