// Пакет validation валидирует ходы по текущему состоянию доски, начальной и конечной клеткам, и фигуре.

package validation

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sadmadrus/chessBox/internal/board"
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
		piece, err := parsePiece(pieceParsed)
		if err != nil {
			log.Printf("%v: %v", errPieceNotExist, pieceParsed)
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
				log.Printf("%v: %v", errPieceNotExist, fromParsed)
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
				log.Printf("%v: %s", errPieceNotExist, toParsed)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		// валидация входных данных: клетки from и to различны
		if from == to {
			log.Printf("%v: %v (from), %v (to)", errFromToSquaresNotDiffer, from, to)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация геометрического движения фигуры без привязки к позиции на доске
		fromSquare := newSquare(int8(from))
		toSquare := newSquare(int8(to))
		err = move(piece, fromSquare, toSquare)
		if err != nil {
			log.Printf("%v: from %v - to %v", err, from, to)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	log.Printf("inside Simple %v: %v", errInvalidHttpMethod, r.Method)
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
		// TODO: board нужно ли проверять валидность доски где-то вообще?
		boardParsed := r.URL.Query().Get("board")
		b, err := board.FromUsFEN(boardParsed)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: клетка from существуют
		fromParsed := r.URL.Query().Get("from")
		var fromSquare square
		fromSquare, err = parseSquare(fromParsed)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: клетка to существуют
		toParsed := r.URL.Query().Get("to")
		var toSquare square
		toSquare, err = parseSquare(toParsed)
		if err != nil {
			log.Printf("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: клетки from и to различны
		if !fromSquare.isEqual(toSquare) {
			log.Printf("%v: %v (from), %v (to)", errFromToSquaresNotDiffer, fromSquare, toSquare)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// валидация входных данных: фигура newpiece принимает q/r/b/n/Q/R/B/N или пустое значение
		newpieceParsed := r.URL.Query().Get("newpiece")
		var newpiece board.Piece
		newpiece, err = parseNewpiece(newpieceParsed)
		if err != nil {
			log.Printf("%v: %v", errNewpieceNotValid, newpieceParsed)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// полная валидация хода с учетом положения на доске, а также возможных рокировок, взятия на проходе и
		// проведения пешки
		newBoard, isValid, err := advancedLogic(*b, fromSquare, toSquare, newpiece)
		if err != nil {
			log.Printf("error occured when trying to validate move: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !isValid {
			log.Printf("move invalid: %v", err)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			boardUsFEN := newBoard.UsFEN()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			data := advancedResponse{boardUsFEN}
			err = json.NewEncoder(w).Encode(data)
			if err != nil {
				log.Printf("error while encoding json: %v", err)
			}
			return
		}
	}

	log.Printf("inside Advanced %v: %v", errInvalidHttpMethod, r.Method)
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
