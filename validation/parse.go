package validation

import (
	"fmt"
	"strconv"

	"github.com/sadmadrus/chessBox/internal/board"
)

// parsePiece переводит строковое представление фигуры типа k/q/r/b/n/p/K/Q/R/B/N/P в тип board.Piece.  Если
// преобразование невозможно, возвращает ошибку.
func parsePiece(piece string) (board.Piece, error) {
	switch piece {
	case "P":
		return board.WhitePawn, nil
	case "p":
		return board.BlackPawn, nil
	case "N":
		return board.WhiteKnight, nil
	case "n":
		return board.BlackKnight, nil
	case "B":
		return board.WhiteBishop, nil
	case "b":
		return board.BlackBishop, nil
	case "R":
		return board.WhiteRook, nil
	case "r":
		return board.BlackRook, nil
	case "Q":
		return board.WhiteQueen, nil
	case "q":
		return board.BlackQueen, nil
	case "K":
		return board.WhiteKing, nil
	case "k":
		return board.BlackKing, nil
	default:
		return 0, fmt.Errorf("%w", errPieceNotExist)
	}
}

// parseSquare переводит строковое представление клетки from/to в клетку структуры square. При невалидных входных данных
// (номер клетки выходит за пределы от 0 до 63 или клетка не вида "а1") выдается ошибку, в противном случае nil.
func parseSquare(squareString string) (squareSquare square, err error) {
	// перевод в тип board.square для цифро-буквенного обозначения клетки (напр., "а1")
	sq := board.Sq(squareString)

	if sq == -1 {
		// перевод в тип board.square для числового значения клетки от 0 до 63
		var sqParsedNum int
		sqParsedNum, err = strconv.Atoi(squareString)
		sq = board.Sq(sqParsedNum)
		if sq == -1 || err != nil {
			return newSquare(-1), fmt.Errorf("%v: %v", errPieceNotExist, squareString)
		}
	}

	return newSquare(int8(sq)), nil
}

// parseNewpiece переводит строковое представление фигуры типа q/r/b/n/Q/R/B/N или пустое значение в тип board.Piece.
// Если преобразование невозможно, возвращает ошибку.
func parseNewpiece(newpiece string) (board.Piece, error) {
	switch newpiece {
	case "N":
		return board.WhiteKnight, nil
	case "n":
		return board.BlackKnight, nil
	case "B":
		return board.WhiteBishop, nil
	case "b":
		return board.BlackBishop, nil
	case "R":
		return board.WhiteRook, nil
	case "r":
		return board.BlackRook, nil
	case "Q":
		return board.WhiteQueen, nil
	case "q":
		return board.BlackQueen, nil
	case "":
		return 0, nil
	default:
		return 0, fmt.Errorf("%w", errNewpieceNotValid)
	}
}
