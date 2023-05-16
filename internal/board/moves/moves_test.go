package moves

import (
	"sort"
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

func TestGetMoves(t *testing.T) {
	brd1FEN := "3r4/p3PB2/2nr4/2k2Pp1/1b4Pq/1n1Q3P/1p1P1N2/RN2K2R w KQ g6 5 6"
	brd1, _ := board.FromFEN(brd1FEN)

	tests := []struct {
		name  string
		brd   board.Board
		from  square
		res   []square
		isErr bool
	}{
		{"not exist 65", *brd1, newSquare(65), []square{}, true},
		{"empty c1", *brd1, newSquare(2), []square{}, false},
		{"R a1", *brd1, newSquare(0), []square{newSquare(8), newSquare(16), newSquare(24),
			newSquare(32), newSquare(40), newSquare(48), newSquare(56), newSquare(1), newSquare(2),
			newSquare(3), newSquare(4), newSquare(5), newSquare(6), newSquare(7)}, false},
		{"N b1", *brd1, newSquare(1), []square{newSquare(16), newSquare(18), newSquare(11)}, false},
		{"K e1", *brd1, newSquare(4), []square{newSquare(2), newSquare(3), newSquare(5),
			newSquare(6), newSquare(11), newSquare(12), newSquare(13)}, false},
		{"R h1", *brd1, newSquare(7), []square{newSquare(0), newSquare(1), newSquare(2),
			newSquare(3), newSquare(4), newSquare(5), newSquare(6), newSquare(15), newSquare(23),
			newSquare(31), newSquare(39), newSquare(47), newSquare(55), newSquare(63)}, false},
		{"p b2", *brd1, newSquare(9), []square{newSquare(0), newSquare(1), newSquare(2)}, false},
		{"P d2", *brd1, newSquare(11), []square{newSquare(18), newSquare(19), newSquare(20),
			newSquare(27)}, false},
		{"N f2", *brd1, newSquare(13), []square{newSquare(3), newSquare(7), newSquare(19),
			newSquare(28), newSquare(30), newSquare(23)}, false},
		{"n b3", *brd1, newSquare(17), []square{newSquare(0), newSquare(2), newSquare(11),
			newSquare(27), newSquare(34), newSquare(32)}, false},
		{"Q d3", *brd1, newSquare(19), []square{newSquare(3), newSquare(11), newSquare(16),
			newSquare(17), newSquare(18), newSquare(20), newSquare(21), newSquare(22), newSquare(23),
			newSquare(27), newSquare(35), newSquare(43), newSquare(51), newSquare(59), newSquare(10),
			newSquare(1), newSquare(12), newSquare(5), newSquare(26), newSquare(33), newSquare(40),
			newSquare(28), newSquare(37), newSquare(46), newSquare(55)}, false},
		{"P h3", *brd1, newSquare(23), []square{newSquare(30), newSquare(31)}, false},
		{"b b4", *brd1, newSquare(25), []square{newSquare(16), newSquare(18), newSquare(11),
			newSquare(4), newSquare(32), newSquare(34), newSquare(43), newSquare(52), newSquare(61)}, false},
		{"P g4", *brd1, newSquare(30), []square{newSquare(37), newSquare(38), newSquare(39)}, false},
		{"q h4", *brd1, newSquare(31), []square{newSquare(7), newSquare(15), newSquare(23),
			newSquare(39), newSquare(47), newSquare(55), newSquare(63), newSquare(24), newSquare(25),
			newSquare(26), newSquare(27), newSquare(28), newSquare(29), newSquare(30), newSquare(22),
			newSquare(13), newSquare(4), newSquare(38), newSquare(45), newSquare(52), newSquare(59)}, false},
		{"k c5", *brd1, newSquare(34), []square{newSquare(25), newSquare(26), newSquare(27),
			newSquare(33), newSquare(35), newSquare(41), newSquare(42), newSquare(43)}, false},
		{"P f5", *brd1, newSquare(37), []square{newSquare(44), newSquare(45), newSquare(46)}, false},
		{"p g5", *brd1, newSquare(38), []square{newSquare(29), newSquare(30), newSquare(31)}, false},
		{"n c6", *brd1, newSquare(42), []square{newSquare(32), newSquare(25), newSquare(27),
			newSquare(36), newSquare(52), newSquare(59), newSquare(57), newSquare(48)}, false},
		{"r d6", *brd1, newSquare(43), []square{newSquare(3), newSquare(11), newSquare(19),
			newSquare(27), newSquare(35), newSquare(51), newSquare(59), newSquare(40), newSquare(41),
			newSquare(42), newSquare(44), newSquare(45), newSquare(46), newSquare(47)}, false},
		{"p a7", *brd1, newSquare(48), []square{newSquare(40), newSquare(41), newSquare(32)}, false},
		{"P e7", *brd1, newSquare(52), []square{newSquare(59), newSquare(60), newSquare(61)}, false},
		{"B f7", *brd1, newSquare(53), []square{newSquare(60), newSquare(62), newSquare(46),
			newSquare(39), newSquare(44), newSquare(35), newSquare(26), newSquare(17), newSquare(8)}, false},
		{"r d8", *brd1, newSquare(59), []square{newSquare(3), newSquare(11), newSquare(19),
			newSquare(27), newSquare(35), newSquare(43), newSquare(51), newSquare(56), newSquare(57),
			newSquare(58), newSquare(60), newSquare(61), newSquare(62), newSquare(63)}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := getMoves(tc.brd, tc.from)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err")
			}
			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}

			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })

			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetPawnMoves(t *testing.T) {
	tests := []struct {
		name  string
		piece board.Piece
		from  square
		res   []square
	}{
		{"P e2", board.WhitePawn, newSquare(12), []square{newSquare(19), newSquare(20), newSquare(21), newSquare(28)}},
		{"P a7", board.WhitePawn, newSquare(48), []square{newSquare(56), newSquare(57)}},
		{"P h5", board.WhitePawn, newSquare(39), []square{newSquare(46), newSquare(47)}},
		{"P e8", board.WhitePawn, newSquare(60), []square{}},
		{"P a1", board.WhitePawn, newSquare(0), []square{}},

		{"p a7", board.BlackPawn, newSquare(48), []square{newSquare(40), newSquare(41), newSquare(32)}},
		{"p e2", board.BlackPawn, newSquare(12), []square{newSquare(3), newSquare(4), newSquare(5)}},
		{"p h2", board.BlackPawn, newSquare(15), []square{newSquare(6), newSquare(7)}},
		{"p d1", board.BlackPawn, newSquare(3), []square{}},
		{"p h8", board.BlackPawn, newSquare(63), []square{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getPawnMoves(tc.piece, tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })

			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetKnightMoves(t *testing.T) {
	tests := []struct {
		name string
		from square
		res  []square
	}{
		{"c3", newSquare(18), []square{newSquare(35), newSquare(28), newSquare(12), newSquare(3), newSquare(1), newSquare(8), newSquare(24), newSquare(33)}},
		{"a1", newSquare(0), []square{newSquare(17), newSquare(10)}},
		{"h1", newSquare(7), []square{newSquare(13), newSquare(22)}},
		{"a8", newSquare(56), []square{newSquare(41), newSquare(50)}},
		{"h8", newSquare(63), []square{newSquare(53), newSquare(46)}},
		{"e1", newSquare(4), []square{newSquare(10), newSquare(19), newSquare(21), newSquare(14)}},
		{"h6", newSquare(47), []square{newSquare(30), newSquare(37), newSquare(53), newSquare(62)}},
		{"b8", newSquare(57), []square{newSquare(40), newSquare(42), newSquare(51)}},
		{"a2", newSquare(8), []square{newSquare(25), newSquare(18), newSquare(2)}},
		{"g7", newSquare(54), []square{newSquare(60), newSquare(44), newSquare(37), newSquare(39)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getKnightMoves(tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })

			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetBishopMoves(t *testing.T) {
	tests := []struct {
		name string
		from square
		res  []square
	}{
		{"a1", newSquare(0), []square{newSquare(9), newSquare(18), newSquare(27), newSquare(36),
			newSquare(45), newSquare(54), newSquare(63)}},
		{"e8", newSquare(60), []square{newSquare(51), newSquare(42), newSquare(33), newSquare(24),
			newSquare(53), newSquare(46), newSquare(39)}},
		{"e4", newSquare(28), []square{newSquare(35), newSquare(42), newSquare(49), newSquare(56),
			newSquare(19), newSquare(10), newSquare(1), newSquare(37), newSquare(46), newSquare(55),
			newSquare(21), newSquare(14), newSquare(7)}},
		{"h2", newSquare(15), []square{newSquare(6), newSquare(22), newSquare(29), newSquare(36),
			newSquare(43), newSquare(50), newSquare(57)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getBishopMoves(tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })
			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetRookMoves(t *testing.T) {
	tests := []struct {
		name string
		from square
		res  []square
	}{
		{"a1", newSquare(0), []square{newSquare(8), newSquare(16), newSquare(24), newSquare(32),
			newSquare(40), newSquare(48), newSquare(56), newSquare(1), newSquare(2), newSquare(3),
			newSquare(4), newSquare(5), newSquare(6), newSquare(7)}},
		{"e8", newSquare(60), []square{newSquare(52), newSquare(44), newSquare(36), newSquare(28),
			newSquare(20), newSquare(12), newSquare(4), newSquare(56), newSquare(57), newSquare(58),
			newSquare(59), newSquare(61), newSquare(62), newSquare(63)}},
		{"c3", newSquare(18), []square{newSquare(16), newSquare(17), newSquare(19), newSquare(20),
			newSquare(21), newSquare(22), newSquare(23), newSquare(10), newSquare(2), newSquare(26),
			newSquare(34), newSquare(42), newSquare(50), newSquare(58)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getRookMoves(tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })
			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetQueenMoves(t *testing.T) {
	tests := []struct {
		name string
		from square
		res  []square
	}{
		{"a1", newSquare(0), []square{newSquare(8), newSquare(16), newSquare(24), newSquare(32),
			newSquare(40), newSquare(48), newSquare(56), newSquare(1), newSquare(2), newSquare(3),
			newSquare(4), newSquare(5), newSquare(6), newSquare(7), newSquare(9), newSquare(18),
			newSquare(27), newSquare(36), newSquare(45), newSquare(54), newSquare(63)}},
		{"e8", newSquare(60), []square{newSquare(52), newSquare(44), newSquare(36), newSquare(28),
			newSquare(20), newSquare(12), newSquare(4), newSquare(56), newSquare(57), newSquare(58),
			newSquare(59), newSquare(61), newSquare(62), newSquare(63), newSquare(51), newSquare(42),
			newSquare(33), newSquare(24), newSquare(53), newSquare(46), newSquare(39)}},
		{"c3", newSquare(18), []square{newSquare(16), newSquare(17), newSquare(19), newSquare(20),
			newSquare(21), newSquare(22), newSquare(23), newSquare(10), newSquare(2), newSquare(26),
			newSquare(34), newSquare(42), newSquare(50), newSquare(58), newSquare(0), newSquare(9),
			newSquare(11), newSquare(4), newSquare(25), newSquare(32), newSquare(27), newSquare(36),
			newSquare(45), newSquare(54), newSquare(63)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getQueenMoves(tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })
			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestGetKingMoves(t *testing.T) {
	tests := []struct {
		name  string
		piece board.Piece
		from  square
		res   []square
	}{
		{"K a1", board.WhiteKing, newSquare(0), []square{newSquare(8), newSquare(9), newSquare(1)}},
		{"K e1", board.WhiteKing, newSquare(4), []square{newSquare(2), newSquare(3), newSquare(5),
			newSquare(6), newSquare(11), newSquare(12), newSquare(13)}},
		{"k e1", board.BlackKing, newSquare(4), []square{newSquare(3), newSquare(5),
			newSquare(11), newSquare(12), newSquare(13)}},
		{"k e8", board.BlackKing, newSquare(60), []square{newSquare(58), newSquare(59),
			newSquare(61), newSquare(62), newSquare(51), newSquare(52), newSquare(53)}},
		{"K e8", board.WhiteKing, newSquare(60), []square{newSquare(59), newSquare(61),
			newSquare(51), newSquare(52), newSquare(53)}},
		{"K a6", board.WhiteKing, newSquare(40), []square{newSquare(48), newSquare(49),
			newSquare(41), newSquare(33), newSquare(32)}},
		{"k c3", board.BlackKing, newSquare(18), []square{newSquare(9), newSquare(10),
			newSquare(11), newSquare(17), newSquare(19), newSquare(25), newSquare(26), newSquare(27)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getKingMoves(tc.piece, tc.from)
			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }
			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return res[i].toInt() < res[j].toInt() })
			sort.Slice(tc.res, func(i, j int) bool { return tc.res[i].toInt() < tc.res[j].toInt() })
			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}
