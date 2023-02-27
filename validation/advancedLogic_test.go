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

func TestIsKingChecked(t *testing.T) {
	var (
		brd1FEN = "8/5P2/1QN1k3/8/1b1q2R1/2np3B/2r5/4K3 w - - 5 6"
		brd2FEN = "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 w - - 5 6"
		brd3FEN = "8/5P2/1Q2k3/8/1b1q2R1/2np3B/8/2r1K3 w - - 5 6"
		brd4FEN = "8/2N5/4k3/8/3K4/8/2n5/8 w - - 5 6"
		brd5FEN = "2B5/1P6/k7/8/6p1/7K/8/8 w - - 5 6"
	)
	brd1, _ := board.FromFEN(brd1FEN)
	brd2, _ := board.FromFEN(brd2FEN)
	brd3, _ := board.FromFEN(brd3FEN)
	brd4, _ := board.FromFEN(brd4FEN)
	brd5, _ := board.FromFEN(brd5FEN)

	tests := []struct {
		name      string
		brd       board.Board
		king      board.Piece
		isChecked bool
	}{
		{"k none", *brd1, board.BlackKing, false},
		{"K none", *brd1, board.WhiteKing, false},
		{"k by B", *brd2, board.BlackKing, true},
		{"K by b", *brd2, board.WhiteKing, true},
		{"k by Q", *brd3, board.BlackKing, true},
		{"K by r", *brd3, board.WhiteKing, true},
		{"k by N", *brd4, board.BlackKing, true},
		{"K by n", *brd4, board.WhiteKing, true},
		{"k none", *brd5, board.BlackKing, false},
		{"K by p", *brd5, board.WhiteKing, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := isKingChecked(tc.brd, tc.king)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isChecked {
				t.Fatalf("want %v, got %v", tc.isChecked, res)
			}
		})
	}
}

func TestIsSquareChecked(t *testing.T) {
	brdWhiteFEN := "r3k2r/7P/2N5/7B/5b2/1q1Qp1n1/2P5/R3K2R w - - 5 6"
	brdWhite, _ := board.FromFEN(brdWhiteFEN)

	tests := []struct {
		name      string
		brd       board.Board
		sq        square
		king      board.Piece
		isChecked bool
	}{
		{"b1 by q", *brdWhite, newSquare(1), board.WhiteKing, true},
		{"c1 none", *brdWhite, newSquare(2), board.WhiteKing, false},
		{"d1 none", *brdWhite, newSquare(3), board.WhiteKing, false},
		{"f1 by n", *brdWhite, newSquare(5), board.WhiteKing, true},
		{"g1 none", *brdWhite, newSquare(6), board.WhiteKing, false},
		{"e5 by b", *brdWhite, newSquare(36), board.WhiteKing, true},
		{"d4 none", *brdWhite, newSquare(27), board.WhiteKing, false},
		{"d2 by p", *brdWhite, newSquare(11), board.WhiteKing, true},

		{"b8 by N", *brdWhite, newSquare(57), board.BlackKing, true},
		{"c8 none", *brdWhite, newSquare(58), board.BlackKing, false},
		{"d8 by N", *brdWhite, newSquare(59), board.BlackKing, true},
		{"e8 by B", *brdWhite, newSquare(60), board.BlackKing, true},
		{"f8 none", *brdWhite, newSquare(61), board.BlackKing, false},
		{"g8 by P", *brdWhite, newSquare(62), board.BlackKing, true},
		{"a8 by R", *brdWhite, newSquare(24), board.BlackKing, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := isSquareChecked(tc.brd, tc.sq, tc.king)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isChecked {
				t.Fatalf("want %v, got %v", tc.isChecked, res)
			}
		})
	}
}

func TestIsSquareCheckedByKnights(t *testing.T) {
	brdWhiteFEN := "8/8/8/8/2b5/1k6/3K4/1r3N2 w - - 5 6"
	brdBlackFEN := "k7/8/1NK5/8/8/8/8/8 b - - 5 6"

	brdWhite, _ := board.FromFEN(brdWhiteFEN)
	brdBlack, _ := board.FromFEN(brdBlackFEN)

	tests := []struct {
		name        string
		brd         board.Board
		kingSquare  square
		enemyKnight board.Piece
		isChecked   bool
	}{
		{"no black knights for white king", *brdWhite, newSquare(11), board.BlackKnight, false},
		{"white knight for black king", *brdBlack, newSquare(56), board.WhiteKnight, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := isSquareCheckedByKnights(tc.brd, tc.kingSquare, tc.enemyKnight)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isChecked {
				t.Fatalf("want %v, got %v", tc.isChecked, res)
			}
		})
	}
}

