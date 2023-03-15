package validation

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/board/position"
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

	// 2. проверяем корректность данных о проходе пешки
	isOk := checkPawnPromotion(piece, to, newpiece)
	if !isOk {
		log.Printf("%v or %v: %v %v %v %v (piece, from, to, newpiece)", errNewpieceNotExist, errNewpieceExist, piece, from, to, newpiece)
		return newBoard, isValid, nil
	}

	// 3. проверяем, что фигура принадлежит той стороне, чья очередь хода. Иначе ход невалиден.
	isOk = checkPieceColor(b, piece)
	if !isOk {
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
	var areSquaresEmpty bool
	areSquaresEmpty, err = checkSquaresToBePassed(b, squaresToBePassed)
	if err != nil || !areSquaresEmpty {
		return newBoard, isValid, err
	}

	// 6. Проверяем наличие и цвет фигур в клетке to. Обработка корректного взятия пешки на проходе (en Passant)
	// валидирована здесь.
	isOk, err = checkToSquare(&b, piece, from, to)
	if err != nil || !isOk {
		return newBoard, isValid, err
	}

	// 7. Валидация рокировки: проверка, что при рокировке король не проходит через битое поле, рокировка разрешена
	// (король и ладья еще не двигались, между ними все клетки пустые).
	if (piece == board.WhiteKing || piece == board.BlackKing) && (abs(from.diffColumn(to)) == 2) {
		isOk, err = checkCastling(&b, piece, from, to)
		if err != nil || !isOk {
			return newBoard, isValid, err
		}
	}

	// 8. На текущем этапе ход возможен. Генерируем новое положение доски newBoard. Так как до текущего положения ход
	// валидирован, ошибок не ожидаем.
	newBoard, err = getNewBoard(b, piece, from, to, newpiece)
	if err != nil {
		return newBoard, isValid, err
	}

	// 9. Проверяем, что при новой позиции на доске не появился шах для собственного короля. На этом шаге также
	// проверяем, что король не находится вплотную к чужому королю - такой ход будет запрещен.
	king := "k"
	if b.NextToMove() {
		king = "K"
	}
	kingSquare, _ := getSquareByPiece(newBoard, king)
	checks := position.ThreatsTo(board.Sq(kingSquare.toInt()), newBoard)
	if len(checks) != 0 {
		log.Printf("%v", errKingChecked)
		return newBoard, isValid, nil
	}

	isValid = true
	return newBoard, isValid, nil
}

// checkPawnPromotion проверяет, что указанная пользователем новая фигура newpiece корректна относительно хода.
// Возвращает true (фигура указана корректно) или false (некорректно).
func checkPawnPromotion(piece board.Piece, to square, newpiece board.Piece) (isOk bool) {
	// Если этим ходом проводится пешка, должна быть указана фигура. Если пешка не проводится, фигура не должна быть
	// указана.
	if (piece == board.WhitePawn && to.row == 7) || (piece == board.BlackPawn && to.row == 0) {
		if newpiece != 0 {
			isOk = true
		}
	} else {
		if newpiece == 0 {
			isOk = true
		}
	}

	return isOk
}

// checkPieceColor проверяет, что очередь хода и цвет фигуры p, которую хотят передвинуть, совпадают.
// Возвращает true в случае успеха, false в противном случае.
func checkPieceColor(b board.Board, p board.Piece) (isOk bool) {
	var pieceIsWhite bool

	switch p {
	case board.WhitePawn, board.WhiteKnight, board.WhiteBishop, board.WhiteRook, board.WhiteKing, board.WhiteQueen:
		pieceIsWhite = true
	case board.BlackPawn, board.BlackKnight, board.BlackBishop, board.BlackRook, board.BlackKing, board.BlackQueen:
		pieceIsWhite = false
	}

	if (b.NextToMove() && pieceIsWhite) || (!b.NextToMove() && !pieceIsWhite) {
		isOk = true
	}

	return isOk
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
	if len(squaresToBePassed) == 0 {
		areSquaresEmpty = true
		return areSquaresEmpty, nil
	}

	for _, sq := range squaresToBePassed {
		var piece board.Piece
		piece, err = b.Get(board.Sq(sq.toInt()))
		if err != nil {
			return areSquaresEmpty, err
		}
		if piece != 0 {
			log.Printf("%v: %v", errPiecesStayInTheWay, sq)
			return areSquaresEmpty, nil
		}
	}

	areSquaresEmpty = true
	return areSquaresEmpty, nil
}

