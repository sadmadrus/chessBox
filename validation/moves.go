package validation

import (
	"fmt"

	"github.com/sadmadrus/chessBox/internal/board"
)

// Кастомные ошибки валидации

var (
	errPieceNotExist                 = fmt.Errorf("piece does not exist")
	errNewpieceNotValid              = fmt.Errorf("newpiece is not valid")
	errInvalidHttpMethod             = fmt.Errorf("method is not supported")
	errFromToSquaresNotDiffer        = fmt.Errorf("from and to squares are not different")
	errPawnMoveNotValid              = fmt.Errorf("pawn move is not valid")
	errKnightMoveNotValid            = fmt.Errorf("knight move is not valid")
	errBishopMoveNotValid            = fmt.Errorf("bishop move is not valid")
	errRookMoveNotValid              = fmt.Errorf("rook move is not valid")
	errQueenMoveNotValid             = fmt.Errorf("queen move is not valid")
	errKingMoveNotValid              = fmt.Errorf("king move is not valid")
	errNoPieceOnFromSquare           = fmt.Errorf("no piece on from square")
	errPieceWrongColor               = fmt.Errorf("piece has wrong color")
	errClashWithPieceOfSameColor     = fmt.Errorf("clash with piece of the same color")
	errClashWithKing                 = fmt.Errorf("clash with king")
	errClashWithPawn                 = fmt.Errorf("pawn can not clash with another piece when moving vertically")
	errNewpieceExist                 = fmt.Errorf("newpiece exists with no pawn promotion")
	errNewpieceNotExist              = fmt.Errorf("newpiece does not exist but pawn promotion required")
	errPieceNotExistOnBoard          = fmt.Errorf("piece does not exist on board")
	errKingChecked                   = fmt.Errorf("king checked after move")
	errCastlingThroughCheckedSquare  = fmt.Errorf("castling is not valid through square under check")
	errCastlingThroughOccupiedSquare = fmt.Errorf("castling is not valid through square occupied by other pieces")
	errPiecesStayInTheWay            = fmt.Errorf("piece or pieces stay in the way of figure move")
	errPiecetypeNotExist             = fmt.Errorf("piece type does not exist")
	errInternalErrorIsSquareChecked  = fmt.Errorf("internal error occured while performing isSquareChecked")
	errBoardNotValid                 = fmt.Errorf("board has no valid position")
)

// Структуры клетки

// square клетка доски, моделирует ряд row и колонку column на шахматной доске в форматах int.
type square struct {
	row    int
	column int
}

// newSquare создает новый экземпляр клетки доски square из представления s int8 пакета board.
func newSquare(s int8) square {
	return square{
		row:    int(s / 8),
		column: int(s % 8),
	}
}

// abs возвращает модуль числа n.
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// diffRow дает разницу в рядах между клетками s1 и s2.
func (s1 *square) diffRow(s2 square) int {
	return s1.row - s2.row
}

// diffColumn дает разницу в колонках между клетками s1 и s2.
func (s1 *square) diffColumn(s2 square) int {
	return s1.column - s2.column
}

// isEqual возвращает true, если клетки s1, s2 одинаковы; в противном случае false.
func (s1 *square) isEqual(s2 square) bool {
	return s1.toInt() == s2.toInt()
}

// toInt возвращает представление клетки s1 в формате int.
func (s1 *square) toInt() int {
	return s1.row*8 + s1.column
}

// toInt8 возвращает представление клетки s1 в формате int8.
func (s1 *square) toInt8() int8 {
	return int8(s1.toInt())
}

// методы move

// move для каждого типа фигуры. В случае невозможности сделать ход, возвращают ошибку, иначе возвращают nil.
func move(piece board.Piece, from, to square) (err error) {
	switch piece {
	case board.WhitePawn, board.BlackPawn:
		err = movePawn(piece, from, to)

	case board.WhiteKnight, board.BlackKnight:
		err = moveKnight(from, to)

	case board.WhiteBishop, board.BlackBishop:
		err = moveBishop(from, to)

	case board.WhiteRook, board.BlackRook:
		err = moveRook(from, to)

	case board.WhiteQueen, board.BlackQueen:
		err = moveQueen(from, to)

	case board.WhiteKing, board.BlackKing:
		err = moveKing(piece, from, to)
	}

	return err
}

