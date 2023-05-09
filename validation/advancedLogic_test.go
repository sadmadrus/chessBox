package validation

import (
	"sort"
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestCanMove(t *testing.T) {
	var (
		startBrd1WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3K1NR w KQ - 5 6"
		startBrd1BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PPPp1/QP2k1P1/P6P/R3K1NR b KQ f3 5 6"
		startBrd2BlackFEN = "rn2k2r/8/8/8/3q4/6n1/8/R3K2R b KQq - 5 6"
		startBrd2WhiteFEN = "rn2k2r/8/8/8/3q4/6n1/8/R3K2R w KQq - 5 6"
		invalidBrdFEN     = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3KKNR w KQ - 5 6"
	)

	startBrd1White, _ := board.FromFEN(startBrd1WhiteFEN)
	startBrd1Black, _ := board.FromFEN(startBrd1BlackFEN)
	startBrd2Black, _ := board.FromFEN(startBrd2BlackFEN)
	startBrd2White, _ := board.FromFEN(startBrd2WhiteFEN)
	invalidBrd, _ := board.FromFEN(invalidBrdFEN)

	tests := []struct {
		name      string
		brd       board.Board
		from      board.Square
		to        board.Square
		promoteTo board.Piece
		isErr     bool
	}{
		{"from not valid", *startBrd1White, board.Sq(-1), board.Sq(33), 0, true},
		{"to not valid", *startBrd1White, board.Sq(32), board.Sq(64), 0, true},
		{"from and to equal", *startBrd1White, board.Sq(32), board.Sq(32), 0, true},
		{"crazy promote", *startBrd1White, board.Sq(50), board.Sq(57), 15, true},
		{"invalid board with 2 Kings", *invalidBrd, board.Sq(32), board.Sq(33), 0, true},

		{"no piece", *startBrd1White, board.Sq(32), board.Sq(33), 0, true},
		{"no promotion, promoteTo indicated", *startBrd1White, board.Sq(16), board.Sq(24), board.WhiteBishop, true},
		{"promotion, wrong color promoteTo", *startBrd1White, board.Sq(50), board.Sq(57), board.BlackBishop, true},
		{"promotion, no promoteTo indicated", *startBrd1White, board.Sq(50), board.Sq(57), 0, true},
		{"pawn promotion successful", *startBrd1White, board.Sq(50), board.Sq(57), board.WhiteBishop, false},
		{"white turn, black move", *startBrd1White, board.Sq(57), board.Sq(40), 0, true},
		{"Knight try diagonal move", *startBrd1White, board.Sq(6), board.Sq(13), 0, true},
		{"Q turn, P in the way", *startBrd1White, board.Sq(16), board.Sq(18), 0, true},
		{"B turn, p in the way", *startBrd1White, board.Sq(37), board.Sq(23), 0, true},
		{"P up, clash with p", *startBrd1White, board.Sq(22), board.Sq(30), 0, true},
		{"P up, clash with B", *startBrd1White, board.Sq(29), board.Sq(37), 0, true},
		{"R to P", *startBrd1White, board.Sq(0), board.Sq(8), 0, true},
		{"Q to p", *startBrd1White, board.Sq(16), board.Sq(43), 0, false},
		{"K O-O, N in the way", *startBrd1White, board.Sq(4), board.Sq(6), 0, true},
		{"K O-O-O successful", *startBrd1White, board.Sq(4), board.Sq(2), 0, false},
		{"K too close to k", *startBrd1White, board.Sq(4), board.Sq(11), 0, true},

		{"p g4-f3 successful enPassant", *startBrd1Black, board.Sq(30), board.Sq(21), 0, false},
		{"p g4-h3 enPassant not allowed", *startBrd1Black, board.Sq(30), board.Sq(23), 0, true},
		{"k f3 under self-check", *startBrd1Black, board.Sq(20), board.Sq(21), 0, true},
		{"K O-O through checked cells", *startBrd2White, board.Sq(4), board.Sq(6), 0, true},
		{"K O-O-O through checked cells", *startBrd2White, board.Sq(4), board.Sq(2), 0, true},
		{"k O-O not allowed", *startBrd2White, board.Sq(60), board.Sq(62), 0, true},
		{"k O-O-O through busy cells", *startBrd2Black, board.Sq(60), board.Sq(58), 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CanMove(tc.brd, tc.from, tc.to, tc.promoteTo)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err: %v", err)
			}
			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}
		})
	}
}

