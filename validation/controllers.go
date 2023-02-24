// Пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sadmadrus/chessBox/internal/board"
	log "github.com/sirupsen/logrus"
)

var (
	errPieceNotExist                = fmt.Errorf("piece does not exist")
	errInvalidHttpMethod            = fmt.Errorf("method is not supported")
	errFromToSquaresNotDiffer       = fmt.Errorf("from and to squares are not different")
	errPawnMoveNotValid             = fmt.Errorf("pawn move is not valid")
	errKnightMoveNotValid           = fmt.Errorf("knight move is not valid")
	errBishopMoveNotValid           = fmt.Errorf("bishop move is not valid")
	errRookMoveNotValid             = fmt.Errorf("rook move is not valid")
	errQueenMoveNotValid            = fmt.Errorf("queen move is not valid")
	errKingMoveNotValid             = fmt.Errorf("king move is not valid")
	errNoPieceOnFromSquare          = fmt.Errorf("no piece on from square")
	errPieceWrongColor              = fmt.Errorf("piece has wrong color")
	errPieceFound                   = fmt.Errorf("piece found on the square")
	errClashWithPieceOfSameColor    = fmt.Errorf("clash with piece of the same color")
	errClashWithKing                = fmt.Errorf("clash with king")
	errClashWithPawn                = fmt.Errorf("pawn can not clash with another piece when moving vertically")
	errPawnPromotionNotValid        = fmt.Errorf("pawn promotion to pawn or king is not valid")
	errNewpieceExist                = fmt.Errorf("newpiece exists with no pawn promotion")
	errNewpieceNotExist             = fmt.Errorf("newpiece does not exist but pawn promotion required")
	errPieceNotExistOnBoard         = fmt.Errorf("piece does not exist on board")
	errKingChecked                  = fmt.Errorf("king checked after move")
	errKingsAdjacent                = fmt.Errorf("kings are adjacent")
	errCastlingThroughCheckedSquare = fmt.Errorf("castling is not valid through square under check")
)

// http хендлеры

