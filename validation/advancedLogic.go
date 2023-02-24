package validation

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sadmadrus/chessBox/internal/board"
)

// advancedLogic обрабывает общую логику валидации хода. Возвращает доску и флаг валидации хода (true валиден, false нет).
// В случае ошибки при обработке хода (некорректная клетка, фигура и т.п.), возвращает ошибку, в противном случае nil.
func advancedLogic(b board.Board, from, to square, newpiece board.Piece) (newBoard board.Board, isValid bool, err error) {
	// Логика валидации хода пошагово.
	// 1. получаем фигуру, находящуюся в клетке from. Если в этой клетке фигуры нет, ход невалиден.
	var piece board.Piece
	piece, err = b.Get(board.Sq(from.toInt()))
	if err != nil {
		return newBoard, isValid, err
	}
	if piece == 0 {
		log.Printf("%v: %v", errNoPieceOnFromSquare, from)
		return newBoard, isValid, nil
	}

	// 2a. проверяем, что пользователь указал, какую новую фигуру выставить в случае проведения пешки. Если фигура
	// не указана, ход невалиден. TODO: или логичнее выдавать ошибку?
	if newpiece == 0 && ((piece == board.WhitePawn && to.row == 7) || (piece == board.BlackPawn && to.row == 0)) {
		log.Printf("%v", errNewpieceNotExist)
		return newBoard, isValid, nil
	}

	// 2b. проверяем, что пользователь не захотел выставить нового фигуру в неуместном для этого случае. В таком случае
	// ход будет невалиден. TODO: или логичнее выдавать ошибку?
	if newpiece != 0 && ((piece != board.WhitePawn && to.row != 7) || (piece != board.BlackPawn && to.row != 0)) {
		log.Printf("%v", errNewpieceExist)
		return newBoard, isValid, nil
	}

	// 3. проверяем, что фигура принадлежит той стороне, чья очередь хода. Иначе ход невалиден.
	isFigureRightColor := checkFigureColor(b, piece)
	if !isFigureRightColor {
		log.Printf("%v", errPieceWrongColor)
		return newBoard, isValid, nil
	}

	// 4. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	err = move(piece, from, to)
	if err != nil {
		log.Printf("%v", err)
		return newBoard, isValid, nil
	}

	// 5. Проверяем, что по пути фигуры с клетки from до to (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	squaresToBePassed := getSquaresToBePassed(piece, from, to)
	if len(squaresToBePassed) > 0 {
		var areSquaresEmpty bool
		areSquaresEmpty, err = checkSquaresToBePassed(b, squaresToBePassed)
		if err != nil {
			return newBoard, isValid, err
		}
		if !areSquaresEmpty {
			log.Printf("%v", errPiecesStayInTheWay)
			return newBoard, isValid, nil
		}
	}

	// TODO: остановилась здесь

	// 6. Проверяем наличие и цвет фигур в клетке to.
	err = checkToSquare(&b, piece, from, to)
	if err != nil {
		return newBoard, isValid, err
	}

	// 7. Проверка, что при рокировке король не проходит через битое поле.
	var king board.Piece
	if b.NextToMove() {
		king = board.WhiteKing
	} else {
		king = board.BlackKing
	}
	if (piece == board.WhiteKing || piece == board.BlackKing) && (abs(from.diffColumn(to)) == 2) {
		squareToBePassed := newSquare(from.toInt8() + ((to.toInt8() - from.toInt8()) / 2))
		var squareChecked bool
		squareChecked, err = isSquareChecked(b, squareToBePassed, king)
		if err != nil {
			return newBoard, isValid, err
		}
		if squareChecked {
			return newBoard, isValid, fmt.Errorf("%v", errCastlingThroughCheckedSquare)
		}
	}

	// 8. На текущем этапе ход возможен. Генерируем новое положение доски newBoard. Выдается ошибка при некорректном
	// проведении пешки, некорректной рокировке, некорректном взятии на проходе по логике из пакета board.
	newBoard, err = getNewBoard(b, piece, from, to, newpiece)
	if err != nil {
		return newBoard, isValid, fmt.Errorf("%w", err)
	}

	// 9. Проверяем, что при новой позиции на доске не появился шах для собственного короля.
	var kingChecked bool
	kingChecked, err = isKingChecked(b, king)
	if err != nil {
		return newBoard, isValid, err
	}
	if kingChecked {
		return newBoard, isValid, fmt.Errorf("%v", errKingChecked)
	}

	// 10. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.
	if piece == board.WhiteKing || piece == board.BlackKing {
		var isEnemyKingAdjacent bool
		isEnemyKingAdjacent, err = checkDistanceToEnemyKing(b)
		if err != nil {
			return newBoard, isValid, err
		}
		if isEnemyKingAdjacent {
			return newBoard, isValid, fmt.Errorf("%v", errKingsAdjacent)
		}
	}

	isValid = true
	return newBoard, isValid, nil
}