func TestAvailableMoves(t *testing.T) {
	var (
		brdWhiteFEN = "3r4/p3PB2/2nr4/2k2Pp1/1b4Pq/1n1Q3P/1p1P1N2/RN2K2R w KQ g6 5 6"
		brdBlackFEN = "3r4/p3PB2/2nr4/2k2Pp1/1b4Pq/1n1Q3P/1p1P1N2/RN2K2R b KQ - 5 6"
		//invalidboardFEN = "dfahdfk 8jerq ~ + dfak"
	)

	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	brdBlack, _ := board.FromFEN(brdBlackFEN)
	//invalidBrd, _ := board.FromFEN(invalidboardFEN)

	tests := []struct {
		name string
		brd  board.Board
		from board.Square
		res  []board.Square
	}{
		//{"invalid board", *invalidBrd, board.Sq("e2"), []board.Square{}},
		{"not exist 65", *brdWhite, board.Sq("65"), []board.Square{}},
		{"empty c1", *brdWhite, board.Sq("c1"), []board.Square{}},

		{"R a1 w", *brdWhite, board.Sq("a1"), []board.Square{board.Sq(8), board.Sq(16), board.Sq(24),
			board.Sq(32), board.Sq(40), board.Sq(48)}},
		{"R a1 b", *brdBlack, board.Sq("a1"), []board.Square{}},

		{"N b1 w", *brdWhite, board.Sq(1), []board.Square{board.Sq(16), board.Sq(18)}},
		{"N b1 b", *brdBlack, board.Sq(1), []board.Square{}},

		{"K e1 w", *brdWhite, board.Sq("e1"), []board.Square{board.Sq(3), board.Sq(5),
			board.Sq(6), board.Sq(12)}},
		{"K e1 b", *brdBlack, board.Sq("e1"), []board.Square{}},

		{"R h1 w", *brdWhite, board.Sq(7), []board.Square{board.Sq(5), board.Sq(6), board.Sq(15)}},
		{"R h1 b", *brdBlack, board.Sq(7), []board.Square{}},

		{"p b2 w", *brdWhite, board.Sq("b2"), []board.Square{}},
		{"p b2 b", *brdBlack, board.Sq("b2"), []board.Square{board.Sq(0)}},

		{"P d2 w", *brdWhite, board.Sq(11), []board.Square{}},
		{"P d2 b", *brdBlack, board.Sq(11), []board.Square{}},

		{"N f2 w", *brdWhite, board.Sq("f2"), []board.Square{}},
		{"N f2 b", *brdBlack, board.Sq("f2"), []board.Square{}},

		{"n b3 w", *brdWhite, board.Sq(17), []board.Square{}},
		{"n b3 b", *brdBlack, board.Sq(17), []board.Square{board.Sq(0), board.Sq(2),
			board.Sq(11), board.Sq(27), board.Sq(32)}},

		{"Q d3 w", *brdWhite, board.Sq("d3"), []board.Square{board.Sq(5), board.Sq(10), board.Sq(12),
			board.Sq(17), board.Sq(18), board.Sq(20), board.Sq(21), board.Sq(22), board.Sq(26),
			board.Sq(27), board.Sq(28), board.Sq(33), board.Sq(35), board.Sq(40), board.Sq(43)}},
		{"Q d3 b", *brdBlack, board.Sq("d3"), []board.Square{}},

		{"P h3 w", *brdWhite, board.Sq(23), []board.Square{}},
		{"P h3 b", *brdBlack, board.Sq(23), []board.Square{}},

		{"b b4 w", *brdWhite, board.Sq("b4"), []board.Square{}},
		{"b b4 b", *brdBlack, board.Sq("b4"), []board.Square{board.Sq(11), board.Sq(16),
			board.Sq(18), board.Sq(32)}},

		{"P g4 w", *brdWhite, board.Sq(30), []board.Square{}},
		{"P g4 b", *brdBlack, board.Sq(30), []board.Square{}},

		{"q h4 w", *brdWhite, board.Sq("h4"), []board.Square{}},
		{"q h4 b", *brdBlack, board.Sq("h4"), []board.Square{board.Sq(13), board.Sq(22),
			board.Sq(23), board.Sq(30), board.Sq(39), board.Sq(47), board.Sq(55), board.Sq(63)}},

		{"k c5 w", *brdWhite, board.Sq(34), []board.Square{}},
		{"k c5 b", *brdBlack, board.Sq(34), []board.Square{board.Sq(41)}},

		{"P f5 w", *brdWhite, board.Sq("f5"), []board.Square{board.Sq(45), board.Sq(46)}},
		{"P f5 b", *brdBlack, board.Sq("f5"), []board.Square{}},

		{"p g5 w", *brdWhite, board.Sq(38), []board.Square{}},
		{"p g5 b", *brdBlack, board.Sq(38), []board.Square{}},

		{"n c6 w", *brdWhite, board.Sq("c6"), []board.Square{}},
		{"n c6 b", *brdBlack, board.Sq("c6"), []board.Square{board.Sq(27), board.Sq(32), board.Sq(36),
			board.Sq(52), board.Sq(57)}},

		{"r d6 w", *brdWhite, board.Sq(43), []board.Square{}},
		{"r d6 b", *brdBlack, board.Sq(43), []board.Square{board.Sq(19), board.Sq(27), board.Sq(35),
			board.Sq(44), board.Sq(45), board.Sq(46), board.Sq(47), board.Sq(51)}},

		{"p a7 w", *brdWhite, board.Sq("a7"), []board.Square{}},
		{"p a7 b", *brdBlack, board.Sq("a7"), []board.Square{board.Sq(32), board.Sq(40)}},

		{"P e7 w", *brdWhite, board.Sq(52), []board.Square{board.Sq(59), board.Sq(60)}},
		{"P e7 b", *brdBlack, board.Sq(52), []board.Square{}},

		{"B f7 w", *brdWhite, board.Sq("f7"), []board.Square{board.Sq(17), board.Sq(26), board.Sq(35),
			board.Sq(39), board.Sq(44), board.Sq(46), board.Sq(60), board.Sq(62)}},
		{"B f7 b", *brdBlack, board.Sq("f7"), []board.Square{}},

		{"r d8 w", *brdWhite, board.Sq(59), []board.Square{}},
		{"r d8 b", *brdBlack, board.Sq(59), []board.Square{board.Sq(51), board.Sq(56), board.Sq(57),
			board.Sq(58), board.Sq(60), board.Sq(61), board.Sq(62), board.Sq(63)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := AvailableMoves(tc.brd, tc.from)

			fail := func() { t.Fatalf("want %d, got %d", tc.res, res) }

			if len(res) != len(tc.res) {
				fail()
			}

			sort.Slice(res, func(i, j int) bool { return int(res[i]) < int(res[j]) })
			sort.Slice(tc.res, func(i, j int) bool { return int(tc.res[i]) < int(tc.res[j]) })

			for i := range res {
				if res[i] != tc.res[i] {
					fail()
				}
			}
		})
	}
}

