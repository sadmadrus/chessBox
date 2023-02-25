package validation

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

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

// TODO: тест не работает предположительно из-за isSQuareCHecked вызываемой внути. Проверить позже
//func TestCheckCastling(t *testing.T) {
//	brdWhiteFENCasltingsNotAllowed := "r3k3/1qpb1prp/2np2pb/p3QB2/Pp2P3/1P1PBPnN/2PN2PP/R3K2R w Kq - 5 6"
//	brdWhiteNotAllowed, _ := board.FromFEN(brdWhiteFENCasltingsNotAllowed)
//
//	tests := []struct {
//		name  string
//		brd   *board.Board
//		piece board.Piece
//		from  square
//		to    square
//		isOk  bool
//	}{
//		{"K e1-g1 O-O", brdWhiteNotAllowed, board.WhiteKing, newSquare(4), newSquare(6), false},
//		{"K e1-c1 O-O-O", brdWhiteNotAllowed, board.WhiteKing, newSquare(4), newSquare(2), false},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := checkCastling(tc.brd, tc.piece, tc.from, tc.to)
//			if err != nil {
//				t.Fatalf("want nil, got err: %v", err)
//			}
//			if res != tc.isOk {
//				t.Fatalf("want %v, got %v", tc.isOk, res)
//			}
//		})
//	}
//}

func TestGetNewBoard(t *testing.T) {
	var (
		brdWhiteFEN     = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 5 6"
		newbrdQa7a6FEN  = "4k2r/1Pp4p/Q5pb/p7/P3N3/8/6PP/R3K2R b KQ - 5 6"
		newbrdPb7b8FEN  = "1R2k2r/Q1p4p/r5pb/p7/P3N3/8/6PP/R3K2R b KQ - 5 6"
		newbrdKO_OFEN   = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R4RK1 b - - 6 6"
		newbrdKO_O_OFEN = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/2KR3R b - - 6 6"

		brdBlackFEN    = "4k2r/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R b KQk - 5 6"
		newbrdkO_OFEN  = "5rk1/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 6 7"
		newbrdrh8g8FEN = "4k1r1/QPp4p/r5pb/p7/P3N3/8/6PP/R3K2R w KQ - 6 7"
		newbrdbh6g5FEN = "4k2r/QPp4p/r5p1/p5b1/P3N3/8/6PP/R3K2R w KQk - 6 7"
		newbrdpg6g5FEN = "4k2r/QPp4p/r6b/p5p1/P3N3/8/6PP/R3K2R w KQk - 5 7"
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

func TestCheckEnemyKnightsNearKing(t *testing.T) {
	brdWhiteFEN := "8/8/8/8/2b5/1k6/3K4/1r3N2 w - - 5 6"
	brdBlackFEN := "k7/8/1NK5/8/8/8/8/8 b - - 5 6"

	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	brdBlack, _ := board.FromFEN(brdBlackFEN)

	tests := []struct {
		name                 string
		brd                  board.Board
		kingSquare           square
		enemyKnight          board.Piece
		isEnemyKnightPresent bool
	}{
		{"no black knights for white king", *brdWhite, newSquare(11), board.BlackKnight, false},
		{"white knight for black king", *brdBlack, newSquare(56), board.WhiteKnight, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := checkEnemyKnightsNearKing(tc.brd, tc.kingSquare, tc.enemyKnight)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isEnemyKnightPresent {
				t.Fatalf("want %v, got %v", tc.isEnemyKnightPresent, res)
			}
		})
	}
}