// movePawn логика движения пешки. Может двигаться вверх (белый) или вниз (черный) на 1 или 2 клетки. Может съедать
// по диагонали проходные пешки или фигуры соперника. Возвращает ошибку, если движение невалидно.
func movePawn(piece board.Piece, from, to square) error {
	var (
		isVerticalValid bool // разрешено ли движение пешки по вертикали
		isDiagonalValid bool // разрешено ли движение пешки по диагонали
	)

	switch piece {
	case board.WhitePawn:
		isVerticalValid = (from.diffColumn(to) == 0) && // верикаль не изменяется
			(from.row != 0) && // пешка не может стартовать с 0 ряда
			(to.row > 1) && // пешка не может прийти на 0 или 1 ряд
			((from.diffRow(to) == -1) || (to.row == 3 && from.row == 1)) // движение вверх на 1 клетку, либо на 2 клетки (с 1 на 3 ряд)
		isDiagonalValid = (from.row != 0) && // пешка не может стартовать с 0 ряда
			(to.row > 1) && // пешка не может прийти на 0 или 1 ряд
			(from.diffRow(to) == -1) && (abs(from.diffColumn(to)) == 1) // движение вверх по диагонали на 1 клетку

	case board.BlackPawn:
		isVerticalValid = (from.diffColumn(to) == 0) && // верикаль не изменяется
			(from.row != 7) && // пешка не может стартовать с 7 ряда
			(to.row < 6) && // пешка не может прийти на 7 или 6 ряд
			((from.diffRow(to) == 1) || (to.row == 4 && from.row == 6)) // движение вниз на 1 клетку, либо на 2 клетки (с 6 на 4 ряд)
		isDiagonalValid = (from.row != 7) && // пешка не может стартовать с 7 ряда
			(to.row < 6) && // пешка не может прийти на 7 или 6 ряд
			(from.diffRow(to) == 1) && (abs(from.diffColumn(to)) == 1) // движение вниз по диагонали на 1 клетку
	}

	if !isVerticalValid && !isDiagonalValid {
		return fmt.Errorf("%w", errPawnMoveNotValid)
	}
	return nil
}

// moveKnight логика движения коня без привязки к позиции на доске. Может двигаться буквой Г. То есть +/- 2 клетки
// в одном направлении и +/- 1 клетка в перпендикулярном направлении. Возвращает ошибку, если движение невалидно.
func moveKnight(from, to square) error {
	// разрешено ли движение конем на +/- 2 клетки в одном направлении и +/- 1 клетку в перпендикулярном ему направлении
	isValid := (abs(from.diffRow(to)) == 2 && abs(from.diffColumn(to)) == 1) ||
		(abs(from.diffRow(to)) == 1 && abs(from.diffColumn(to)) == 2)

	if !isValid {
		return fmt.Errorf("%w", errKnightMoveNotValid)
	}
	return nil
}

// moveBishop логика движения слона без привязки к позиции на доске. Может двигаться по всем диагоналям. Возвращает
// ошибку, если движение невалидно.
func moveBishop(from, to square) error {
	// разрешено ли движение слоном по диагоналям
	isValid := abs(from.diffRow(to)) == abs(from.diffColumn(to))

	if !isValid {
		return fmt.Errorf("%w", errBishopMoveNotValid)
	}
	return nil
}

// moveRook логика движения ладьи. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток. Возвращает
// ошибку, если движение невалидно.
func moveRook(from, to square) error {
	// разрешено ли движение ладьей
	isValid := (from.diffRow(to) == 0) || // по горизонталям
		(from.diffColumn(to) == 0) // по вертикалям

	if !isValid {
		return fmt.Errorf("%w", errRookMoveNotValid)
	}
	return nil
}

// moveQueen для ферзя. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток. Может двигаться диагонально
// на любое кол-во клеток. Возвращает ошибку, если движение невалидно.
func moveQueen(from, to square) error {
	errBishop := moveBishop(from, to) // может ли ферзь двигаться как слон по диагонаям
	errRook := moveRook(from, to)     // может ли ферзь двигаться как ладья по вертикалям и горизонталям

	if errBishop != nil && errRook != nil {
		return fmt.Errorf("%w", errQueenMoveNotValid)
	}
	return nil
}

// moveKing логика движения короля. Может двигаться вертикально, горизонтально и диагонально только на одну клетку.
// Также король из своего начального положения на доске (row 0 && column 4 для белого; row 7 && column 4 для черного)
// может двигаться: на 2 клетки вправо или 2 клетки влево для рокировок.
func moveKing(piece board.Piece, from, to square) error {
	var (
		isHorizontalValid bool // разрешено ли движение короля по горизонтали
		isVerticalValid   bool // разрешено ли движение короля по вертикали
		isDiagonalValid   bool // разрешено ли движение короля по диагонали
		isCastlingValid   bool // разрешена ли рокировка
	)

	isHorizontalValid = (from.diffRow(to) == 0) && (abs(from.diffColumn(to)) == 1)    // на 1 клетку вправо или влево
	isVerticalValid = (abs(from.diffRow(to)) == 1) && (from.diffColumn(to) == 0)      // на 1 клетку вверх или вниз
	isDiagonalValid = (abs(from.diffRow(to)) == 1) && (abs(from.diffColumn(to)) == 1) // на 1 клетку по любой диагонали

	// определение возможности рокировки
	switch piece {
	case board.WhiteKing:
		if from.row == 0 && from.column == 4 && to.row == 0 && (abs(from.diffColumn(to)) == 2) {
			isCastlingValid = true
		}
	case board.BlackKing:
		if from.row == 7 && from.column == 4 && to.row == 7 && (abs(from.diffColumn(to)) == 2) {
			isCastlingValid = true
		}
	}

	if !isHorizontalValid && !isVerticalValid && !isDiagonalValid && !isCastlingValid {
		return fmt.Errorf("%w", errKingMoveNotValid)
	}
	return nil
}