func TestIsSquareCheckedVerticallyOrHorizontally(t *testing.T) {
	var (
		brdWhite1FEN = "3q4/8/8/8/8/8/3K4/8 w - - 5 6"
		brdWhite2FEN = "3q4/8/3p4/3r4/3b4/8/3K4/8 w - - 5 6"
		brdWhite3FEN = "8/8/8/8/8/8/3K1RQq/8 w - - 5 6"
		brdBlack1FEN = "8/8/8/8/8/1q1Q2k1/8/7K b - - 5 6"
		brdBlack2FEN = "8/8/8/8/8/6kR/8/7K b - - 5 6"
		brdBlack3FEN = "6k1/8/8/6N1/8/8/8/6R1 b - - 5 6"
	)

	brdWhite1, _ := board.FromFEN(brdWhite1FEN)
	brdWhite2, _ := board.FromFEN(brdWhite2FEN)
	brdWhite3, _ := board.FromFEN(brdWhite3FEN)
	brdBlack1, _ := board.FromFEN(brdBlack1FEN)
	brdBlack2, _ := board.FromFEN(brdBlack2FEN)
	brdBlack3, _ := board.FromFEN(brdBlack3FEN)

	tests := []struct {
		name       string
		brd        board.Board
		kingSquare square
		enemyRook  board.Piece
		enemyQueen board.Piece
		isChecked  bool
	}{
		{"q far vertically", *brdWhite1, newSquare(11), board.BlackRook, board.BlackQueen, true},
		{"q, r vertically hidden by b", *brdWhite2, newSquare(11), board.BlackRook, board.BlackQueen, false},
		{"q horizontally hidden by R, Q", *brdWhite3, newSquare(11), board.BlackRook, board.BlackQueen, false},

		{"Q horizontally", *brdBlack1, newSquare(22), board.WhiteRook, board.WhiteQueen, true},
		{"R close horizontally", *brdBlack2, newSquare(22), board.WhiteRook, board.WhiteQueen, true},
		{"R vertically hidden by N, b", *brdBlack3, newSquare(62), board.WhiteRook, board.WhiteQueen, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := isSquareCheckedVerticallyOrHorizontally(tc.brd, tc.kingSquare, tc.enemyRook, tc.enemyQueen)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isChecked {
				t.Fatalf("want %v, got %v", tc.isChecked, res)
			}
		})
	}
}

func TestIsSquareCheckedDiagonally(t *testing.T) {
	var (
		brdWhite1FEN = "7q/8/8/8/8/8/8/K7 w - - 5 6"
		brdWhite2FEN = "8/1N6/8/3K4/8/1b6/8/7Q w - - 5 6"
		brdWhite3FEN = "8/8/3K4/2p1p3/8/8/8/8 w - - 5 6"
		brdBlack1FEN = "8/8/8/3r3q/8/5k2/4P3/8 b - - 5 6"
		brdBlack2FEN = "8/1B6/2P5/5P2/4k3/5b2/2R5/1Q6 b - - 5 6"
		brdBlack3FEN = "8/5P2/1QN1k3/8/1b1q4/3p3B/2r5/4K3 b - - 5 6"
	)

	brdWhite1, _ := board.FromFEN(brdWhite1FEN)
	brdWhite2, _ := board.FromFEN(brdWhite2FEN)
	brdWhite3, _ := board.FromFEN(brdWhite3FEN)
	brdBlack1, _ := board.FromFEN(brdBlack1FEN)
	brdBlack2, _ := board.FromFEN(brdBlack2FEN)
	brdBlack3, _ := board.FromFEN(brdBlack3FEN)

	tests := []struct {
		name        string
		brd         board.Board
		kingSquare  square
		enemyQueen  board.Piece
		enemyBishop board.Piece
		enemyPawn   board.Piece
		isChecked   bool
	}{
		{"q far up-right", *brdWhite1, newSquare(0), board.BlackQueen, board.BlackBishop, board.BlackPawn, true},
		{"b down-left", *brdWhite2, newSquare(35), board.BlackQueen, board.BlackBishop, board.BlackPawn, true},
		{"p down-left, down-right", *brdWhite3, newSquare(35), board.BlackQueen, board.BlackBishop, board.BlackPawn, false},

		{"P down-left", *brdBlack1, newSquare(21), board.WhiteQueen, board.WhiteBishop, board.WhitePawn, true},
		{"B, Q hidden by P, R; P ud-right", *brdBlack2, newSquare(28), board.WhiteQueen, board.WhiteBishop, board.WhitePawn, false},
		{"B down right", *brdBlack3, newSquare(44), board.WhiteQueen, board.WhiteBishop, board.WhitePawn, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := isSquareCheckedDiagonally(tc.brd, tc.kingSquare, tc.enemyQueen, tc.enemyBishop, tc.enemyPawn)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isChecked {
				t.Fatalf("want %v, got %v", tc.isChecked, res)
			}
		})
	}
}

func TestCheckDistanceToEnemyKing(t *testing.T) {
	var (
		brd1FEN = "8/8/8/8/8/3k4/1K6/8 w - - 5 6"
		brd2FEN = "8/8/8/8/8/2k5/1K6/8 w - - 5 6"
		brd3FEN = "8/8/8/8/8/8/1K6/1k6 w - - 5 6"
	)

	brd1, _ := board.FromFEN(brd1FEN)
	brd2, _ := board.FromFEN(brd2FEN)
	brd3, _ := board.FromFEN(brd3FEN)

	tests := []struct {
		name                string
		brd                 board.Board
		isEnemyKingAdjacent bool
	}{
		{"knight position", *brd1, false},
		{"up-right", *brd2, true},
		{"down", *brd3, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := checkDistanceToEnemyKing(tc.brd)
			if err != nil {
				t.Fatalf("want nil, got error: %s", err)
			}
			if res != tc.isEnemyKingAdjacent {
				t.Fatalf("want %v, got %v", tc.isEnemyKingAdjacent, res)
			}
		})
	}
}