func TestAdvancedLogic(t *testing.T) {
	var (
		emptyBrd board.Board

		startBrd1WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3K1NR w KQ - 5 6"
		endBrd1_2WhiteFEN = "rBbq1bnr/pp6/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3K1NR b KQ - 0 6"
		endBrd1_3WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/2KR2NR b - - 6 6"
		endBrd1_4WhiteFEN = "rnbq1bnr/ppP5/3Q4/4pBBp/3PPPp1/1P2k1P1/P6P/R3K1NR b KQ - 0 6"
		endBrd1_5WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P2K3P/R5NR b - - 6 6"

		startBrd1BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PPPp1/QP2k1P1/P6P/R3K1NR b KQ f3 5 6"
		endBrd1_2BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PP3/QP2kpP1/P6P/R3K1NR w KQ - 0 7"
		endBrd1_3BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PPPp1/QP3kP1/P6P/R3K1NR w KQ - 6 7"

		startBrd2BlackFEN = "rn2k2r/8/8/8/3q4/6n1/8/R3K2R b KQq - 5 6"
		startBrd2WhiteFEN = "rn2k2r/8/8/8/3q4/6n1/8/R3K2R w KQq - 5 6"
	)

	startBrd1White, _ := board.FromFEN(startBrd1WhiteFEN)
	endBrd1_2White, _ := board.FromFEN(endBrd1_2WhiteFEN)
	endBrd1_3White, _ := board.FromFEN(endBrd1_3WhiteFEN)
	endBrd1_4White, _ := board.FromFEN(endBrd1_4WhiteFEN)
	endBrd1_5White, _ := board.FromFEN(endBrd1_5WhiteFEN)

	startBrd1Black, _ := board.FromFEN(startBrd1BlackFEN)
	endBrd1_2Black, _ := board.FromFEN(endBrd1_2BlackFEN)
	endBrd1_3Black, _ := board.FromFEN(endBrd1_3BlackFEN)

	startBrd2Black, _ := board.FromFEN(startBrd2BlackFEN)
	startBrd2White, _ := board.FromFEN(startBrd2WhiteFEN)

	tests := []struct {
		name      string
		brd       board.Board
		from      square
		to        square
		promoteTo board.Piece
		newBoard  board.Board
		isValid   bool
		isErr     bool
	}{
		{"no piece", *startBrd1White, newSquare(32), newSquare(33), 0, emptyBrd, false, false},
		{"no promotion, promoteTo indicated", *startBrd1White, newSquare(16), newSquare(24), board.WhiteBishop, emptyBrd, false, false},
		{"promotion, wrong color promoteTo", *startBrd1White, newSquare(50), newSquare(57), board.BlackBishop, *startBrd1White, false, true},
		{"promotion, no promoteTo indicated", *startBrd1White, newSquare(50), newSquare(57), 0, emptyBrd, false, false},
		{"pawn promotion successful", *startBrd1White, newSquare(50), newSquare(57), board.WhiteBishop, *endBrd1_2White, true, false},
		{"white turn, black move", *startBrd1White, newSquare(57), newSquare(40), 0, emptyBrd, false, false},
		{"Knight try diagonal move", *startBrd1White, newSquare(6), newSquare(13), 0, emptyBrd, false, false},
		{"Q turn, P in the way", *startBrd1White, newSquare(16), newSquare(18), 0, emptyBrd, false, false},
		{"B turn, p in the way", *startBrd1White, newSquare(37), newSquare(23), 0, emptyBrd, false, false},
		{"P up, clash with p", *startBrd1White, newSquare(22), newSquare(30), 0, emptyBrd, false, false},
		{"P up, clash with B", *startBrd1White, newSquare(29), newSquare(37), 0, emptyBrd, false, false},
		{"R to P", *startBrd1White, newSquare(0), newSquare(8), 0, emptyBrd, false, false},
		{"Q to p", *startBrd1White, newSquare(16), newSquare(43), 0, *endBrd1_4White, true, false},
		{"K O-O, N in the way", *startBrd1White, newSquare(4), newSquare(6), 0, emptyBrd, false, false},
		{"K O-O-O successful", *startBrd1White, newSquare(4), newSquare(2), 0, *endBrd1_3White, true, false},
		{"K too close to k", *startBrd1White, newSquare(4), newSquare(11), 0, *endBrd1_5White, false, false},

		{"p g4-f3 successful enPassant", *startBrd1Black, newSquare(30), newSquare(21), 0, *endBrd1_2Black, true, false},
		{"p g4-h3 enPassant not allowed", *startBrd1Black, newSquare(30), newSquare(23), 0, emptyBrd, false, false},
		{"k f3 under self-check", *startBrd1Black, newSquare(20), newSquare(21), 0, *endBrd1_3Black, false, false},
		{"K O-O through checked cells", *startBrd2White, newSquare(4), newSquare(6), 0, emptyBrd, false, false},
		{"K O-O-O through checked cells", *startBrd2White, newSquare(4), newSquare(2), 0, emptyBrd, false, false},
		{"k O-O not allowed", *startBrd2White, newSquare(60), newSquare(62), 0, emptyBrd, false, false},
		{"k O-O-O through busy cells", *startBrd2Black, newSquare(60), newSquare(58), 0, emptyBrd, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			newBoardRes, isValidRes, err := advancedLogic(tc.brd, tc.from, tc.to, tc.promoteTo)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got err: %v", err)
			}
			if err == nil && tc.isErr {
				t.Fatalf("want err, got nil")
			}
			if newBoardRes != tc.newBoard {
				t.Fatalf("want newBoard %v, got %v", tc.newBoard, newBoardRes)
			}
			if isValidRes != tc.isValid {
				t.Fatalf("want isValid %v, got %v", tc.isValid, isValidRes)
			}
		})
	}
}

