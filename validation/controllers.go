// Пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package validation

import (
	"fmt"
	"github.com/sadmadrus/chessBox/internal/usfen"
	"net/http"
	"strconv"

	"github.com/sadmadrus/chessBox/internal/board"
	log "github.com/sirupsen/logrus"
)

var (
	errPieceNotExist             = fmt.Errorf("piece does not exist")
	errInvalidHttpMethod         = fmt.Errorf("method is not supported")
	errFromToSquaresNotDiffer    = fmt.Errorf("from and to squares are not different")
	errPawnMoveNotValid          = fmt.Errorf("pawn move is not valid")
	errKnightMoveNotValid        = fmt.Errorf("knight move is not valid")
	errBishopMoveNotValid        = fmt.Errorf("bishop move is not valid")
	errRookMoveNotValid          = fmt.Errorf("rook move is not valid")
	errQueenMoveNotValid         = fmt.Errorf("queen move is not valid")
	errKingMoveNotValid          = fmt.Errorf("king move is not valid")
	errNoPieceOnFromSquare       = fmt.Errorf("no piece on from square")
	errPieceWrongColor           = fmt.Errorf("piece has wrong color")
	errPieceFound                = fmt.Errorf("piece found on the square")
	errClashWithPieceOfSameColor = fmt.Errorf("clash with piece of the same color")
	errClashWithKing             = fmt.Errorf("clash with king")
	errClashWithPawn             = fmt.Errorf("pawn can not clash with another piece when moving vertically")
	errPawnPromotionNotValid     = fmt.Errorf("pawn promotion to pawn or king is not valid")
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
	// TODO написать логику
	if r.Method == "GET" || r.Method == "HEAD" {
		// валидация входных данных: доска board существует
		// TODO board нужно ли проверять валидность позиции?
		boardParsed := r.URL.Query().Get("board")
		boardFen := usfen.ToFen(boardParsed)
		b, err := board.FromFEN(boardFen)
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
		newBoard, err := advancedMoveValidation(*b, fromSquare, toSquare, newpiece)
		if err != nil {
			log.Errorf("move invalid: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			// TODO возвращать в теле новую доску
			w.WriteHeader(http.StatusOK)
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
// TODO add tests
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

// advancedMoveValidation обрабывает общую логику валидации хода.
func advancedMoveValidation(b board.Board, from, to square, newpiece board.Piece) (newBoard board.Board, err error) {
	// TODO написать логику
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

	// 2. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(b, piece)
	if !isFigureRightColor {
		return newBoard, fmt.Errorf("%v: %v", errPieceWrongColor, from)
	}

	// 3. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	err = move(piece, from, to)
	if err != nil {
		return newBoard, err
	}

	// 4. Проверяем, что по пути фигуры с клетки from до to (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	squaresToBePassed := getSquaresToBePassed(piece, from, to)
	if len(squaresToBePassed) > 0 {
		err = checkSquaresToBePassed(b, squaresToBePassed)
		if err != nil {
			return newBoard, err
		}
	}

	// 5. Проверяем наличие и цвет фигур в клетке to.
	err = checkToSquare(&b, piece, from, to)
	if err != nil {
		return newBoard, err
	}

	// TODO: остановилась здесь
	// TODO: проверяем, что пользователь не захотел выставить нового ферзя в неуместном для этого случае.

	// 6. На текущем этапе ход возможен. Генерируем новое положение доски newBoard.
	newBoard, err = getNewBoard(b, piece, from, to, newpiece)
	if err != nil {
		return newBoard, fmt.Errorf("move invalid: %w", err)
	}

	// 7. Проверяем, что при новой позиции на доске не появился шах для собственного короля оппонента o.
	err = checkSelfCheck(b)
	if err != nil {
		return newBoard, fmt.Errorf("move invalid: %w", err)
	}

	// 8. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.

	// 9. В случае если ход делается пешкой, проверяем, попала ли пешка на последнюю линию и потребуется ли
	// трансформация в другую фигуру.

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
	switch piece {
	case board.WhitePawn:
		if to.row == 7 {
			err := b.Promote(board.Sq(from.toInt()), board.Sq(to.toInt()), newpiece)
			if err != nil {
				return b, err
			}
		}

	case board.BlackPawn:
		if to.row == 0 {
			err := b.Promote(board.Sq(from.toInt()), board.Sq(to.toInt()), newpiece)
			if err != nil {
				return b, err
			}
		}

		// TODO рокировка
	case board.WhiteKing:

	}

}

// checkSelfCheck проверяет, что при новой позиции на доске нет шаха собственному королю.
func checkSelfCheck(b board.Board) error {
	// TODO написать логику

	// Находим клетку с собственным королем.
	kingSquare := getKingSquare(b)

	// проверяем, есть ли на расстоянии буквы Г от этой клетки вражеские кони. Если да, выдаем ошибку.
	err := checkEnemyKnightsNearKing(b, kingSquare)
	if err != nil {
		return fmt.Errorf("self-check by enemy knight")
	}

	// проверяем, есть ли по вертикали или горизонтали в качестве ближайших фигур вражеские ладьи и ферзи.
	// Если да, выдаем ошибку.
	err = checkEnemiesVerticallyAndHorizontally(b, kingSquare)
	if err != nil {
		return fmt.Errorf("self-check by enemy rook or queen")
	}

	// проверяем, есть ли по диагоналям в качестве ближайших фигур вражеские ферзи, слоны и пешки.
	// Если да, выдаем ошибку.
	err = checkEnemiesDiagonally(b, kingSquare)
	if err != nil {
		return fmt.Errorf("self-check by enemy queen or pawn or bishop")
	}

	return nil

}

// getKingSquare возвращает клетку, на которой находится свой король оппонента.
func getKingSquare(b board.Board) int {
	// TODO написать логику
	return 0
}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemyKnightsNearKing(b board.Board, kingSquare int) error {
	// TODO написать логику
	return nil
}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingSquare, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b board.Board, kingSquare int) error {
	// TODO написать логику
	return nil
}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingSquare, на
// которых находятся вражеские слоны, ферзи и пешки. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesDiagonally(b board.Board, kingSquare int) error {
	// должна быть реализована дополнтельная проверка на пешки - их битое поле только по ходу движения!
	// TODO написать логику
	return nil
}