// checkFigureColor проверяет, что очередь хода и цвет фигуры p, которую хотят передвинуть, совпадают.
// Возвращает true в случае успеха, false в противном случае.
func checkFigureColor(b board.Board, p board.Piece) bool {
	var pieceIsWhite bool

	switch p {
	case board.WhitePawn, board.WhiteKnight, board.WhiteBishop, board.WhiteRook, board.WhiteKing, board.WhiteQueen:
		pieceIsWhite = true
	case board.BlackPawn, board.BlackKnight, board.BlackBishop, board.BlackRook, board.BlackKing, board.BlackQueen:
		pieceIsWhite = false
	}

	if (b.NextToMove() && pieceIsWhite) || (!b.NextToMove() && !pieceIsWhite) {
		return true
	}

	return false
}

// getSquaresToBePassed возвращает массив клеток, через которые проходит фигура p при движении из клетки from на клетку
// to. Если таких клеток нет, то возвращаемый массив пуст.
func getSquaresToBePassed(p board.Piece, from, to square) (squaresToBePassed []square) {
	var verticalDirection, horizontalDirection int8

	if from.diffRow(to) < 0 {
		verticalDirection = 1
	} else if from.diffRow(to) > 0 {
		verticalDirection = -1
	} else {
		verticalDirection = 0
	}

	if from.diffColumn(to) < 0 {
		horizontalDirection = 1
	} else if from.diffColumn(to) > 0 {
		horizontalDirection = -1
	} else {
		horizontalDirection = 0
	}

	squaresToBePassedAmount := abs(from.diffRow(to))
	if abs(from.diffColumn(to)) > abs(from.diffRow(to)) {
		squaresToBePassedAmount = abs(from.diffColumn(to))
	}

	switch p {
	case board.WhitePawn, board.BlackPawn, board.WhiteBishop, board.BlackBishop, board.WhiteRook, board.BlackRook,
		board.WhiteKing, board.BlackKing, board.WhiteQueen, board.BlackQueen:
		var i int8 = 1
		for squaresToBePassedAmount > 1 {
			squareToBeAdded := newSquare(from.toInt8() + (verticalDirection * 8 * i) + (horizontalDirection * i))
			squaresToBePassed = append(squaresToBePassed, squareToBeAdded)
			squaresToBePassedAmount--
			i++
		}
	}

	return squaresToBePassed
}

// checkSquaresToBePassed проверяет, есть ли на клетках из массива squaresToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается флаг false (иначе true). Если при обработке запроса возникает ошибка, она также
// возвращается (иначе nil).
func checkSquaresToBePassed(b board.Board, squaresToBePassed []square) (areSquaresEmpty bool, err error) {
	for _, sq := range squaresToBePassed {
		var piece board.Piece
		piece, err = b.Get(board.Sq(sq.toInt()))
		if err != nil {
			log.Printf("%v", err)
			return areSquaresEmpty, err
		}
		if piece != 0 {
			log.Printf("%v: %v", errPieceFound, sq)
			return areSquaresEmpty, nil
		}
	}

	areSquaresEmpty = true
	return areSquaresEmpty, nil
}

//TODO: остановилась здесь

// checkToSquare проверяет наличие фигуры на клетке to на предмет совместимости хода. Возвращает ошибку при
// несовместимости хода или nil в случае успеха.
func checkToSquare(b *board.Board, pieceFrom board.Piece, from, to square) error {
	pieceTo, _ := b.Get(board.Sq(to.toInt()))

	// Если в клетке to нет фигур, ход возможен.
	if pieceTo == 0 {
		return nil
	}

	// Если фигура в to принадлежит самому участнику, ход невозможен.
	if checkFigureColor(*b, pieceTo) {
		return fmt.Errorf("%v", errClashWithPieceOfSameColor)

		// Если фигура в to принадлежит сопернику, проверка, возможно ли взятие
	} else {
		// ни одна фигура не может взять короля
		switch pieceTo {
		case board.WhiteKing, board.BlackKing:
			return fmt.Errorf("%v", errClashWithKing)
		}

		// пешка не может взять ни одну фигуру при движении вертикально
		switch pieceFrom {
		case board.WhitePawn, board.BlackPawn:
			if from.diffColumn(to) == 0 {
				return fmt.Errorf("%v", errClashWithPawn)
			}
		}
	}

	return nil
}