func TestCheckPromotion(t *testing.T) {
	tests := []struct {
		name     string
		piece    board.Piece
		to       square
		newpiece board.Piece
		isOk     bool
	}{
		{"p a1 to b", board.BlackPawn, newSquare(0), board.BlackBishop, true},
		{"P h8 to 0", board.WhitePawn, newSquare(63), 0, false},
		{"b b8 to 0", board.BlackBishop, newSquare(57), 0, true},
		{"n c8 to q", board.BlackKnight, newSquare(58), board.BlackQueen, false},
		{"P e4 to N", board.WhitePawn, newSquare(28), board.WhiteKnight, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := checkPawnPromotion(tc.piece, tc.to, tc.newpiece)
			if res != tc.isOk {
				t.Fatalf("want %v, got %v", tc.isOk, res)
			}
		})
	}
}

func TestCheckFigureColor(t *testing.T) {
	tests := []struct {
		name  string
		brd   board.Board
		piece board.Piece
		isOk  bool
	}{
		{"R classical", *board.Classical(), board.WhiteRook, true},
		{"P classical", *board.Classical(), board.WhitePawn, true},
		{"k classical", *board.Classical(), board.BlackKing, false},
		{"b classical", *board.Classical(), board.BlackBishop, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := checkPieceColor(tc.brd, tc.piece)
			if res != tc.isOk {
				t.Fatalf("want %v, got %v", tc.isOk, res)
			}
		})
	}
}

