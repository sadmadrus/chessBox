package validation

import (
	"testing"

	"github.com/sadmadrus/chessBox/internal/board"
)

func TestAdvancedLogic(t *testing.T) {
	var (
		emptyBrd board.Board

		startBrd1WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3K1NR w KQ - 5 6"
		endBrd1_2WhiteFEN = "rBbq1bnr/pp6/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/R3K1NR b KQ - 5 6"
		endBrd1_3WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P6P/2KR2NR b - - 6 6"
		endBrd1_4WhiteFEN = "rnbq1bnr/ppP5/3Q4/4pBBp/3PPPp1/1P2k1P1/P6P/R3K1NR b KQ - 5 6"
		endBrd1_5WhiteFEN = "rnbq1bnr/ppP5/3p4/4pBBp/3PPPp1/QP2k1P1/P2K3P/R5NR b - - 6 6"

		startBrd1BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PPPp1/QP2k1P1/P6P/R3K1NR b KQ f3 5 6"
		endBrd1_2BlackFEN = "rnbq1bnr/ppP5/3p4/4pB1p/3PP3/QP2kpP1/P6P/R3K1NR w KQ - 5 7"
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
		name     string
		brd      board.Board
		from     square
		to       square
		newpiece board.Piece
		newBoard board.Board
		isValid  bool
		isErr    bool
	}{
		{"no piece", *startBrd1White, newSquare(32), newSquare(33), 0, emptyBrd, false, false},
		{"no promotion, newpiece indicated", *startBrd1White, newSquare(16), newSquare(24), board.WhiteBishop, emptyBrd, false, false},
		{"promotion, wrong color newpiece", *startBrd1White, newSquare(50), newSquare(57), board.BlackBishop, *startBrd1White, false, true},
		{"promotion, no newpiece indicated", *startBrd1White, newSquare(50), newSquare(57), 0, emptyBrd, false, false},
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
			newBoardRes, isValidRes, err := advancedLogic(tc.brd, tc.from, tc.to, tc.newpiece)
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
