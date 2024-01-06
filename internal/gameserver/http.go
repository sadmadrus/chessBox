// Пакет gameserver реализует HTTP API для сервиса игровой сессии.
package gameserver

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sadmadrus/chessBox/internal/chess"
)

const urlencoded = "application/x-www-form-urlencoded"

// HandleRoot возвращает ручку для обработки запросов к сервису в целом.
func HandleRoot(st Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handleGame(st, w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			fmt.Fprint(w, "The game server is online and working.")
		case http.MethodPost:
			creator(st, w, r)
		case http.MethodOptions:
			w.Header().Set("Allow", "GET, POST, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "GET, POST, OPTIONS")
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		}
	}
}

// creator — ручка для создания новой игры.
func creator(st Storage, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, p := range []string{"notify", "position", "timing", "move1", "timewhite", "timeblack"} {
		if _, ok := data[p]; ok {
			http.Error(w, "Not (yet) implemented.", http.StatusNotImplemented)
			return
		}
	}

	g := newGame(st)

	w.Header().Add("location", "/"+string(g.ID))
	w.WriteHeader(http.StatusCreated)
	serveState(w, g, nil)
}

// handleGame обрабатывает запросы к играм.
func handleGame(st Storage, w http.ResponseWriter, r *http.Request) {
	id := ID(strings.TrimPrefix(r.URL.Path, "/"))
	g, err := st.LoadGame(id)
	if err != nil {
		if errors.Is(err, ErrGameNotFound) {
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
		} else {
			log.Printf("unexpected error when loading game: %v\n", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	handler(st, g)(w, r)
}

// handler — "ручка" для конкретной игры.
func handler(st Storage, game Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			serveState(w, game, nil)
		case http.MethodPut:
			handlePut(st, game, w, r)
		case http.MethodOptions:
			w.Header().Set("Allow", "GET, PUT, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "GET, PUT, OPTIONS")
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		}
	}
}

// parseUrlEncoded возвращает данные из www-url-encoded.
func parseUrlEncoded(r *http.Request) (url.Values, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the request's body: %w", err)
	}
	data, err := url.ParseQuery(string(b))
	if err != nil {
		return data, fmt.Errorf("could not parse data: %w", err)
	}
	return data, nil
}

// serveState пишет в ответ состояние игры.
func serveState(w http.ResponseWriter, g Game, data url.Values) {
	w.Header().Set("Content-Type", urlencoded)

	if data == nil {
		data = make(url.Values)
	}

	data.Set("game", g.ID.String())

	pos := g.CurrentPosition()
	data.Set("position", pos.FEN())

	if g.State != chess.Ongoing {
		data.Set("result", g.State.String())
	}

	// TODO: add other fields

	_, err := fmt.Fprint(w, data.Encode())
	if err != nil {
		log.Printf("unexpected error: %v\n", err)
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
	}
}

// handlePut обрабатывает PUT-запросы для игры.
func handlePut(st Storage, g Game, w http.ResponseWriter, r *http.Request) {
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player := getPlayer(data)

	kind := identifyRequest(data)
	if kind == 0 {
		http.Error(w, "can't parse the request", http.StatusBadRequest)
		return
	}

	state := make(url.Values)

	switch kind {
	case MakeMove:
		var move chess.Move
		move, err = chess.ParseUCI(data.Get("move"))
		if err != nil {
			http.Error(w, "can't parse move", http.StatusBadRequest)
			return
		}
		if player == 0 {
			http.Error(w, "No Player Has Been Specified", http.StatusBadRequest)
			return
		}
		err = g.MakeMove(move, player)
		if err != nil {
			switch {
			case errors.Is(err, chess.ErrWrongTurn) || errors.Is(err, chess.ErrGameOver):
				w.WriteHeader(http.StatusConflict)
			case errors.Is(err, chess.ErrInvalidMove):
				w.WriteHeader(http.StatusForbidden)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			err = st.StoreMove(g.ID, len(g.Moves)-1, move)
			if err != nil {
				switch {
				case errors.Is(err, ErrMoveNoMismatch):
					http.Error(w, "Confilcting Data about Move", http.StatusConflict)
				case errors.Is(err, ErrGameNotFound):
					http.Error(w, "404 Game Not Found", http.StatusNotFound)
				default:
					http.Error(w, "500 Unknown Server Error", http.StatusInternalServerError)
				}
				return
			}
		}
		state.Set("movemade", player.String())
		state.Set("move", data.Get("move"))
	case Forfeit:
		err := g.Forfeit(player)
		if err != nil {
			switch {
			case errors.Is(err, chess.ErrGameOver):
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			err = st.StoreResult(g.ID, g.State)
			if err != nil {
				switch {
				case errors.Is(err, ErrGameNotFound):
					http.Error(w, "404 Game Not Found", http.StatusNotFound)
				default:
					http.Error(w, "500 Unknown Server Error", http.StatusInternalServerError)
				}
			}
		}
		state.Set("result", g.State.String())
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}

	serveState(w, g, state)
}

func getPlayer(data url.Values) chess.Player {
	switch data.Get("player") {
	case "white":
		return chess.White
	case "black":
		return chess.Black
	default:
		return 0
	}
}

// RequestType показывает, что нужно сделать с игрой.
type RequestType int

const (
	MakeMove RequestType = iota + 1
	TakebackMove
	OfferDraw
	OfferAdjourn
	Forfeit
	ShowState
)

// identifyRequest возвращает тип запроса.
func identifyRequest(data url.Values) RequestType {
	var have []RequestType

	keys := map[string]RequestType{
		"move":      MakeMove,
		"takeback":  TakebackMove,
		"drawoffer": OfferDraw,
		"adjourn":   OfferAdjourn,
		"forfeit":   Forfeit,
	}

	for k, v := range keys {
		if data.Get(k) != "" {
			have = append(have, v)
		}
	}

	if len(have) != 1 {
		return 0
	}

	return have[0]
}