func TestGetSquaresToBePassed(t *testing.T) {
	tests := []struct {
		name   string
		piece  board.Piece
		from   square
		to     square
		result []square
	}{
		{"P e2-e4", board.WhitePawn, newSquare(12), newSquare(28), []square{newSquare(20)}},
		{"p a7-a6", board.BlackPawn, newSquare(48), newSquare(40), []square{}},
		{"N g1-h3", board.WhiteKnight, newSquare(6), newSquare(23), []square{}},
		{"n b8-d7", board.BlackKnight, newSquare(57), newSquare(51), []square{}},
		{"B f1-a6", board.WhiteBishop, newSquare(5), newSquare(40),
			[]square{newSquare(12), newSquare(19), newSquare(26), newSquare(33)}},
		{"b c8-a6", board.BlackBishop, newSquare(58), newSquare(40), []square{newSquare(49)}},
		{"R h1-e1", board.WhiteRook, newSquare(7), newSquare(4),
			[]square{newSquare(6), newSquare(5)}},
		{"r a8-a7", board.BlackRook, newSquare(56), newSquare(48), []square{}},
		{"Q d1-d4", board.WhiteQueen, newSquare(3), newSquare(27),
			[]square{newSquare(11), newSquare(19)}},
		{"q h8-f6", board.BlackQueen, newSquare(63), newSquare(45), []square{newSquare(54)}},
		{"K e5-e6", board.WhiteKing, newSquare(36), newSquare(44), []square{}},
		{"k e8-g8 O-O", board.BlackKing, newSquare(60), newSquare(62), []square{newSquare(61)}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := getSquaresToBePassed(tc.piece, tc.from, tc.to)
			for i := 0; i < len(res); i++ {
				if res[i] != tc.result[i] {
					t.Fatalf("want %v, got %v", tc.result, res)
				}
			}
		})
	}
}

func TestCheckSquaresToBePassed(t *testing.T) {
	tests := []struct {
		name    string
		brd     board.Board
		squares []square
		result  bool
	}{
		{"a3-a6 classical", *board.Classical(),
			[]square{newSquare(16), newSquare(24), newSquare(32), newSquare(40)}, true},
		{"a6-a7 classical", *board.Classical(), []square{newSquare(40), newSquare(48)}, false},
		{"h8 classical", *board.Classical(), []square{newSquare(63)}, false},
		{"h3 classical", *board.Classical(), []square{newSquare(23)}, true},
		{"a1-b2 classical", *board.Classical(), []square{newSquare(0), newSquare(9)}, false},
		{"h4-g5 classical", *board.Classical(), []square{newSquare(31), newSquare(38)}, true},
		{"no squares", *board.Classical(), []square{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := checkSquaresToBePassed(tc.brd, tc.squares)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}

			if res != tc.result {
				t.Fatalf("want %v, got %v", tc.result, res)
			}
		})
	}
}

func TestCheckToSquare(t *testing.T) {
	brdWhiteFEN := "r3k3/1qpb1prp/2np2pb/p3QB2/Pp2P3/1P1PBPnN/2PN2PP/R3K2R w KQq - 5 6"
	brdBlackFEN := "r3n3/1qpb1prp/2np2pb/p3QB2/Pp2P3/1P1PBPPN/2PN2kP/R3K2R b KQq a3 5 6"
	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	brdBlack, _ := board.FromFEN(brdBlackFEN)

	tests := []struct {
		name      string
		brd       *board.Board
		pieceFrom board.Piece
		from      square
		to        square
		isOk      bool
	}{
		{"Q e5-e8 to king", brdWhite, board.WhiteQueen, newSquare(36), newSquare(60), false},
		{"B f5-h3 to Knight", brdWhite, board.WhiteBishop, newSquare(37), newSquare(23), false},
		{"P a4-a5 to pawn", brdWhite, board.WhitePawn, newSquare(24), newSquare(32), false},
		{"P e4-e5 to Queen", brdWhite, board.WhitePawn, newSquare(28), newSquare(36), false},
		{"N h3-g5 to empty", brdWhite, board.WhiteKnight, newSquare(23), newSquare(38), true},
		{"B f5-g6 to pawn", brdWhite, board.WhiteBishop, newSquare(37), newSquare(46), true},

		{"k g2-f2 close to King", brdBlack, board.BlackKing, newSquare(14), newSquare(13), true},
		{"K g2-h1 to Rook", brdBlack, board.BlackKing, newSquare(14), newSquare(7), true},
		{"p b4-a3 to empty", brdBlack, board.BlackPawn, newSquare(25), newSquare(16), true},
		{"n c6-e5 to Queen", brdBlack, board.BlackKnight, newSquare(42), newSquare(36), true},
		{"n c6-b4 to pawn", brdBlack, board.BlackKnight, newSquare(42), newSquare(25), false},
		{"p b4-b3 to Pawn", brdBlack, board.BlackPawn, newSquare(25), newSquare(17), false},

		{"p b4-a3 en Passant", brdBlack, board.BlackPawn, newSquare(25), newSquare(16), true},
		{"p b4-c3 en Passant not allowed", brdBlack, board.BlackPawn, newSquare(25), newSquare(18), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := checkToSquare(tc.brd, tc.pieceFrom, tc.from, tc.to)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isOk {
				t.Fatalf("want %v, got %v", tc.isOk, res)
			}
		})
	}
}

