package moves

import (
	"fmt"
	"sort"

	"github.com/sadmadrus/chessBox/internal/board"
)

// Кастомные ошибки валидации

var (
	ErrBoardNotValid          = fmt.Errorf("board has no valid position")
	ErrSquareNotExist         = fmt.Errorf("square does not exist")
	ErrPromoteToNotValid      = fmt.Errorf("promoteTo is not valid")
	ErrFromToSquaresNotDiffer = fmt.Errorf("from and to squares are not different")
	ErrNoPieceOnFromSquare    = fmt.Errorf("no piece on from square")
	ErrPieceWrongColor        = fmt.Errorf("piece has wrong color")

	errPawnMoveNotValid              = fmt.Errorf("pawn move is not valid")
	errKnightMoveNotValid            = fmt.Errorf("knight move is not valid")
	errBishopMoveNotValid            = fmt.Errorf("bishop move is not valid")
	errRookMoveNotValid              = fmt.Errorf("rook move is not valid")
	errQueenMoveNotValid             = fmt.Errorf("queen move is not valid")
	errKingMoveNotValid              = fmt.Errorf("king move is not valid")
	errClashWithPieceOfSameColor     = fmt.Errorf("clash with piece of the same color")
	errClashWithKing                 = fmt.Errorf("clash with king")
	errClashWithPawn                 = fmt.Errorf("pawn can not clash with another piece when moving vertically")
	errPieceNotExistOnBoard          = fmt.Errorf("piece does not exist on board")
	errCastlingThroughCheckedSquare  = fmt.Errorf("castling is not valid through square under check")
	errCastlingThroughOccupiedSquare = fmt.Errorf("castling is not valid through square occupied by other pieces")
	errPiecesStayInTheWay            = fmt.Errorf("piece or pieces stay in the way of figure move")
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

// методы getMoves

// getMoves возвращает все допустимые ходы фигуры с клетки from, без привязки к доске. Срез клеток возвращается
// отсортированным по возрастанию.
func getMoves(b board.Board, from square) (moves []square, err error) {
	var piece board.Piece
	piece, err = b.Get(board.Sq(from.toInt()))
	if err != nil {
		return moves, err
	}

	switch piece {
	case board.WhitePawn, board.BlackPawn:
		moves = getPawnMoves(piece, from)
	case board.WhiteKnight, board.BlackKnight:
		moves = getKnightMoves(from)
	case board.WhiteBishop, board.BlackBishop:
		moves = getBishopMoves(from)
	case board.WhiteRook, board.BlackRook:
		moves = getRookMoves(from)
	case board.WhiteQueen, board.BlackQueen:
		moves = getQueenMoves(from)
	case board.WhiteKing, board.BlackKing:
		moves = getKingMoves(piece, from)
	}

	sort.Slice(moves, func(i, j int) bool { return moves[i].toInt() < moves[j].toInt() })

	return moves, nil
}

// getPawnMoves возвращает все допустимые ходы пешки с клетки from, без привязки к доске.
func getPawnMoves(piece board.Piece, from square) (moves []square) {
	switch piece {
	case board.WhitePawn:
		if from.row == 1 {
			moves = append(moves, newSquare(from.toInt8()+16))
		}
		if from.row < 7 && from.row > 0 {
			if from.column > 0 {
				moves = append(moves, newSquare(from.toInt8()+7))
			}
			moves = append(moves, newSquare(from.toInt8()+8))
			if from.column < 7 {
				moves = append(moves, newSquare(from.toInt8()+9))
			}
		}

	case board.BlackPawn:
		if from.row == 6 {
			moves = append(moves, newSquare(from.toInt8()-16))
		}
		if from.row > 0 && from.row < 7 {
			if from.column > 0 {
				moves = append(moves, newSquare(from.toInt8()-9))
			}
			moves = append(moves, newSquare(from.toInt8()-8))
			if from.column < 7 {
				moves = append(moves, newSquare(from.toInt8()-7))
			}
		}
	}

	return moves
}

// getKnightMoves возвращает все допустимые ходы коня с клетки from, без привязки к доске.
func getKnightMoves(from square) (moves []square) {
	if from.row < 6 && from.column < 7 {
		moves = append(moves, newSquare(from.toInt8()+17))
	}
	if from.row < 7 && from.column < 6 {
		moves = append(moves, newSquare(from.toInt8()+10))
	}
	if from.row > 0 && from.column < 6 {
		moves = append(moves, newSquare(from.toInt8()-6))
	}
	if from.row > 1 && from.column < 7 {
		moves = append(moves, newSquare(from.toInt8()-15))
	}

	if from.row > 1 && from.column > 0 {
		moves = append(moves, newSquare(from.toInt8()-17))
	}
	if from.row > 0 && from.column > 1 {
		moves = append(moves, newSquare(from.toInt8()-10))
	}
	if from.row < 7 && from.column > 1 {
		moves = append(moves, newSquare(from.toInt8()+6))
	}
	if from.row < 6 && from.column > 0 {
		moves = append(moves, newSquare(from.toInt8()+15))
	}

	return moves
}

// getBishopMoves возвращает все допустимые ходы слона с клетки from, без привязки к доске.
func getBishopMoves(from square) (moves []square) {
	row, column := from.row, from.column
	for row <= 7 && column <= 7 {
		if row != from.row && column != from.column {
			moves = append(moves, newSquare(int8(row*8+column)))
		}
		row++
		column++
	}

	row, column = from.row, from.column
	for row <= 7 && column >= 0 {
		if row != from.row && column != from.column {
			moves = append(moves, newSquare(int8(row*8+column)))
		}
		row++
		column--
	}

	row, column = from.row, from.column
	for row >= 0 && column >= 0 {
		if row != from.row && column != from.column {
			moves = append(moves, newSquare(int8(row*8+column)))
		}
		row--
		column--
	}

	row, column = from.row, from.column
	for row >= 0 && column <= 7 {
		if row != from.row && column != from.column {
			moves = append(moves, newSquare(int8(row*8+column)))
		}
		row--
		column++
	}

	return moves
}

// getRookMoves возвращает все допустимые ходы ладьи с клетки from, без привязки к доске.
func getRookMoves(from square) (moves []square) {
	row := from.row
	for row <= 7 {
		if row != from.row {
			moves = append(moves, newSquare(int8(row*8+from.column)))
		}
		row++
	}

	row = from.row
	for row >= 0 {
		if row != from.row {
			moves = append(moves, newSquare(int8(row*8+from.column)))
		}
		row--
	}

	column := from.column
	for column <= 7 {
		if column != from.column {
			moves = append(moves, newSquare(int8(from.row*8+column)))
		}
		column++
	}

	column = from.column
	for column >= 0 {
		if column != from.column {
			moves = append(moves, newSquare(int8(from.row*8+column)))
		}
		column--
	}

	return moves
}

// getQueenMoves возвращает все допустимые ходы ферзя с клетки from, без привязки к доске.
func getQueenMoves(from square) (moves []square) {
	bishopMoves := getBishopMoves(from)
	rookMoves := getRookMoves(from)

	moves = append(moves, bishopMoves...)
	moves = append(moves, rookMoves...)

	return moves
}

// getKingMoves возвращает все допустимые ходы короля с клетки from, без привязки к доске.
func getKingMoves(piece board.Piece, from square) (moves []square) {
	// рокировки
	switch piece {
	case board.WhiteKing:
		if from.row == 0 && from.column == 4 {
			moves = append(moves, newSquare(2))
			moves = append(moves, newSquare(6))
		}

	case board.BlackKing:
		if from.row == 7 && from.column == 4 {
			moves = append(moves, newSquare(58))
			moves = append(moves, newSquare(62))
		}
	}

	// вертикали и горизонтали
	if from.row > 0 {
		moves = append(moves, newSquare(from.toInt8()-8))
	}
	if from.row < 7 {
		moves = append(moves, newSquare(from.toInt8()+8))
	}
	if from.column > 0 {
		moves = append(moves, newSquare(from.toInt8()-1))
	}
	if from.column < 7 {
		moves = append(moves, newSquare(from.toInt8()+1))
	}

	// диагонали
	if from.row > 0 && from.column > 0 {
		moves = append(moves, newSquare(from.toInt8()-9))
	}
	if from.row > 0 && from.column < 7 {
		moves = append(moves, newSquare(from.toInt8()-7))
	}
	if from.row < 7 && from.column > 0 {
		moves = append(moves, newSquare(from.toInt8()+7))
	}
	if from.row < 7 && from.column < 7 {
		moves = append(moves, newSquare(from.toInt8()+9))
	}

	return moves
}
