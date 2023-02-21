package validation

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestMove(t *testing.T) {
	tests := []struct {
		name  string
		piece board.Piece
		from  square
		to    square
		isOk  bool
	}{
		{"P 8-24", board.WhitePawn, newSquare(8), newSquare(24), true},
		{"p 24-8", board.BlackPawn, newSquare(24), newSquare(8), false},
		{"N 2-34", board.WhiteKnight, newSquare(2), newSquare(34), false},
		{"n 18-35", board.BlackKnight, newSquare(18), newSquare(35), true},
		{"B 9-18", board.WhiteBishop, newSquare(9), newSquare(18), true},
		{"b 9-25", board.BlackBishop, newSquare(9), newSquare(25), false},
		{"R 9-63", board.WhiteRook, newSquare(9), newSquare(63), false},
		{"r 63-7", board.BlackRook, newSquare(63), newSquare(7), true},
		{"Q 19-1", board.WhiteQueen, newSquare(19), newSquare(1), true},
		{"q 18-28", board.BlackQueen, newSquare(18), newSquare(28), false},
		{"K 20-30", board.WhiteKing, newSquare(20), newSquare(30), false},
		{"k 60-58", board.BlackKing, newSquare(60), newSquare(58), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := move(tc.piece, tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}

func TestMovePawn(t *testing.T) {
	tests := []struct {
		name  string
		piece board.Piece
		from  square
		to    square
		isOk  bool
	}{
		{"white-up-6-7", board.WhitePawn, newSquare(48), newSquare(56), true},
		{"white-up-1-3", board.WhitePawn, newSquare(8), newSquare(24), true},
		{"white-up-5-7", board.WhitePawn, newSquare(40), newSquare(56), false},
		{"white-up-0-1", board.WhitePawn, newSquare(0), newSquare(8), false},
		{"white-down", board.WhitePawn, newSquare(24), newSquare(16), false},
		{"white-up-right", board.WhitePawn, newSquare(8), newSquare(17), true},
		{"white-up-left", board.WhitePawn, newSquare(35), newSquare(42), true},
		{"white-diagonal-far", board.WhitePawn, newSquare(18), newSquare(36), false},
		{"white-horizontal", board.WhitePawn, newSquare(28), newSquare(29), false},
		{"white-knight", board.WhitePawn, newSquare(18), newSquare(28), false},

		{"black-down-6-5", board.BlackPawn, newSquare(48), newSquare(40), true},
		{"black-down-6-4", board.BlackPawn, newSquare(48), newSquare(32), true},
		{"black-down-3-1", board.BlackPawn, newSquare(24), newSquare(8), false},
		{"black-down-7-6", board.BlackPawn, newSquare(56), newSquare(48), false},
		{"black-up", board.BlackPawn, newSquare(32), newSquare(40), false},
		{"black-down-right", board.BlackPawn, newSquare(18), newSquare(11), true},
		{"black-down-left", board.BlackPawn, newSquare(18), newSquare(9), true},
		{"black-diagonal-far", board.BlackPawn, newSquare(18), newSquare(4), false},
		{"black-horizontal", board.BlackPawn, newSquare(27), newSquare(28), false},
		{"black-knight", board.BlackPawn, newSquare(18), newSquare(12), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := movePawn(tc.piece, tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}

func TestMoveKnight(t *testing.T) {
	tests := []struct {
		name string
		from square
		to   square
		isOk bool
	}{
		{"up-up-right", newSquare(18), newSquare(35), true},
		{"up-right-right", newSquare(18), newSquare(28), true},
		{"down-right-right", newSquare(18), newSquare(12), true},
		{"down-down-right", newSquare(18), newSquare(3), true},
		{"down-down-left", newSquare(18), newSquare(1), true},
		{"down-left-left", newSquare(18), newSquare(8), true},
		{"up-left-left", newSquare(18), newSquare(24), true},
		{"up-up-left", newSquare(18), newSquare(33), true},
		{"vertical", newSquare(2), newSquare(34), false},
		{"horizontal", newSquare(0), newSquare(6), false},
		{"diagonal", newSquare(9), newSquare(36), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := moveKnight(tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}

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
			err := moveBishop(tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatal("want error, got nil")
			}
		})
	}
}

func TestMoveRook(t *testing.T) {
	tests := []struct {
		name string
		from square
		to   square
		isOk bool
	}{
		{"right", newSquare(17), newSquare(22), true},
		{"left", newSquare(17), newSquare(16), true},
		{"up", newSquare(10), newSquare(34), true},
		{"down", newSquare(63), newSquare(7), true},
		{"diagonal", newSquare(9), newSquare(63), false},
		{"knight", newSquare(18), newSquare(28), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := moveRook(tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}

func TestMoveQueen(t *testing.T) {
	tests := []struct {
		name string
		from square
		to   square
		isOk bool
	}{
		{"right", newSquare(17), newSquare(22), true},
		{"left", newSquare(17), newSquare(16), true},
		{"up", newSquare(10), newSquare(34), true},
		{"down", newSquare(63), newSquare(7), true},
		{"up-right", newSquare(9), newSquare(63), true},
		{"down-right", newSquare(9), newSquare(2), true},
		{"up-left", newSquare(6), newSquare(27), true},
		{"down-left", newSquare(19), newSquare(1), true},
		{"knight", newSquare(18), newSquare(28), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := moveQueen(tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}

func TestMoveKing(t *testing.T) {
	tests := []struct {
		name  string
		piece board.Piece
		from  square
		to    square
		isOk  bool
	}{
		{"white-right", board.WhiteKing, newSquare(12), newSquare(13), true},
		{"black-left", board.BlackKing, newSquare(12), newSquare(11), true},
		{"white-up", board.WhiteKing, newSquare(27), newSquare(35), true},
		{"black-down", board.BlackKing, newSquare(61), newSquare(53), true},
		{"white-up-right", board.WhiteKing, newSquare(21), newSquare(30), true},
		{"black-up-left", board.BlackKing, newSquare(34), newSquare(41), true},
		{"white-down-left", board.WhiteKing, newSquare(27), newSquare(18), true},
		{"black-down-right", board.BlackKing, newSquare(37), newSquare(30), true},

		{"white-castling-K", board.WhiteKing, newSquare(4), newSquare(6), true},
		{"white-castling-Q", board.WhiteKing, newSquare(4), newSquare(2), true},
		{"black-castling-K", board.BlackKing, newSquare(60), newSquare(62), true},
		{"black-castling-Q", board.BlackKing, newSquare(60), newSquare(58), true},

		{"white-horizontal-far", board.WhiteKing, newSquare(61), newSquare(63), false},
		{"black-horizontal-far", board.BlackKing, newSquare(28), newSquare(26), false},
		{"white-diagonal-far", board.WhiteKing, newSquare(36), newSquare(54), false},
		{"black-diagonal-far", board.BlackKing, newSquare(29), newSquare(43), false},
		{"white-vertical-far", board.WhiteKing, newSquare(28), newSquare(44), false},
		{"black-vertical-far", board.BlackKing, newSquare(36), newSquare(20), false},
		{"white-knight", board.WhiteKing, newSquare(20), newSquare(30), false},
		{"black-knight", board.BlackKing, newSquare(62), newSquare(52), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := moveKing(tc.piece, tc.from, tc.to)
			if tc.isOk && err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if !tc.isOk && err == nil {
				t.Fatalf("want error, got nil")
			}
		})
	}
}