func TestCheckCastling(t *testing.T) {
	brdWhiteFENCasltingsNotAllowed := "r3k3/1qpb1prp/2np2pb/p3QB2/Pp2P3/1P1PBPnN/2PN2PP/R3K2R w Kq - 5 6"
	brdWhiteNotAllowed, _ := board.FromFEN(brdWhiteFENCasltingsNotAllowed)

	tests := []struct {
		name    string
		brd     *board.Board
		piece   board.Piece
		from    square
		to      square
		isValid bool
	}{
		{"K e1-g1 O-O", brdWhiteNotAllowed, board.WhiteKing, newSquare(4), newSquare(6), false},
		{"K e1-c1 O-O-O", brdWhiteNotAllowed, board.WhiteKing, newSquare(4), newSquare(2), false},
		{"k e8-g8 O-O", brdWhiteNotAllowed, board.BlackKing, newSquare(60), newSquare(62), false},
		{"k e8-c8 O-O-O", brdWhiteNotAllowed, board.BlackKing, newSquare(60), newSquare(58), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := checkCastling(tc.brd, tc.piece, tc.from, tc.to)
			if err != nil {
				t.Fatalf("want nil, got err: %v", err)
			}
			if res != tc.isValid {
				t.Fatalf("want %v, got %v", tc.isValid, res)
			}
		})
	}
}

func TestGetNewBoard(t *testing.T) {
	var (
		brdWhiteFEN     = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 5 6"
		newbrdQa7a6FEN  = "4k2r/1Pp4p/Q5pb/p7/P3N3/8/6PP/R3K2R b KQ - 0 6"
		newbrdPb7b8FEN  = "1R2k2r/Q1p4p/r5pb/p7/P3N3/8/6PP/R3K2R b KQ - 0 6"
		newbrdKO_OFEN   = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R4RK1 b - - 6 6"
		newbrdKO_O_OFEN = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/2KR3R b - - 6 6"

		brdBlackFEN    = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R b KQk - 5 6"
		newbrdkO_OFEN  = "5rk1/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 6 7"
		newbrdrh8g8FEN = "4k1r1/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 6 7"
		newbrdbh6g5FEN = "4k2r/QPp4p/r5p1/p5b1/P3N3/8/6PP/R3K2R w KQk - 6 7"
		newbrdpg6g5FEN = "4k2r/QPp4p/r6b/p5p1/P3N3/8/6PP/R3K2R w KQk - 0 7"
	)

	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	newbrdQa7a6, _ := board.FromFEN(newbrdQa7a6FEN)
	newbrdPb7b8, _ := board.FromFEN(newbrdPb7b8FEN)
	newbrdKO_O, _ := board.FromFEN(newbrdKO_OFEN)
	newbrdKO_O_O, _ := board.FromFEN(newbrdKO_O_OFEN)

	brdBlack, _ := board.FromFEN(brdBlackFEN)
	newbrdkO_O, _ := board.FromFEN(newbrdkO_OFEN)
	newbrdrh8g8, _ := board.FromFEN(newbrdrh8g8FEN)
	newbrdbh6g5, _ := board.FromFEN(newbrdbh6g5FEN)
	newbrdpg6g5, _ := board.FromFEN(newbrdpg6g5FEN)

	tests := []struct {
		name     string
		brd      *board.Board
		piece    board.Piece
		from     square
		to       square
		newpiece board.Piece
		newBrd   *board.Board
	}{
		{"Q a7-a6", brdWhite, board.WhiteQueen, newSquare(48), newSquare(40), 0, newbrdQa7a6},
		{"P b7-b8 to Rook", brdWhite, board.WhitePawn, newSquare(49), newSquare(57), board.WhiteRook, newbrdPb7b8},
		{"K O-O", brdWhite, board.WhiteKing, newSquare(4), newSquare(6), 0, newbrdKO_O},
		{"K O-O-O", brdWhite, board.WhiteKing, newSquare(4), newSquare(2), 0, newbrdKO_O_O},

		{"k O-O", brdBlack, board.BlackKing, newSquare(60), newSquare(62), 0, newbrdkO_O},
		{"r h8-g8", brdBlack, board.BlackRook, newSquare(63), newSquare(62), 0, newbrdrh8g8},
		{"b h6-g5", brdBlack, board.BlackBishop, newSquare(47), newSquare(38), 0, newbrdbh6g5},
		{"p g6-g5", brdBlack, board.BlackPawn, newSquare(46), newSquare(38), 0, newbrdpg6g5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := getNewBoard(*tc.brd, tc.piece, tc.from, tc.to, tc.newpiece)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res.FEN() != tc.newBrd.FEN() {
				t.Fatalf("want %v, got %v", tc.newBrd.FEN(), res.FEN())
			}
		})
	}
}

