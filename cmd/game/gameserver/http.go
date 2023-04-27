// Пакет gameserver реализует HTTP API для сервиса игровой сессии.
package gameserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sadmadrus/chessBox/internal/game"
)

// RootHandler отвечает за обработку запросов к сервису в целом.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		gameHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "The game server is online and working.")
	case http.MethodPost:
		creator(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

// creator — http.HandlerFunc для создания новой игры.
func creator(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manager := data.Get("notify")
	white := data.Get("white")
	black := data.Get("black")
	if manager == "" || white == "" || black == "" {
		http.Error(w, "Required parameter missing.", http.StatusBadRequest)
		return
	}

	for _, p := range []string{"position", "timing", "move1", "timewhite", "timeblack"} {
		if _, ok := data[p]; ok {
			http.Error(w, "Not (yet) implemented.", http.StatusNotImplemented)
			return
		}
	}

	g, err := game.New(manager, white, black)
	if err != nil {
		// TODO проверка, не вернуть ли 408
		http.Error(w, fmt.Sprintf("Couldn't create the game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("location", "/"+string(g))
	w.WriteHeader(http.StatusCreated)
	serveCurrentState(g, w)
}

// gameHandler обрабатывает запросы к играм.
func gameHandler(w http.ResponseWriter, r *http.Request) {
	g := game.ID(strings.TrimPrefix(r.URL.Path, "/"))
	if !g.Exists() {
		http.Error(w, "404 Game Not Found", http.StatusNotFound)
		return
	}
	handler(g)(w, r)
}

// handler — "ручка" для конкретной игры.
func handler(game game.ID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			serveCurrentState(game, w)
		case http.MethodPut:
			handlePut(game, w, r)
		case http.MethodDelete:
			deleteGame(game)
		case http.MethodOptions:
			w.Header().Set("Allow", "GET, PUT, DELETE, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE, OPTIONS")
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

// serveCurrentState запрашиевает и пишет в ответ текущее состояние игры.
func serveCurrentState(g game.ID, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	r := game.Request{Kind: game.ShowState}
	res, err := g.Do(r)
	if err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
			return
		}
		if errors.Is(err, game.ErrGameRequestTimeout) {
			http.Error(w, "503 Status Unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		log.Printf("unexpected error: %v", err)
		return
	}
	serveState(w, res)
}

// serveState пишет в ответ текущее состояние игры.
func serveState(w http.ResponseWriter, s game.State) {
	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, "Failed to encode the current game state.", http.StatusInternalServerError)
		log.Printf("error encoding game state: %v", err)
	}
}

// deleteGame удаляет игру из сервиса.
func deleteGame(g game.ID) {
	log.Printf("Removing game: %s", g.String())
	r := game.Request{Kind: game.Delete}
	_, _ = g.Do(r)
}

// handlePut обрабатывает PUT-запросы для игры.
func handlePut(g game.ID, w http.ResponseWriter, r *http.Request) {
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := game.Request{Player: getPlayer(data)}

	req.Kind = identifyRequest(data)
	if req.Kind == 0 {
		http.Error(w, "can't parse the request", http.StatusBadRequest)
		return
	}

	switch req.Kind {
	case game.MakeMove:
		req.Move, err = game.ParseUCI(data.Get("move"))
		if err != nil {
			http.Error(w, "can't parse move", http.StatusBadRequest)
			return
		}
	case game.Forfeit:
	default:
		w.WriteHeader(http.StatusNotImplemented)
		serveCurrentState(g, w)
		return
	}

	res, err := g.Do(req)
	if err != nil {
		switch {
		case errors.Is(err, game.ErrGameNotFound):
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
			return
		case errors.Is(err, game.ErrGameRequestTimeout):
			http.Error(w, "503 Status Unavailable", http.StatusServiceUnavailable)
			return
		case errors.Is(err, game.ErrWrongTurn) || errors.Is(err, game.ErrGameOver):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(err, game.ErrInvalidMove):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	serveState(w, res)
}

func getPlayer(data url.Values) game.Player {
	switch data.Get("player") {
	case "white":
		return game.White
	case "black":
		return game.Black
	default:
		return 0
	}
}

// identifyRequest возвращает тип запроса.
func identifyRequest(data url.Values) game.RequestType {
	var have []game.RequestType

	keys := map[string]game.RequestType{
		"move":      game.MakeMove,
		"takeback":  game.TakebackMove,
		"drawoffer": game.OfferDraw,
		"adjourn":   game.OfferAdjourn,
		"forfeit":   game.Forfeit,
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