// Simple сервис отвечает за простую валидацию хода по начальной (from) и конечной (to) клетке
// и фигуре (piece) (GET, HEAD). Валидирует корректность геометрического перемещения фигуры без привязки к положению
// на доске. Возвращает заголовок HttpResponse 200 (ход валиден) или HttpsResponse 403 (ход невалиден). Возвращает
// HttpResponse 400 при некорректном методе запроса и некорректных входных данных.
// Входящие URL параметры:
// * фигура piece (k/q/r/b/n/p/K/Q/R/B/N/P)
// * начальная клетка предполагаемого хода from (число от 0 до 63, либо строка вида a1, c7 и т.п)
// * конечная клетка предполагаемого хода to (число от 0 до 63, либо строка вида a1, c7 и т.п).
func Simple(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" || r.Method == "HEAD" {
		// валидация входных данных: фигура piece существует
		pieceParsed := r.URL.Query().Get("piece")
		piece, err := parsePieceFromLetter(pieceParsed)
		if err != nil {
			log.Errorf("%v: %v", errPieceNotExist, pieceParsed)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: клетка from существуют
		fromParsed := r.URL.Query().Get("from")
		// перевод в тип board.square для цифро-буквенного обозначения клетки (напр., "а1")
		from := board.Sq(fromParsed)
		if from == -1 {
			// перевод в тип board.square для числового значения клетки от 0 до 63
			var fromParsedNum int
			fromParsedNum, err = strconv.Atoi(fromParsed)
			from = board.Sq(fromParsedNum)
			if from == -1 || err != nil {
				log.Errorf("%v: %v", errPieceNotExist, fromParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// валидация входных данных: клетка to существуют
		toParsed := r.URL.Query().Get("to")
		// перевод в тип board.square для цифро-буквенного обозначения клетки (напр., "а1")
		to := board.Sq(toParsed)
		if to == -1 {
			// перевод в тип board.square для числового значения клетки от 0 до 63
			var toParsedNum int
			toParsedNum, err = strconv.Atoi(toParsed)
			to = board.Sq(toParsedNum)
			if to == -1 || err != nil {
				log.Errorf("%v: %s", errPieceNotExist, toParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// валидация входных данных: клетки from и to различны
		if from == to {
			log.Errorf("%v: %v (from), %v (to)", errFromToSquaresNotDiffer, from, to)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация геометрического движения фигуры без привязки к позиции на доске
		fromSquare := newSquare(int8(from))
		toSquare := newSquare(int8(to))
		err = move(piece, fromSquare, toSquare)
		if err != nil {
			log.Errorf("%v: from %v - to %v", err, from, to)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	log.Errorf("inside Simple %v: %v", errInvalidHttpMethod, r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

type advancedResponse struct {
	Board string `json:"board"`
}

// Advanced сервис отвечает за сложную валидацию хода по начальной и конечной клетке, а также по текущему состоянию
// доски в нотации FEN. Также принимает на вход URL-параметр newpiece (это новая фигура, в которую нужно превратить
// пешку при достижении последнего ряда), в формате pieceВозвращает заголовок HttpResponse 200 (ход валиден) или
// HttpsResponse 403 (ход невалиден). Возвращает HttpResponse 400 при некорректном методе запроса и некорректных
// входных данных. Возвращает в теле JSON с конечной доской board в форате FEN.
// Входящие URL параметры:
// * доска board в формате UsFen (например, "rnbqkbnr~pppppppp~8~8~8~8~PPPPPPPP~RNBQKBNR+w+KQkq+-+0+1")
// * начальная клетка предполагаемого хода from (число от 0 до 63, либо строка вида a1, c7 и т.п)
// * конечная клетка предполагаемого хода to (число от 0 до 63, либо строка вида a1, c7 и т.п).
// * фигура newpiece (q/r/b/n/Q/R/B/N или пустое значение)
func Advanced(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		// валидация входных данных: доска board существует
		// TODO board нужно ли проверять валидность доски где-то вообще?
		boardParsed := r.URL.Query().Get("board")
		b, err := board.FromUsFEN(boardParsed)
		if err != nil {
			log.Errorf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: клетка from существуют
		fromParsed := r.URL.Query().Get("from")
		// перевод в тип board.square для цифро-буквенного обозначения клетки (напр., "а1")
		from := board.Sq(fromParsed)
		if from == -1 {
			// перевод в тип board.square для числового значения клетки от 0 до 63
			var fromParsedNum int
			fromParsedNum, err = strconv.Atoi(fromParsed)
			from = board.Sq(fromParsedNum)
			if from == -1 || err != nil {
				log.Errorf("%v: %v", errPieceNotExist, fromParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// валидация входных данных: клетка to существуют
		toParsed := r.URL.Query().Get("to")
		// перевод в тип board.square для цифро-буквенного обозначения клетки (напр., "а1")
		to := board.Sq(toParsed)
		if to == -1 {
			// перевод в тип board.square для числового значения клетки от 0 до 63
			var toParsedNum int
			toParsedNum, err = strconv.Atoi(toParsed)
			to = board.Sq(toParsedNum)
			if to == -1 || err != nil {
				log.Errorf("%v: %s", errPieceNotExist, toParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// валидация входных данных: клетки from и to различны
		if from == to {
			log.Errorf("%v: %v (from), %v (to)", errFromToSquaresNotDiffer, from, to)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: фигура newpiece принимает q/r/b/n/Q/R/B/N или пустое значение
		newpieceParsed := r.URL.Query().Get("newpiece")
		var newpiece board.Piece
		if newpieceParsed != "" {
			newpiece, err = parsePieceFromLetter(newpieceParsed)
			if err != nil {
				log.Errorf("%v: %v", errPieceNotExist, newpieceParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			switch newpiece {
			case board.WhitePawn, board.BlackPawn, board.WhiteKing, board.BlackKing:
				log.Errorf("%v: %v", errPawnPromotionNotValid, newpieceParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// полная валидация хода с учетом положения на доске, а также возможных рокировок, взятия на проходе и
		// проведения пешки
		fromSquare := newSquare(int8(from))
		toSquare := newSquare(int8(to))
		newBoard, err := advancedValidationLogic(*b, fromSquare, toSquare, newpiece)
		if err != nil {
			log.Errorf("move invalid: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			boardUsFEN := newBoard.UsFEN()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			data := advancedResponse{boardUsFEN}
			err = json.NewEncoder(w).Encode(data)
			if err != nil {
				log.Errorf("error while encoding json: %v", err)
			}
			return
		}
	}

	log.Errorf("inside Advanced %v: %v", errInvalidHttpMethod, r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

// AvailableMoves сервис отвечает за оплучение всех возможных ходов для данной позиции доски в нотации FEN и начальной клетке.
// Возвращает заголовок HttpResponse 200 (в случае непустого массива клеток) или HttpsResponse 403 (клетка пустая или
// с фигурой, которой не принадлежит ход или массив клеток пуст). Возвращает HttpResponse 400 при некорректном методе
// запроса и некорректных входных данных. Возвращает в теле JSON массив всех клеток, движение на которые валидно
// для данной фигуры.
func AvailableMoves(w http.ResponseWriter, r *http.Request) {
	// TODO написать логику
}

// Вспомогательные функции
// TODO по мере написания сервисов вспомогательные функции могут быть реорганизованы в другие файлы этого пакета для удобства!

// parsePieceFromLetter переводит строковое представление фигуры типа k/q/r/b/n/p/K/Q/R/B/N/P в тип board.Piece.  Если
// преобразование невозможно, возвращает ошибку.
// TODO add tests to all functions below
func parsePieceFromLetter(piece string) (board.Piece, error) {
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

// advancedValidationLogic обрабывает общую логику валидации хода.
func advancedValidationLogic(b board.Board, from, to square, newpiece board.Piece) (newBoard board.Board, err error) {
	// Логика валидации хода пошагово.
	// 1. получаем фигуру, находящуюся в клетке from. Если в этой клетке фигуры нет, возвращаем ошибку
	var piece board.Piece
	piece, err = b.Get(board.Sq(from.toInt()))
	if err != nil {
		return newBoard, err
	}
	if piece == 0 {
		return newBoard, fmt.Errorf("%v: %v", errNoPieceOnFromSquare, from)
	}

	// 2a. проверяем, что пользователь указал, какую новую фигуру выставить в случае проведения пешки.
	if newpiece == 0 && ((piece == board.WhitePawn && to.row == 7) || (piece == board.BlackPawn && to.row == 0)) {
		return newBoard, fmt.Errorf("%v", errNewpieceNotExist)
	}

	// 2b. проверяем, что пользователь не захотел выставить нового фигуру в неуместном для этого случае.
	if newpiece != 0 && ((piece != board.WhitePawn && to.row != 7) || (piece != board.BlackPawn && to.row != 0)) {
		return newBoard, fmt.Errorf("%v", errNewpieceExist)
	}

	// 3. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(b, piece)
	if !isFigureRightColor {
		return newBoard, fmt.Errorf("%v: %v", errPieceWrongColor, from)
	}

	// 4. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	err = move(piece, from, to)
	if err != nil {
		return newBoard, err
	}

	// 5. Проверяем, что по пути фигуры с клетки from до to (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	squaresToBePassed := getSquaresToBePassed(piece, from, to)
	if len(squaresToBePassed) > 0 {
		err = checkSquaresToBePassed(b, squaresToBePassed)
		if err != nil {
			return newBoard, err
		}
	}

	// 6. Проверяем наличие и цвет фигур в клетке to.
	err = checkToSquare(&b, piece, from, to)
	if err != nil {
		return newBoard, err
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
			return newBoard, err
		}
		if squareChecked {
			return newBoard, fmt.Errorf("%v", errCastlingThroughCheckedSquare)
		}
	}

	// 8. На текущем этапе ход возможен. Генерируем новое положение доски newBoard. Выдается ошибка при некорректном
	// проведении пешки, некорректной рокировке, некорректном взятии на проходе по логике из пакета board.
	newBoard, err = getNewBoard(b, piece, from, to, newpiece)
	if err != nil {
		return newBoard, fmt.Errorf("%w", err)
	}

	// 9. Проверяем, что при новой позиции на доске не появился шах для собственного короля.
	var kingChecked bool
	kingChecked, err = isKingChecked(b, king)
	if err != nil {
		return newBoard, err
	}
	if kingChecked {
		return newBoard, fmt.Errorf("%v", errKingChecked)
	}

	// 10. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.
	if piece == board.WhiteKing || piece == board.BlackKing {
		var isEnemyKingAdjacent bool
		isEnemyKingAdjacent, err = checkDistanceToEnemyKing(b)
		if err != nil {
			return newBoard, err
		}
		if isEnemyKingAdjacent {
			return newBoard, fmt.Errorf("%v", errKingsAdjacent)
		}
	}

	return newBoard, nil
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
	var movingUp, movingRight int8

	if from.diffRow(to) < 0 {
		movingUp = 1
	} else if from.diffRow(to) > 0 {
		movingUp = -1
	} else {
		movingUp = 0
	}

	if from.diffColumn(to) < 0 {
		movingRight = 1
	} else if from.diffColumn(to) > 0 {
		movingRight = -1
	} else {
		movingRight = 0
	}

	squaresToBePassedAmount := abs(from.diffRow(to))
	if abs(from.diffColumn(to)) > abs(from.diffRow(to)) {
		squaresToBePassedAmount = abs(from.diffColumn(to))
	}

	switch p {
	case board.WhitePawn, board.BlackPawn, board.WhiteBishop, board.BlackBishop, board.WhiteRook, board.BlackRook, board.WhiteKing, board.BlackKing, board.WhiteQueen, board.BlackQueen:
		for squaresToBePassedAmount > 1 {
			squareToBeAdded := newSquare(from.toInt8() + (movingUp * 8) + movingRight)
			squaresToBePassed = append(squaresToBePassed, squareToBeAdded)
			squaresToBePassedAmount--
		}
	}

	return squaresToBePassed
}

// checkSquaresToBePassed проверяет, есть ли на клетках из массива squaresToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается ошибка. Иначе возвращается nil.
func checkSquaresToBePassed(b board.Board, squaresToBePassed []square) error {
	for _, sq := range squaresToBePassed {
		piece, _ := b.Get(board.Sq(sq.toInt()))
		if piece == 0 {
			return fmt.Errorf("%v: %v", errPieceFound, sq)
		}
	}
	return nil
}

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