func TestGetSquareByPiece(t *testing.T) {
	var (
		brdFEN  = "1nbk2b1/6P1/8/2nN4/8/3N2Q1/2K5/8 w KQkq - 5 6"
		brd1FEN = "8/5P2/1QN1k3/8/1b1q2R1/2np3B/2r5/4K3 w - - 5 6"
		brd2FEN = "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 w - - 5 6"
	)
	brd, _ := board.FromFEN(brdFEN)
	brd1, _ := board.FromFEN(brd1FEN)
	brd2, _ := board.FromFEN(brd2FEN)

	tests := []struct {
		name        string
		brd         board.Board
		pieceString string
		pieceSquare square
		isErr       bool
	}{
		{"n b8(57)", *brd, "n", newSquare(57), false},
		{"S", *brd, "S", newSquare(0), true},
		{"b c8(58)", *brd, "b", newSquare(58), false},
		{"q", *brd, "q", newSquare(0), true},
		{"k d8(59)", *brd, "k", newSquare(59), false},
		{"N d5(35)", *brd, "N", newSquare(35), false},
		{"K c2(10)", *brd, "K", newSquare(10), false},
		{"Q g3(22)", *brd, "Q", newSquare(22), false},
		{"-", *brd, "-", newSquare(0), true},

		{"k", *brd1, "k", newSquare(44), false},
		{"K", *brd1, "K", newSquare(4), false},
		{"k", *brd2, "k", newSquare(44), false},
		{"B", *brd2, "B", newSquare(23), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := getSquareByPiece(tc.brd, tc.pieceString)
			if err != nil && !tc.isErr {
				t.Fatalf("want nil, got error: %s", err)
			}

			if err == nil && tc.isErr {
				t.Fatal("want error, got nil")
			}

			if res != tc.pieceSquare {
				t.Fatalf("want %v, got %v", tc.pieceSquare, res)
			}
		})
	}
}

func TestIsSquareChecked(t *testing.T) {
	brdWhiteFEN := "r3k2r/7P/2N5/7B/5b2/1q1Qp1n1/2P5/R3K2R w - - 5 6"
	brdWhite, _ := board.FromFEN(brdWhiteFEN)

	tests := []struct {
		name        string
		brd         board.Board
		sq          int
		playerWhite bool
		want        bool
	}{
		{"b1 by q", *brdWhite, 1, true, true},
		{"c1 none", *brdWhite, 2, true, false},
		{"d1 none", *brdWhite, 3, true, false},
		{"f1 by n", *brdWhite, 5, true, true},
		{"g1 none", *brdWhite, 6, true, false},
		{"e5 by b", *brdWhite, 36, true, true},
		{"d4 none", *brdWhite, 27, true, false},
		{"d2 by p", *brdWhite, 11, true, true},

		{"b8 by N", *brdWhite, 57, false, true},
		{"c8 none", *brdWhite, 58, false, false},
		{"d8 by N", *brdWhite, 59, false, true},
		{"e8 by B", *brdWhite, 60, false, true},
		{"f8 none", *brdWhite, 61, false, false},
		{"g8 by P", *brdWhite, 62, false, true},
		{"a8 by R", *brdWhite, 24, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := isSquareChecked(tc.brd, board.Sq(tc.sq), tc.playerWhite)
			if res != tc.want {
				t.Fatalf("want %v, got %v", tc.want, res)
			}
		})
	}
}

