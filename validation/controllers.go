// пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре

package validation

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sadmadrus/chessBox/internal/board"
	log "github.com/sirupsen/logrus"
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

// Simple сервис отвечает за простую валидацию хода по начальной (from) и конечной (to) клетке
// и фигуре (piece) (GET, HEAD). Валидирует корректность геометрического перемещения фигуры без привязки к положению
// на доске. Возвращает заголовок HttpResponse 200 (ход валиден) или HttpsResponse 403 (ход невалиден). Возвращает
// HttpResponse 400 при некорректном методе запроса и некорректных входных данных.
// Входящие URL параметры:
// * фигура piece (k/q/r/b/n/p/K/Q/R/B/N/P)
// * начальная клетка предполагаемого хода from (число от 0 до 63, либо строка вида a1, c7 и т.п)
// * конечная клетка предполагаемого хода to (число от 0 до 63, либо строка вида a1, c7 и т.п)
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
				log.Errorf("ошибка при указании клетки: %s", fromParsed)
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
				log.Errorf("ошибка при указании клетки: %s", toParsed)
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
		toSquare := newSquare(int8(from))
		err = move(piece, fromSquare, toSquare)
		if err != nil {
			log.Errorf("%v: from %v - to %v", err, from, to)
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

	log.Errorf("inside Simple %v: %v", errInvalidHttpMethod, r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

// Advanced сервис отвечает за сложную валидацию хода по начальной и конечной клетке, а также по текущему состоянию
// доски в нотации FEN. Также принимает на вход URL-параметр newpiece (это новая фигура, в которую нужно превратить
// пешку при достижении последнего ряда, в формате pieceВозвращает заголовок HttpResponse 200 (ход валиден) или
// HttpsResponse 403 (ход невалиден). Возвращает в теле JSON с конечной доской board в форате FEN.
func Advanced(w http.ResponseWriter, r *http.Request) {
	// TODO написать логику
}

// AvailableMoves сервис отвечает за оплучение всех возможных ходов для данной позиции доски в нотации FEN и начальной клетке.
// Возвращает заголовок HttpResponse 200 (в случае непустого массива клеток) или HttpsResponse 403 (клетка пустая или
// с фигурой, которой не принадлежит ход или массив клеток пуст). Возвращает в теле
// JSON массив всех клеток, движение на которые валидно для данной фигуры.
func AvailableMoves(w http.ResponseWriter, r *http.Request) {
	// TODO написать логику
}

// Вспомогательные функции
// TODO по мере написания сервисов вспомогательные функции могут быть реорганизованы в другие файлы этого пакета для удобства!

// parsePieceFromLetter переводит строковое представление фигуры типа k/q/r/b/n/p/K/Q/R/B/N/P в тип board.Piece.  Если
// преобразование невозможно, возвращает ошибку.
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

// ValidateMove обрабывает общую логику валидации хода
func ValidateMove(b board.Board, from, to int) error {
	// TODO написать логику
	// Логика валидации хода пошагово.

	// 1. получаем фигуру, находящуюся в клетке startCell. Если в этой клетке фигуры нет, возвращаем ошибку
	figure, err := b.Get(board.Sq(from))
	if err != nil {
		return fmt.Errorf("move invalid: startCell %d doesn't have any figures", from)
	}

	// 2. проверяем, что фигура принадлежит той стороне, чья очередь хода
	isFigureRightColor := checkFigureColor(b)
	if !isFigureRightColor {
		return fmt.Errorf("move invalid: startCell %d has figure of wrong color", from)
	}

	// 3. Проверяем, что фигура в принципе может двигаться в этом направлении (f.e. диагонально для слона, и т.д.)
	var cellsToBePassed []int
	cellsToBePassed, err = checkFigureMove(figure, from, to)
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
	err = checkEndCell(&b, to)
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// TODO: проверяем, что пользователь не захотел выставить нового ферзя в неуместном для этого случае.

	// 6. На текущем этапе ход возможен. Генерируем новое положение доски newBoard.
	err = b.Move(board.Sq(from), board.Sq(to)) // TODO учесть промоушен пешки
	if err != nil {
		return fmt.Errorf("move invalid: %w", err)
	}

	// 7. Проверяем, что при новой позиции на доске не появился шах для собственного короля оппонента o.
	err = checkSelfCheck(b)
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
func checkFigureColor(b board.Board) bool {
	// TODO написать логику
	return true
}

// checkFigureMove проверяет, что фигура f может двигаться в этом направлении с клетки startCell на клетку
// endCell. Возвращает ошибку, если движение невозможно, или nil в противном случае. Также возвращают массив
// cellsToBePassed []Cell, в который входят все "промежуточные" клетки, которые "проходит" эта фигура на своем
// пути с startCell до endCell, если такие есть. В противном случае возвращает nil.
func checkFigureMove(p board.Piece, from, to int) (cellsToBePassed []int, err error) {
	// TODO написать логику

	// через цепочку if идет проверка на тип Figure (или на поле name в структуре Figure) и перенаправление на
	// соответствующий метод Move для этого типа фигуры
	return cellsToBePassed, nil
}

// checkCellsToBePassed проверяет, есть ли на клетках из массива cellsToBePassed какие-либо фигуры. Если хотя бы на одной
// клетке есть фигура, возвращается ошибка. Иначе возвращается nil.
func checkCellsToBePassed(b board.Board, cellsToBePassed []int) error {
	// TODO написать логику

	for _, cell := range cellsToBePassed {
		_, err := b.Get(board.Sq(cell))
		if err != nil {
			return fmt.Errorf("figure present on cell %d", cell)
		}
	}
	return nil
}

// checkEndCell проверяет наличие фигуры на клетке endCell на предмет совместимости хода. Возвращает ошибку при
// несовместимости хода или nil в случае успеха.
func checkEndCell(b *board.Board, to int) error {
	// TODO написать логику

	// логика пошагово:
	// 1. Если в клетке endCell нет фигур, ход возможен:
	_, err := b.Get(board.Sq(to)) // TODO: перевод типов в клетку (square неэкспортируемый)
	if err == nil {
		return nil
	}

	// 2. Если фигура в endCell принадлежит самому участнику, ход невозможен.
	//

	// 3. Если фигура в endCell принадлежит сопернику, проверка, возможно ли взятие
	// (f.e. пешка e2-e4 не может взять коня на e4).
	//
	return nil
}

// checkSelfCheck проверяет, что при новой позиции на доске нет шаха собственному королю.
func checkSelfCheck(b board.Board) error {
	// TODO написать логику

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
func getKingCell(b board.Board) int {
	// TODO написать логику
	return 0
}

// checkEnemyKnightsNearKing проверяет ближайшие клетки в расположении буквой Г к своему королю на наличие на них
// вражеского коня. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemyKnightsNearKing(b board.Board, kingCell int) error {
	// TODO написать логику
	return nil
}

// checkEnemiesVerticallyAndHorizontally проверяет ближайшие клетки по вертикали (сверху, снизу) и
// горизонтали (слева, справа) по отношению к клетке своего короля kingCell, на которых находятся вражеские ладьи и ферзи.
// Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesVerticallyAndHorizontally(b board.Board, kingCell int) error {
	// TODO написать логику
	return nil
}

// checkEnemiesDiagonally проверяет ближайшие клетки по всем диагоналям по отношению к клетке своего короля kingCell, на
// которых находятся вражеские слоны, ферзи и пешки. Если таких клеток нет, возвращется nil, иначе сообщение об ошибке.
func checkEnemiesDiagonally(b board.Board, kingCell int) error {
	// должна быть реализована дополнтельная проверка на пешки - их битое поле только по ходу движения!
	// TODO написать логику
	return nil
}
