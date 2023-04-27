package game

import (
	"fmt"

	"github.com/sadmadrus/chessBox/internal/board"
)

// fullMove содержит информацию о ходе.
type fullMove struct {
	white, black Move
}

// Move содержит информацию о ходе одного игрока.
type Move interface {
	FromSquare() board.Square
	ToSquare() board.Square
}

type simpleMove struct {
	from, to board.Square
}

// promotion описывает ход с проведением пешки.
type promotion struct {
	simpleMove
	promoteTo board.Piece
}

type castling board.Castling

func (h simpleMove) FromSquare() board.Square {
	return h.from
}

func (h simpleMove) ToSquare() board.Square {
	return h.to
}

func (p promotion) FromSquare() board.Square {
	return p.from
}

func (p promotion) ToSquare() board.Square {
	return p.to
}

func (p promotion) toPiece() board.Piece {
	return p.promoteTo
}

func (c castling) FromSquare() board.Square {
	switch c {
	case castling(board.WhiteKingside), castling(board.WhiteQueenside):
		return board.Sq("e1")
	default:
		return board.Sq("e8")
	}
}

func (c castling) ToSquare() board.Square {
	switch c {
	case castling(board.WhiteKingside):
		return board.Sq("g1")
	case castling(board.WhiteQueenside):
		return board.Sq("c1")
	case castling(board.BlackQueenside):
		return board.Sq("c8")
	default:
		return board.Sq("g8")
	}
}

// ParseUCI парсит ход из UCI-нотации
func ParseUCI(s string) (Move, error) {
	if len(s) != 4 && len(s) != 5 {
		return nil, fmt.Errorf(errCantParse)
	}

	switch s {
	case "e1g1":
		return castling(board.WhiteKingside), nil
	case "e1c1":
		return castling(board.WhiteQueenside), nil
	case "e8c8":
		return castling(board.BlackQueenside), nil
	case "e8g8":
		return castling(board.BlackKingside), nil
	}

	from := board.Sq(s[:2])
	to := board.Sq(s[2:4])
	if from == -1 || to == -1 {
		return nil, fmt.Errorf(errCantParse)
	}

	var promoteTo board.Piece
	if len(s) == 5 {
		pc := s[4]
		s = s[:4]
		switch s[3] {
		case '8':
			switch pc {
			case 'q':
				promoteTo = board.WhiteQueen
			case 'r':
				promoteTo = board.WhiteRook
			case 'b':
				promoteTo = board.WhiteBishop
			case 'n':
				promoteTo = board.WhiteKing
			default:
				return nil, fmt.Errorf(errCantParse)
			}
		case '1':
			switch pc {
			case 'q':
				promoteTo = board.BlackQueen
			case 'r':
				promoteTo = board.BlackRook
			case 'b':
				promoteTo = board.BlackBishop
			case 'n':
				promoteTo = board.BlackKing
			default:
				return nil, fmt.Errorf(errCantParse)
			}
		default:
			return nil, fmt.Errorf(errCantParse)
		}
	}

	simple := simpleMove{from: from, to: to}

	if promoteTo != 0 {
		return promotion{simpleMove: simple, promoteTo: promoteTo}, nil
	}

	return simple, nil
}
