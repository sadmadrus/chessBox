// пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package chessBox

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
)

var (
	errPieceNotExist          = fmt.Errorf("piece does not exist")
	errInvalidHttpMethod      = fmt.Errorf("method is not supported")
	errFromToSquaresNotDiffer = fmt.Errorf("from and to squares are not different")
	errPawnMoveNotValid       = fmt.Errorf("pawn move is not valid")
	errKnightMoveNotValid     = fmt.Errorf("knight move is not valid")
	errBishopMoveNotValid     = fmt.Errorf("bishop move is not valid")
	errRookMoveNotValid       = fmt.Errorf("rook move is not valid")
	errQueenMoveNotValid      = fmt.Errorf("queen move is not valid")
	errKingMoveNotValid       = fmt.Errorf("king move is not valid")
)

// http хендлеры

// SimpleValidation сервис отвечает за простую валидацию хода по начальной (from) и конечной (to) клетке
// и фигуре (piece) (GET, HEAD). Валидирует корректность геометрического перемещения фигуры без привязки к положению
// на доске. Возвращает заголовок HttpResponse 200 (ход валиден) или HttpsResponse 403 (ход невалиден). Возвращает
// HttpResponse 400 при некорректном методе запроса и некорректных входных данных.
// Входящие URL параметры:
// * фигура piece (k/q/r/b/n/p/K/Q/R/B/N/P)
// * начальная клетка предполагаемого хода from (число от 0 до 63, либо строка вида a1, c7 и т.п)
// * конечная клетка предполагаемого хода to (число от 0 до 63, либо строка вида a1, c7 и т.п)
func SimpleValidation(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" || r.Method == "HEAD" {
		// валидация входных данных: фигура piece и клетки from, to существуют
		pieceParsed := r.URL.Query().Get("piece")
		// TODO: дописать перевод из k/q/r/b/n/p/K/Q/R/B/N/P в int константу для фигур
		piece, err := someFunctionThatConvertsPieceLetterToInt(pieceParsed)
		if err != nil {
			log.Errorf("%w: %v", errPieceNotExist, pieceParsed)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if piece < WhitePawn || piece > BlackKing {
			log.Errorf("%w: %v", errPieceNotExist, piece)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fromParsed := r.URL.Query().Get("from")
		from := Sq(fromParsed)
		if from == -1 {
			log.Errorf("%w: %v (from)", errSquareNotExist, fromParsed)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		toParsed := r.URL.Query().Get("to")
		to := Sq(toParsed)
		if to == -1 {
			log.Errorf("%w: %v (to)", errSquareNotExist, toParsed)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if from == to {
			log.Errorf("%w: %v (from), %v (to)", errFromToSquaresNotDiffer, from, to)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация геометрического движения фигуры без привязки к позиции на доске
		fromRow := float64(from / 8)
		fromColumn := float64(from % 8)
		toRow := float64(to / 8)
		toColumn := float64(to % 8)

		switch piece {
		case WhitePawn, BlackPawn:
			err = MovePawn(piece, fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case WhiteKnight, BlackKnight:
			err = MoveKnight(fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case WhiteBishop, BlackBishop:
			err = MoveBishop(fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case WhiteRook, BlackRook:
			err = MoveRook(fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case WhiteQueen, BlackQueen:
			err = MoveQueen(fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case WhiteKing, BlackKing:
			err = MoveKing(piece, fromRow, fromColumn, toRow, toColumn)
			if err != nil {
				log.Errorf("%w: from %v - to %v", err, from, to)
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		default: // если число piece не целое (можно вынести выше в валидацию входных данных)
			log.Errorf("%w: %v", errPieceNotExist, piece)
			w.WriteHeader(http.StatusBadRequest)
		}

	}

	log.Errorf("inside SimpleValidation %w: %v", errInvalidHttpMethod, r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

// Методы Move для каждого типа фигуры. В случае невозможности сделать ход, возвращают ошибку, иначе возвращают nil.

// MovePawn логика движения пешки. Может двигаться вверх (белый) или вниз (черный) на 1 или 2 клетки. Может съедать
// по диагонали проходные пешки или фигуры соперника. Возвращает ошибку, если движение невалидно.
func MovePawn(piece int, fromRow, fromColumn, toRow, toColumn float64) error {
	var isVerticalValid, isDiagonalValid bool

	switch piece {
	case WhitePawn:
		isVerticalValid = (fromColumn-toColumn == 0) && (fromRow != 0) && (toRow > 1) &&
			((toRow-fromRow == 1) || (toRow == 3 && fromRow == 1))
		isDiagonalValid = (fromRow != 0) && (toRow > 1) && (toRow-fromRow == 1) && (toColumn-fromColumn == 1)
	case BlackPawn:
		isVerticalValid = (fromColumn-toColumn == 0) && (fromRow != 7) && (toRow < 6) &&
			((toRow-fromRow == -1) || (toRow == 4 && fromRow == 6))
		isDiagonalValid = (fromRow != 7) && (toRow < 6) && (toRow-fromRow == -1) && (toColumn-fromColumn == -1)
	}

	if isVerticalValid || isDiagonalValid {
		return nil
	}
	return fmt.Errorf("%w", errPawnMoveNotValid)
}

// MoveKnight логика движения коня без привязки к позиции на доске. Может двигаться буквой Г. То есть +/- 2 клетки
// в одном направлении и +/- 1 клетка в перпендикулярном направлении. Возвращает ошибку, если движение невалидно.
func MoveKnight(fromRow, fromColumn, toRow, toColumn float64) error {
	isValid := (math.Abs(fromRow-toRow) == 2 && math.Abs(fromColumn-toColumn) == 1) ||
		(math.Abs(fromRow-toRow) == 1 && math.Abs(fromColumn-toColumn) == 2)

	if !isValid {
		return fmt.Errorf("%w", errKnightMoveNotValid)
	}
	return nil
}

// MoveBishop логика движения слона без привязки к позиции на доске. Может двигаться по всем диагоналям. Возвращает
// ошибку, если движение невалидно.
func MoveBishop(fromRow, fromColumn, toRow, toColumn float64) error {
	isValid := math.Abs(fromRow-toRow) == math.Abs(fromColumn-toColumn)

	if !isValid {
		return fmt.Errorf("%w", errBishopMoveNotValid)
	}
	return nil
}

// MoveRook логика движения ладьи. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток. Возвращает
// ошибку, если движение невалидно.
func MoveRook(fromRow, fromColumn, toRow, toColumn float64) error {
	isValid := (fromRow-toRow == 0) || (fromColumn-toColumn == 0)

	if !isValid {
		return fmt.Errorf("%w", errRookMoveNotValid)
	}
	return nil
}

// MoveQueen для ферзя. Может двигаться вверх, вниз, влево, вправо на любое кол-во клеток. Может двигаться диагонально
// на любое кол-во клеток. Возвращает ошибку, если движение невалидно.
func MoveQueen(fromRow, fromColumn, toRow, toColumn float64) error {
	errBishop := MoveBishop(fromRow, fromColumn, toRow, toColumn)
	errRook := MoveRook(fromRow, fromColumn, toRow, toColumn)

	if errBishop != nil || errRook != nil {
		return fmt.Errorf("%w", errQueenMoveNotValid)
	}
	return nil
}

// MoveKing логика движения короля. Может двигаться вертикально, горизонтально и диагонально только на одну клетку.
// Также король из своего начального положения на доске (row 0 && column 4 для белого; row 7 && column 4 для черного)
// может двигаться: на 2 клетки вправо или 2 клетки влево для рокировок.
func MoveKing(piece int, fromRow, fromColumn, toRow, toColumn float64) error {
	isHorizontalValid := (fromRow-toRow == 0) && (math.Abs(fromColumn-toColumn) == 1)
	isVerticalValid := (math.Abs(fromRow-toRow) == 1) && (fromColumn-toColumn == 0)
	isDiagonalValid := (math.Abs(fromRow-toRow) == 1) && (math.Abs(fromColumn-toColumn) == 1)
	isCastlingValid := false

	switch piece {
	case WhiteKing:
		if fromRow == 0 && fromColumn == 4 && toRow == 0 && (math.Abs(fromColumn-toColumn) == 2) {
			isCastlingValid = true
		}
	case BlackKing:
		if fromRow == 7 && fromColumn == 4 && toRow == 7 && (math.Abs(fromColumn-toColumn) == 2) {
			isCastlingValid = true
		}
	}

	if isHorizontalValid || isVerticalValid || isDiagonalValid || isCastlingValid {
		return nil
	}

	return fmt.Errorf("%w", errKingMoveNotValid)
}

// http хендлеры

// AdvancedValidation сервис отвечает за сложную валидацию хода по начальной и конечной клетке, а также по текущему состоянию
// доски в нотации FEN. Также принимает на вход URL-параметр newpiece (это новая фигура, в которую нужно превратить
// пешку при достижении последнего ряда, в формате pieceВозвращает заголовок HttpResponse 200 (ход валиден) или
// HttpsResponse 403 (ход невалиден). Возвращает в теле JSON с конечной доской board в форате FEN.
func AdvancedValidation(w http.ResponseWriter, r *http.Request) {

}

// GetAvailableMoves сервис отвечает за оплучение всех возможных ходов для данной позиции доски в нотации FEN и начальной клетке.
// Возвращает заголовок HttpResponse 200 (в случае непустого массива клеток) или HttpsResponse 403 (клетка пустая или
// с фигурой, которой не принадлежит ход или массив клеток пуст). Возвращает в теле
// JSON массив всех клеток, движение на которые валидно для данной фигуры.
func GetAvailableMoves(w http.ResponseWriter, r *http.Request) {

}

// Вспомогательные функции

// ValidateMove обрабывает общую логику валидации хода
func ValidateMove(b Board, from, to s) error {
	// Логика валидации хода пошагово.

	// 1. получаем фигуру, находящуюся в клетке startCell. Если в этой клетке фигуры нет, возвращаем ошибку
	figure := b.GetFigure(from)
	if figure == nil {
		return fmt.Errorf("move invalid: startCell %d doesn't have any figures", from)
	}

	// 2. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(figure)
	if !isFigureRightColor {
		return fmt.Errorf("move invalid: startCell %d has figure of wrong color", from)
	}

	// 3. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	var cellsToBePassed []square
	cellsToBePassed, err := checkFigureMove(figure, from, to)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// 4. Проверяем, что по пути фигуры с клетки startCell до endCell (не включительно) нет других фигур
	// (f.e. слон а1-h8, но на b2 стоит конь - так запрещено).
	if len(cellsToBePassed) > 0 {
		err = checkCellsToBePassed(b, cellsToBePassed)
		if err != nil {
			return fmt.Errorf("move invalid: %w", err)
		}
	}

	// 5. Проверяем наличие и цвет фигур в клетке endCell.
	err = checkEndCell(b, to)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// TODO: проверяем, что пользователь не захотел выставить нового ферзя в неуместном для этого случае.

	// 6. На текущем этапе ход возможен. Генерируем новое положение доски newBoard.
	newBoard := b.GenerateBoardAfterMove(from, to)

	// 7. Проверяем, что при новой позиции на доске не появился шах для собственного короля оппонента o.
	err = checkSelfCheck(newBoard)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// 8. В случае если ход делается королем, проверяем, что он не подступил вплотную к чужому королю - такой ход
	// будет запрещен.

	// 9. В случае если ход делается пешкой, проверяем, попала ли пешка на последнюю линию и потребуется ли
	// трансформация в другую фигуру. Если да, выставляем pawnTransformation = true

	return nil
}

// checkFigureColor проверяет, что очередь хода и цвет фигуры f, которую хотят передвинуть, совпадают.
// Возвращает true в случае успеха, false в противном случае.
func checkFigureColor(b Board) bool {

}

// checkFigureMove проверяет, что фигура f может двигаться в этом направлении с клетки startCell на клетку
// endCell. Возвращает ошибку, если движение невозможно, или nil в противном случае. Также возвращают массив
// cellsToBePassed []Cell, в который входят все "промежуточные" клетки, которые "проходит" эта фигура на своем
// пути с startCell до endCell, если такие есть. В противном случае возвращает nil.
func checkFigureMove(f Figure, from, to s) (cellsToBePassed []Cell, err error) {
	// через цепочку if идет проверка на тип Figure (или на поле name в структуре Figure) и перенаправление на
	// соответствующий метод Move для этого типа фигуры
}

// checkCellsToBePassed проверяет, есть ли на клетках из массива cellsToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается ошибка. Иначе возвращается nil.
func checkCellsToBePassed(b Board, cellsToBePassed []Cell) error {
	for _, cell := range cellsToBePassed {
		f := b.GetFigure(cell)
		if f != nil {
			return fmt.Errorf("figure present on cell %d", cell)
		}
	}
	return nil
}

// checkEndCell проверяет наличие фигуры на клетке endCell на предмет совместимости хода. Возвращает ошибку при
// несовместимости хода или nil в случае успеха.
func checkEndCell(b *Board, to s) error {
	// логика пошагово:

	// 1. Если в клетке endCell нет фигур, ход возможен:
	f := b.GetFigure(endCell)
	if f == nil {
		return nil
	}

	// 2. Если фигура в endCell принадлежит самому участнику o, ход невозможен.
	//

	// 3. Если фигура в endCell принадлежит сопернику, проверка, возможно ли взятие
	// (f.e. пешка e2-e4 не может взять коня на e4).
	//
}

// checkSelfCheck проверяет, что при новой позиции на доске нет шаха собственному королю.
func checkSelfCheck(b Board) error {
	// Находим клетку с собственным королем.
	kingCell := getKingCell(b)

	// проверяем, есть ли на расстоянии буквы Г от этой клетки вражеские кони. Если да, выдаем ошибку.
	err := checkEnemyKnightsNearKing(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy knight")
	}

	// проверяем, есть ли по вертикали или горизонтали в качестве ближайших фигур вражеские ладьи и ферзи.
	// Если да, выдаем ошибку.
	err = checkEnemiesVerticallyAndHorizontally(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy rook or queen")
	}

	// проверяем, есть ли по диагоналям в качестве ближайших фигур вражеские ферзи, слоны и пешки.
	// Если да, выдаем ошибку.
	err = checkEnemiesDiagonally(b, kingCell)
	if err != nil {
		return fmt.Errorf("self-check by enemy queen or pawn or bishop")
	}

	return nil

}

// getKingCell возвращает клетку, на которой находится свой король оппонента.
func getKingCell(b Board) s {

}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemyKnightsNearKing(b Board, kingCell Cell) error {

}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingCell, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b Board, kingCell Cell) error {

}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingCell, на
// которых находятся вражеские слоны, ферзи и пешки. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesDiagonally(b Board, kingCell Cell) error {
	// должна быть реализована дополнтельная проверка на пешки - их битое поле только по ходу движения!
}