// getNewBoard генерирует новое положение доски с учетом рокировок, взятия на прохоже и проведения пешки. Возвращает
// ошибку при некорректных входных данных.
func getNewBoard(b board.Board, piece board.Piece, from, to square, newpiece board.Piece) (board.Board, error) {
	// обработка проведения белой и черной пешки
	if (piece == board.WhitePawn && to.row == 7) || (piece == board.BlackPawn && to.row == 0) {
		err := b.Promote(board.Sq(from.toInt()), board.Sq(to.toInt()), newpiece)
		if err != nil {
			return b, err
		}
		return b, nil
	}

	// обработка рокировки белого и черного короля
	if (piece == board.WhiteKing || piece == board.BlackKing) && (abs(from.diffColumn(to)) == 2) {
		// сторона короля
		if to.column == 6 {
			var castling board.Castling
			switch piece {
			case board.WhiteKing:
				castling = board.WhiteKingside
			case board.BlackKing:
				castling = board.BlackKingside
			}

			err := b.Castle(castling)
			if err != nil {
				return b, err
			}

			// сторона королевы
		} else if to.column == 2 {
			var castling board.Castling
			switch piece {
			case board.WhiteKing:
				castling = board.WhiteQueenside
			case board.BlackKing:
				castling = board.BlackQueenside
			}

			err := b.Castle(castling)
			if err != nil {
				return b, err
			}
		}
	}

	// обработка всех остальных ходов
	err := b.Move(board.Sq(from.toInt()), board.Sq(to.toInt()))
	if err != nil {
		return b, err
	}
	return b, nil
}

// isKingChecked проверяет, что свой король не под шахом. Если шах есть, возвращает true,
// иначе false. Возвращает ошибку, если возникла при обработке, иначе nil.
func isKingChecked(b board.Board, king board.Piece) (isKingChecked bool, err error) {
	pieceString := king.String()
	var kingSquare square
	kingSquare, err = getSquareByPiece(b, pieceString)
	if err != nil {
		return isKingChecked, err
	}

	isKingChecked, err = isSquareChecked(b, kingSquare, king)
	if err != nil {
		return isKingChecked, err
	}

	return isKingChecked, nil
}

// isSquareChecked проверяет, что на доске b нет шаха королю king, когда он находится в клетке sq. Если шах есть,
// возвращает true, иначе false. Возвращает ошибку, если возникла при обработке, иначе nil.
func isSquareChecked(b board.Board, sq square, king board.Piece) (isSquareChecked bool, err error) {
	var enemyKnight, enemyRook, enemyQueen, enemyBishop, enemyPawn board.Piece
	if king == board.WhiteKing {
		enemyKnight = board.BlackKnight
		enemyRook = board.BlackRook
		enemyQueen = board.BlackQueen
		enemyBishop = board.BlackBishop
		enemyPawn = board.BlackPawn
	} else {
		enemyKnight = board.WhiteKnight
		enemyRook = board.WhiteRook
		enemyQueen = board.WhiteQueen
		enemyBishop = board.WhiteBishop
		enemyPawn = board.WhitePawn
	}

	// проверяем, есть ли на расстоянии буквы Г от клетки sq вражеские кони.
	isSquareChecked, err = checkEnemyKnightsNearKing(b, sq, enemyKnight)
	if err != nil {
		return isSquareChecked, err
	}
	if isSquareChecked {
		return isSquareChecked, nil
	}

	// проверяем, есть ли по вертикали или горизонтали в качестве ближайших фигур вражеские ладьи и ферзи.
	isSquareChecked, err = checkEnemiesVerticallyAndHorizontally(b, sq, enemyRook, enemyQueen)
	if err != nil {
		return isSquareChecked, err
	}
	if isSquareChecked {
		return isSquareChecked, nil
	}

	// проверяем, есть ли по диагоналям в качестве ближайших фигур вражеские ферзи, слоны и пешки.
	isSquareChecked, err = checkEnemiesDiagonally(b, sq, enemyQueen, enemyBishop, enemyPawn)
	if err != nil {
		return isSquareChecked, err
	}

	return isSquareChecked, nil

}