// checkToSquare проверяет наличие фигуры на клетке to на предмет совместимости хода. Возвращает флаг (true если фигура
// в клетке to совместима с цветом и типом фигуры pieceFrom, false в противном случае.  Обработка корректного взятия
// пешки на проходе (en Passant) валидирована здесь. Если при обработке возникает ошибка, возвращает ее (иначе nil).
func checkToSquare(b *board.Board, pieceFrom board.Piece, from, to square) (isOk bool, err error) {
	var pieceTo board.Piece
	pieceTo, err = b.Get(board.Sq(to.toInt()))
	if err != nil {
		return isOk, err
	}

	// Если в клетке to нет фигур, ход возможен.
	if pieceTo == 0 {

		// Исключение: взятие пешки на проходе (диагонально) в том случае, если клетка to не указана как en Passant
		if (pieceFrom == board.WhitePawn || pieceFrom == board.BlackPawn) && abs(from.diffColumn(to)) == 1 {
			if !b.IsEnPassant(board.Sq(to.toInt())) {
				return isOk, nil
			}
		}

		isOk = true
		return isOk, nil
	}

	// Если фигура в to принадлежит самому участнику, ход невозможен.
	if checkPieceColor(*b, pieceTo) {
		log.Printf("%v", errClashWithPieceOfSameColor)
		return isOk, nil

		// Если фигура в to принадлежит сопернику, проверка, возможно ли взятие
	} else {
		// ни одна фигура не может взять короля
		switch pieceTo {
		case board.WhiteKing, board.BlackKing:
			log.Printf("%v", errClashWithKing)
			return isOk, nil
		}

		// пешка не может взять ни одну фигуру при движении вертикально
		switch pieceFrom {
		case board.WhitePawn, board.BlackPawn:
			if from.diffColumn(to) == 0 {
				log.Printf("%v", errClashWithPawn)
				return isOk, nil
			}
		}
	}

	isOk = true
	return isOk, nil
}

// checkCastling валидирует рокировку: проверка, что при рокировке король не проходит через битое поле, рокировка
// разрешена (король и ладья еще не двигались, между ними все клетки пустые). НО: Поле, на котором король оказывается
// после рокировки проверяется на наличие шаха в другой функции. Если при обработке возникает ошибка, возвращает ее
// (иначе nil).
func checkCastling(b *board.Board, piece board.Piece, from, to square) (isValid bool, err error) {
	var rookSquare square
	var castling board.Castling
	switch piece {
	case board.WhiteKing:
		switch to.column {
		case 2:
			rookSquare = newSquare(0)
			castling = board.WhiteQueenside
		case 6:
			rookSquare = newSquare(7)
			castling = board.WhiteKingside
		}
	case board.BlackKing:
		switch to.column {
		case 2:
			rookSquare = newSquare(56)
			castling = board.BlackQueenside
		case 6:
			rookSquare = newSquare(63)
			castling = board.BlackKingside
		}
	}

	// 1. Проверка, что указанная рокировка (castling) валидна (ни король, ни ладья еще не двигались).
	if !b.HaveCastling(castling) {
		return isValid, nil
	}

	// 2. проверка, что между клеткой короля (from) и клеткой ладьи (rookSquare) все клетки пустые.
	squaresBetweenKingAndRook := getSquaresToBePassed(board.WhiteRook, from, rookSquare)
	for _, sq := range squaresBetweenKingAndRook {
		var pc board.Piece
		pc, err = b.Get(board.Sq(sq.toInt()))
		if err != nil {
			return isValid, err
		}
		if pc != 0 {
			log.Printf("%v", errCastlingThroughOccupiedSquare)
			return isValid, nil
		}
	}

	// 3. проверка, что король не проходит через битое поле (под шахом).
	squareToBePassed := newSquare(from.toInt8() + ((to.toInt8() - from.toInt8()) / 2))
	if isSquareChecked(*b, board.Sq(squareToBePassed.toInt()), piece == board.WhiteKing) {
		log.Printf("%v", errCastlingThroughCheckedSquare)
		return isValid, nil
	}

	isValid = true
	return isValid, nil
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
		var castling board.Castling
		if to.column == 6 {
			switch piece {
			case board.WhiteKing:
				castling = board.WhiteKingside
			case board.BlackKing:
				castling = board.BlackKingside
			}

			// сторона королевы
		} else if to.column == 2 {
			switch piece {
			case board.WhiteKing:
				castling = board.WhiteQueenside
			case board.BlackKing:
				castling = board.BlackQueenside
			}
		}

		err := b.Castle(castling)
		if err != nil {
			return b, err
		}
		return b, nil
	}

	// обработка всех остальных ходов
	err := b.Move(board.Sq(from.toInt()), board.Sq(to.toInt()))
	if err != nil {
		return b, err
	}
	return b, nil
}

// getSquareByPiece возвращает клетку, на которой находится заданная фигура. Если такой фигуры на доске нет,
// возвращает ошибку. Если таких фигур несколько, возвращает первую встретившуюся по пути с верхнего ряда к нижнему,
// слево направо.
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
		boardRow = boardRow[slashIndex+1:]
		rowsCount--
	}
	slashIndex := strings.Index(boardRow, "/")
	if slashIndex != -1 {
		boardRow = boardRow[:slashIndex]
	}

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

// isSquareChecked проверяет, находится ли поле под боем.
func isSquareChecked(b board.Board, s board.Square, weAreWhite bool) bool {
	p := board.BlackPawn
	if weAreWhite {
		p = board.WhitePawn
	}
	_ = b.Put(s, p)
	return len(position.ThreatsTo(s, b)) > 0
}

func getAvailableMoves(b board.Board, from square) (moves []square, err error) {
	var allMoves []square
	allMoves, err = getMoves(b, from)
	if err != nil {
		return moves, err
	}

	var piece board.Piece
	piece, err = b.Get(board.Sq(from.toInt()))

	for _, to := range allMoves {
		var isValid bool
		var newpiece board.Piece
		switch piece {
		case board.WhitePawn:
			newpiece = board.WhiteQueen
		case board.BlackPawn:
			newpiece = board.BlackQueen
		}

		_, isValid, err = advancedLogic(b, from, to, newpiece)
		if err != nil {
			return moves, err
		}
		if isValid {
			moves = append(moves, to)
		}
	}

	return moves, err
}
