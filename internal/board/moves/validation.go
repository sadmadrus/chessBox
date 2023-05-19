// Пакет moves валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package moves

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/board/validation"
)

// IsValid валидирует ход по позиции на доске b, клетке откуда делается ход from, клетке куда делается ход to и новой
// фигуре (в случае проведения пешки) promoteTo. Возвращает флаг валидации (true - ход валиден, false - ход невалиден).
// При некорректных входных данных возвращает флаг false и ошибку. Некорректными входными данными считаются:
// 1. доска b невалидна (ErrBoardNotValid).
// 2. клетка from или to невалидна (ErrSquareNotExist), либо клетки совпадают (ErrFromToSquaresNotDiffer).
// 3. в клетке from нет фигуры (ErrNoPieceOnFromSquare) либо фигура неправильного цвета (ErrPieceWrongColor).
// 4. указанная фигура promoteTo невалидна: то есть не конь, не слон, не ладья, не ферзь; или неправильного цвета; или
// указана фигура для проведения пешки, хотя в клетке from находится не пешка (ErrPromoteToNotValid).
// Если входные данные корректны, возвращается nil.
func IsValid(b board.Board, from, to board.Square, promoteTo board.Piece) (bool, error) {
	err := isDataValid(b, from, to, promoteTo)
	if err != nil {
		return false, err
	}

	fromSquare := newSquare(int8(from))
	toSquare := newSquare(int8(to))

	ok := validationLogic(b, fromSquare, toSquare, promoteTo)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// GetAvailableMoves по текущей позиции на доске b и клетке, с которой делается ход from определяет срез допустимых
// клеток, куда можно передвинуть данную фигуру moves. Срез клеток возвращается отсортированным по возрастанию.
// Пустой срез означает, что либо клетка пустая, либо на клетке стоит фигура, которой не принадлежит ход, либо на клетке
// стоит фигура правильного цвета, для которой нет разрешенных ходов.
// Если в входные данные некорректны, возвращается ошибка, иначе возвращается nil. Некорректными считаются данные:
// 1. доска b невалидна (ErrBoardNotValid).
// 2. клетка from невалидна (ErrSquareNotExist).
func GetAvailableMoves(b board.Board, from board.Square) ([]board.Square, error) {
	// 1. валидность доски b.
	if !validation.IsLegal(b) {
		return nil, ErrBoardNotValid
	}

	// 2. валидность клетки from.
	if !from.IsValid() {
		return nil, ErrSquareNotExist
	}

	fromSquare := newSquare(int8(from))
	allMoves, _ := getMoves(b, fromSquare)

	var moves []board.Square
	var piece board.Piece
	piece, _ = b.Get(from)

	for _, to := range allMoves {
		var promoteTo board.Piece
		if piece == board.WhitePawn && to.row == 7 {
			promoteTo = board.WhiteQueen
		} else if piece == board.BlackPawn && to.row == 0 {
			promoteTo = board.BlackQueen
		}

		isValid, _ := IsValid(b, from, board.Sq(to.toInt()), promoteTo)
		if isValid {
			moves = append(moves, board.Sq(to.toInt()))
		}
	}

	return moves, nil
}

// isDataValid проверяет корректность входных данных и возвращает ошибку или nil. Некорректными входными данными считаются:
// 1. доска b невалидна (ErrBoardNotValid).
// 2. клетка from или to невалидна (ErrSquareNotExist), либо клетки совпадают (ErrFromToSquaresNotDiffer).
// 3. в клетке from нет фигуры (ErrNoPieceOnFromSquare) либо фигура неправильного цвета (ErrPieceWrongColor).
// 4. указанная фигура promoteTo невалидна: то есть не конь, не слон, не ладья, не ферзь; или неправильного цвета; или
// указана фигура для проведения пешки, хотя в клетке from находится не пешка (ErrPromoteToNotValid).
func isDataValid(b board.Board, from, to board.Square, promoteTo board.Piece) error {
	// 1. валидность доски b.
	if !validation.IsLegal(b) {
		return ErrBoardNotValid
	}

	// 2. валидность клеток from, to.
	if !from.IsValid() || !to.IsValid() {
		return ErrSquareNotExist
	}
	if int(from) == int(to) {
		return ErrFromToSquaresNotDiffer
	}

	// 3a. если в клетке from фигуры нет, ход невалиден.
	var piece board.Piece
	piece, _ = b.Get(from)
	if piece == 0 {
		return ErrNoPieceOnFromSquare
	}

	// 3b. если фигура не принадлежит той стороне, чья очередь хода, ход невалиден.
	isOk := checkPieceColor(b, piece)
	if !isOk {
		return ErrPieceWrongColor
	}

	// 4. указанная фигура promoteTo невалидна.
	isOk = checkPawnPromotion(piece, newSquare(int8(to)), promoteTo)
	if !isOk {
		return ErrPromoteToNotValid
	}

	return nil
}

// validationLogic обрабывает общую логику валидации хода. Возвращает флаг валидации хода (true валиден, false нет).
func validationLogic(b board.Board, from, to square, promoteTo board.Piece) bool {
	// Логика валидации хода пошагово.
	piece, _ := b.Get(board.Sq(from.toInt()))

	// 1. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	err := move(piece, from, to)
	if err != nil {
		return false
	}

	// 2. Проверяем, что по пути фигуры с клетки from до to (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	squaresToBePassed := getSquaresToBePassed(piece, from, to)
	var areSquaresEmpty bool
	areSquaresEmpty, err = checkSquaresToBePassed(b, squaresToBePassed)
	if err != nil || !areSquaresEmpty {
		return false
	}

	// 3. Проверяем наличие и цвет фигур в клетке to. Обработка корректного взятия пешки на проходе (en Passant)
	// валидирована здесь.
	var isOk bool
	isOk, err = checkToSquare(&b, piece, from, to)
	if err != nil || !isOk {
		return false
	}

	// 4. Валидация рокировки: проверка, что при рокировке король не проходит через битое поле, рокировка разрешена
	// (король и ладья еще не двигались, между ними все клетки пустые).
	if (piece == board.WhiteKing || piece == board.BlackKing) && (abs(from.diffColumn(to)) == 2) {
		isOk, err = checkCastling(&b, piece, from, to)
		if err != nil || !isOk {
			return false
		}
	}

	// 5. На текущем этапе ход возможен. Генерируем новое положение доски newBoard. Так как до текущего положения ход
	// валидирован, ошибок не ожидаем.
	// Проверяем, что при новой позиции на доске не появился шах для собственного короля. На этом шаге также
	// проверяем, что король не находится вплотную к чужому королю - такой ход будет запрещен.
	var newBoard board.Board
	newBoard, _ = getNewBoard(b, piece, from, to, promoteTo)
	king := "k"
	if b.NextToMove() {
		king = "K"
	}
	kingSquare, _ := getSquareByPiece(newBoard, king)
	checks := validation.CheckedBy(board.Sq(kingSquare.toInt()), newBoard)

	return len(checks) == 0
}

// checkPawnPromotion проверяет, что указанная пользователем новая фигура promoteTo корректна относительно хода (то есть
// для проведения пешки указаны конь или слон или ладья или ферзь нужного цвета; а если пешка не проводится, указан 0).
func checkPawnPromotion(piece board.Piece, to square, promoteTo board.Piece) (isOk bool) {
	if piece == board.WhitePawn && to.row == 7 {
		switch promoteTo {
		case board.WhiteKnight, board.WhiteBishop, board.WhiteRook, board.WhiteQueen:
			isOk = true
		}

	} else if piece == board.BlackPawn && to.row == 0 {
		switch promoteTo {
		case board.BlackKnight, board.BlackBishop, board.BlackRook, board.BlackQueen:
			isOk = true
		}

	} else {
		if promoteTo == 0 {
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
	return len(validation.CheckedBy(s, b)) > 0
}
