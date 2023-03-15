package validation

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestParsePiece(t *testing.T) {
	tests := []struct {
		name        string
		pieceString string
		pieceType   string
		piece       board.Piece
		isErr       bool
	}{
		// piece pieceType
		{"P", "P", "piece", board.WhitePawn, false},
		{"p", "p", "piece", board.BlackPawn, false},
		{"N", "N", "piece", board.WhiteKnight, false},
		{"n", "n", "piece", board.BlackKnight, false},
		{"B", "B", "piece", board.WhiteBishop, false},
		{"b", "b", "piece", board.BlackBishop, false},
		{"R", "R", "piece", board.WhiteRook, false},
		{"r", "r", "piece", board.BlackRook, false},
		{"Q", "Q", "piece", board.WhiteQueen, false},
		{"q", "q", "piece", board.BlackQueen, false},
		{"K", "K", "piece", board.WhiteKing, false},
		{"k", "k", "piece", board.BlackKing, false},
		{"abcd", "abcd", "piece", 0, true},
		{"", "", "piece", 0, true},

		// newpiece pieceType
		{"P", "P", "newpiece", 0, true},
		{"p", "p", "newpiece", 0, true},
		{"N", "N", "newpiece", board.WhiteKnight, false},
		{"n", "n", "newpiece", board.BlackKnight, false},
		{"B", "B", "newpiece", board.WhiteBishop, false},
		{"b", "b", "newpiece", board.BlackBishop, false},
		{"R", "R", "newpiece", board.WhiteRook, false},
		{"r", "r", "newpiece", board.BlackRook, false},
		{"Q", "Q", "newpiece", board.WhiteQueen, false},
		{"q", "q", "newpiece", board.BlackQueen, false},
		{"K", "K", "newpiece", 0, true},
		{"k", "k", "newpiece", 0, true},
		{"abcd", "abcd", "newpiece", 0, true},
		{"", "", "newpiece", 0, false},

		// wrong pieceType
		{"R", "R", "somepiece", 0, true},
		{"r", "r", "", 0, true},
		{"Q", "Q", "asdfg", 0, true},
		{"q", "q", "12345", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parsePiece(tc.pieceString, tc.pieceType)
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