func TestGetAvailableMoves(t *testing.T) {
	var (
		// Одна и та же позиция проверяется для каждой фигуры черных и белых
		brdWhiteFEN = "3r4/p3PB2/2nr4/2k2Pp1/1b4Pq/1n1Q3P/1p1P1N2/RN2K2R w KQ g6 5 6"
		brdBlackFEN = "3r4/p3PB2/2nr4/2k2Pp1/1b4Pq/1n1Q3P/1p1P1N2/RN2K2R b KQ g6 5 6"
	)

	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	brdBlack, _ := board.FromFEN(brdBlackFEN)

	tests := []struct {
		name  string
		brd   board.Board
		from  square
		res   []square
		isErr bool
	}{
		{"not exist 65", *brdWhite, newSquare(65), []square{}, true},
		{"empty c1", *brdWhite, newSquare(2), []square{}, false},

		{"R a1 w", *brdWhite, newSquare(0), []square{newSquare(8), newSquare(16), newSquare(24),
			newSquare(32), newSquare(40), newSquare(48)}, false},
		{"R a1 b", *brdBlack, newSquare(0), []square{}, false},

		{"N b1 w", *brdWhite, newSquare(1), []square{newSquare(16), newSquare(18)}, false},
		{"N b1 b", *brdBlack, newSquare(1), []square{}, false},

		{"K e1 w", *brdWhite, newSquare(4), []square{newSquare(3), newSquare(5),
			newSquare(6), newSquare(12)}, false},
		{"K e1 b", *brdBlack, newSquare(4), []square{}, false},

		{"R h1 w", *brdWhite, newSquare(7), []square{newSquare(5), newSquare(6), newSquare(15)}, false},
		{"R h1 b", *brdBlack, newSquare(7), []square{}, false},

		{"p b2 w", *brdWhite, newSquare(9), []square{}, false},
		{"p b2 b", *brdBlack, newSquare(9), []square{newSquare(0)}, false},

		{"P d2 w", *brdWhite, newSquare(11), []square{}, false},
		{"P d2 b", *brdBlack, newSquare(11), []square{}, false},

		{"N f2 w", *brdWhite, newSquare(13), []square{}, false},
		{"N f2 b", *brdBlack, newSquare(13), []square{}, false},

		{"n b3 w", *brdWhite, newSquare(17), []square{}, false},
		{"n b3 b", *brdBlack, newSquare(17), []square{newSquare(0), newSquare(2), newSquare(11),
			newSquare(27), newSquare(32)}, false},

		{"Q d3 w", *brdWhite, newSquare(19), []square{newSquare(17), newSquare(18), newSquare(20),
			newSquare(21), newSquare(22), newSquare(27), newSquare(35), newSquare(43), newSquare(10),
			newSquare(12), newSquare(5), newSquare(26), newSquare(33), newSquare(40), newSquare(28)}, false},
		{"Q d3 b", *brdBlack, newSquare(19), []square{}, false},

		{"P h3 w", *brdWhite, newSquare(23), []square{}, false},
		{"P h3 b", *brdBlack, newSquare(23), []square{}, false},

		{"b b4 w", *brdWhite, newSquare(25), []square{}, false},
		{"b b4 b", *brdBlack, newSquare(25), []square{newSquare(16), newSquare(18), newSquare(11),
			newSquare(32)}, false},

		{"P g4 w", *brdWhite, newSquare(30), []square{}, false},
		{"P g4 b", *brdBlack, newSquare(30), []square{}, false},

		{"q h4 w", *brdWhite, newSquare(31), []square{}, false},
		{"q h4 b", *brdBlack, newSquare(31), []square{newSquare(23), newSquare(39), newSquare(47),
			newSquare(55), newSquare(63), newSquare(30), newSquare(22), newSquare(13)}, false},

		{"k c5 w", *brdWhite, newSquare(34), []square{}, false},
		{"k c5 b", *brdBlack, newSquare(34), []square{newSquare(41)}, false},

		{"P f5 w", *brdWhite, newSquare(37), []square{newSquare(45), newSquare(46)}, false},
		{"P f5 b", *brdBlack, newSquare(37), []square{}, false},

		{"p g5 w", *brdWhite, newSquare(38), []square{}, false},
		{"p g5 b", *brdBlack, newSquare(38), []square{}, false},

		{"n c6 w", *brdWhite, newSquare(42), []square{}, false},
		{"n c6 b", *brdBlack, newSquare(42), []square{newSquare(32), newSquare(27),
			newSquare(36), newSquare(52), newSquare(57)}, false},

		{"r d6 w", *brdWhite, newSquare(43), []square{}, false},
		{"r d6 b", *brdBlack, newSquare(43), []square{newSquare(19), newSquare(27), newSquare(35),
			newSquare(51), newSquare(44), newSquare(45), newSquare(46), newSquare(47)}, false},

		{"p a7 w", *brdWhite, newSquare(48), []square{}, false},
		{"p a7 b", *brdBlack, newSquare(48), []square{newSquare(40), newSquare(32)}, false},

		{"P e7 w", *brdWhite, newSquare(52), []square{newSquare(59), newSquare(60)}, false},
		{"P e7 b", *brdBlack, newSquare(52), []square{}, false},

		{"B f7 w", *brdWhite, newSquare(53), []square{newSquare(60), newSquare(62), newSquare(46),
			newSquare(39), newSquare(44), newSquare(35), newSquare(26), newSquare(17)}, false},
		{"B f7 b", *brdBlack, newSquare(53), []square{}, false},

		{"r d8 w", *brdWhite, newSquare(59), []square{}, false},
		{"r d8 b", *brdBlack, newSquare(59), []square{newSquare(51), newSquare(56), newSquare(57),
			newSquare(58), newSquare(60), newSquare(61), newSquare(62), newSquare(63)}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := getAvailableMoves(tc.brd, tc.from)
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