// getSquareByPiece возвращает клетку, на которой находится заданная фигура. Если такой фигуры на доске нет,
// возвращает ошибку.
func getSquareByPiece(b board.Board, pieceString string) (pieceSquare square, err error) {
	boardFEN := b.FEN()
	boardFENArr := strings.Split(boardFEN, " ")
	boardFEN = boardFENArr[0]
	index := strings.Index(boardFEN, pieceString)
	if index == -1 {
		return pieceSquare, fmt.Errorf("%v", errPieceNotExistOnBoard)
	}

	rowsCount := strings.Count(boardFEN[:index], "/")
	var squareRow = 7 - int8(rowsCount)

	var boardRow = boardFEN
	for rowsCount > 0 {
		slashIndex := strings.Index(boardRow, "/")
		boardRow = boardRow[slashIndex:]
		rowsCount--
	}
	slashIndex := strings.Index(boardRow, "/")
	boardRow = boardRow[:slashIndex]

	var squareColumn int
	for _, sym := range boardRow {
		if string(sym) == pieceString {
			break
		}

		switch sym {
		case 'P', 'p', 'K', 'k', 'Q', 'q', 'R', 'r', 'N', 'n', 'B', 'b':
			squareColumn++
		default:
			num, _ := strconv.Atoi(string(sym))
			squareColumn += num
		}
	}

	pieceSquare = newSquare(squareRow*8 + int8(squareColumn))
	return pieceSquare, nil
}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если есть хотя бы одна такая клетка, возвращает true, иначе false. Если при проверке клеток
// возникает ошибка, она также возвращается, иначе возвращется nil.
func checkEnemyKnightsNearKing(b board.Board, kingSquare square, enemyKnight board.Piece) (isEnemyKnightPresent bool, err error) {
	var squaresToBeChecked []square
	// +2 клетки вверх, +1 клетка вправо
	if kingSquare.row <= 5 && kingSquare.column <= 6 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()+int8(17)))
	}
	// +1 клетка вверх, +2 клетки вправо
	if kingSquare.row <= 6 && kingSquare.column <= 5 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()+int8(10)))
	}
	// -1 клетка вниз, +2 клетки вправо
	if kingSquare.row >= 1 && kingSquare.column <= 5 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()-int8(6)))
	}
	// -2 клетка вниз, +1 клетки вправо
	if kingSquare.row >= 2 && kingSquare.column <= 6 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()-int8(15)))
	}
	// -2 клетка вниз, -1 клетки влево
	if kingSquare.row >= 2 && kingSquare.column >= 1 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()-int8(17)))
	}
	// -1 клетка вниз, -2 клетки влево
	if kingSquare.row >= 1 && kingSquare.column >= 2 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()-int8(10)))
	}
	// +1 клетка вверх, -2 клетки влево
	if kingSquare.row <= 6 && kingSquare.column >= 2 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()+int8(6)))
	}
	// +2 клетка вверх, -1 клетки влево
	if kingSquare.row <= 5 && kingSquare.column >= 1 {
		squaresToBeChecked = append(squaresToBeChecked, newSquare(kingSquare.toInt8()+int8(15)))
	}

	for _, sq := range squaresToBeChecked {
		var piece board.Piece
		piece, err = b.Get(board.Sq(sq.toInt()))
		if err != nil {
			return isEnemyKnightPresent, err
		}
		if piece == enemyKnight {
			isEnemyKnightPresent = true
		}
	}

	return isEnemyKnightPresent, nil
}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingSquare, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b board.Board, kingSquare square, enemyRook, enemyQueen board.Piece) (isEnemyVerticallyOrHorizontallyPresent bool, err error) {
	// проверка вертикали вверх
	var row = kingSquare.row
	var squareToBeChecked = kingSquare.toInt()
	for row < 7 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked + 8))
		if err != nil {
			return isEnemyVerticallyOrHorizontallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyRook, enemyQueen:
			isEnemyVerticallyOrHorizontallyPresent = true
			return isEnemyVerticallyOrHorizontallyPresent, nil
		default:
			break
		}

		row++
	}

	// проверка вертикали вниз
	row = kingSquare.row
	for row > 0 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked - 8))
		if err != nil {
			return isEnemyVerticallyOrHorizontallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyRook, enemyQueen:
			isEnemyVerticallyOrHorizontallyPresent = true
			return isEnemyVerticallyOrHorizontallyPresent, nil
		default:
			break
		}

		row--
	}

	var column = kingSquare.column
	// проверка горизонтали вправо
	for column < 7 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked + 1))
		if err != nil {
			return isEnemyVerticallyOrHorizontallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyRook, enemyQueen:
			isEnemyVerticallyOrHorizontallyPresent = true
			return isEnemyVerticallyOrHorizontallyPresent, nil
		default:
			break
		}

		column++
	}

	column = kingSquare.column
	// проверка горизонтали влево
	for column > 0 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked - 1))
		if err != nil {
			return isEnemyVerticallyOrHorizontallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyRook, enemyQueen:
			isEnemyVerticallyOrHorizontallyPresent = true
			return isEnemyVerticallyOrHorizontallyPresent, nil
		default:
			break
		}

		column--
	}

	return isEnemyVerticallyOrHorizontallyPresent, nil
}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingSquare, на
// которых находятся вражеские слоны, ферзи и пешки. Если есть хотя бы одна такая клетка, возвращает true, иначе false.
// Если при проверке клеток возникает ошибка, она также возвращается, иначе возвращется nil.
func checkEnemiesDiagonally(b board.Board, kingSquare square, enemyQueen, enemyBishop, enemyPawn board.Piece) (isEnemyDiagonallyPresent bool, err error) {
	var row = kingSquare.row
	var column = kingSquare.column

	// проверка диагонали вправо-вверх
	var squareToBeChecked = kingSquare.toInt()
	for row < 7 && column < 7 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked + 9))
		if err != nil {
			return isEnemyDiagonallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyQueen, enemyBishop:
			isEnemyDiagonallyPresent = true
			return isEnemyDiagonallyPresent, nil
		case enemyPawn:
			if enemyPawn == board.BlackPawn && (abs(row-kingSquare.row) == 1) {
				isEnemyDiagonallyPresent = true
				return isEnemyDiagonallyPresent, nil
			}
		default:
			break
		}

		row++
		column++
	}

	// проверка диагонали вправо-вниз
	row = kingSquare.row
	column = kingSquare.column
	for row > 1 && column < 7 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked - 7))
		if err != nil {
			return isEnemyDiagonallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyQueen, enemyBishop:
			isEnemyDiagonallyPresent = true
			return isEnemyDiagonallyPresent, nil
		case enemyPawn:
			if enemyPawn == board.WhitePawn && (abs(row-kingSquare.row) == 1) {
				isEnemyDiagonallyPresent = true
				return isEnemyDiagonallyPresent, nil
			}
		default:
			break
		}

		row--
		column++
	}

	// проверка диагонали влево-вниз
	row = kingSquare.row
	column = kingSquare.column
	for row > 1 && column > 1 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked - 9))
		if err != nil {
			return isEnemyDiagonallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyQueen, enemyBishop:
			isEnemyDiagonallyPresent = true
			return isEnemyDiagonallyPresent, nil
		case enemyPawn:
			if enemyPawn == board.WhitePawn && (abs(row-kingSquare.row) == 1) {
				isEnemyDiagonallyPresent = true
				return isEnemyDiagonallyPresent, nil
			}
		default:
			break
		}

		row--
		column--
	}

	// проверка диагонали влево-вверх
	row = kingSquare.row
	column = kingSquare.column
	for row < 7 && column > 1 {
		var piece board.Piece
		piece, err = b.Get(board.Sq(squareToBeChecked + 7))
		if err != nil {
			return isEnemyDiagonallyPresent, err
		}

		switch piece {
		case 0:
			continue
		case enemyQueen, enemyBishop:
			isEnemyDiagonallyPresent = true
			return isEnemyDiagonallyPresent, nil
		case enemyPawn:
			if enemyPawn == board.BlackPawn && (abs(row-kingSquare.row) == 1) {
				isEnemyDiagonallyPresent = true
				return isEnemyDiagonallyPresent, nil
			}
		default:
			break
		}

		row++
		column--
	}

	return isEnemyDiagonallyPresent, nil
}

// checkEnemiesDiagonally проверяет ближайшие клетки (вертикально, горизонтально, диагонально к клетке своего короля
// kingSquare на наличие на них вражеского короля. Если вражеский король оказался вплотную к своему королю, возвращает
// true, иначе false. Если при проверке клеток возникает ошибка, она также возвращается, иначе возвращется nil.
func checkDistanceToEnemyKing(b board.Board) (isEnemyKingAdjacent bool, err error) {
	var whiteKingSquare, blackKingSquare square
	whiteKingSquare, err = getSquareByPiece(b, "K")
	if err != nil {
		return isEnemyKingAdjacent, err
	}
	blackKingSquare, err = getSquareByPiece(b, "k")
	if err != nil {
		return isEnemyKingAdjacent, err
	}

	if abs(whiteKingSquare.diffRow(blackKingSquare)) <= 1 && abs(whiteKingSquare.diffColumn(blackKingSquare)) <= 1 {
		isEnemyKingAdjacent = true
	}

	return isEnemyKingAdjacent, nil
}
