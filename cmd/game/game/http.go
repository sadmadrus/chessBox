package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Creator — http.HandlerFunc для создания новой игры.
func Creator(w http.ResponseWriter, r *http.Request) {
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

	g, err := start(manager, white, black, nil)
	if err != nil {
		// TODO проверка, не вернуть ли 408
		http.Error(w, fmt.Sprintf("Couldn't create the game: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Add("location", "/"+string(g))
	w.WriteHeader(http.StatusCreated)
	serveCurrentState(g, w)
}

// GameHandler обрабатывает запросы к играм.
func GameHandler(w http.ResponseWriter, r *http.Request) {
	g := id(strings.TrimPrefix(r.URL.Path, "/"))
	_, ok := games.Load(g)
	if !ok {
		http.Error(w, "404 Game Not Found", http.StatusNotFound)
		return
	}
	handler(g)(w, r)
}

// handler — "ручка" для конкретной игры.
func handler(game id) http.HandlerFunc {
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
func serveCurrentState(game id, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	r := request{kind: showState}
	res, err := requestWithTimeout(r, game)
	if err != nil {
		if errors.Is(err, errGameNotFound) {
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
			return
		}
		if errors.Is(err, errGameRequestTimeout) {
			http.Error(w, "503 Status Unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		log.Printf("unexpected error: %v", err)
		return
	}
	serveState(w, res.state)
}

// serveState пишет в ответ текущее состояние игры.
func serveState(w http.ResponseWriter, s gameState) {
	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, "Failed to encode the current game state.", http.StatusInternalServerError)
		log.Printf("error encoding game state: %v", err)
	}
}

// deleteGame удаляет игру из сервиса.
func deleteGame(game id) {
	log.Printf("Removing game: %s", game.string())
	games.Delete(game)
	r := request{kind: stopGame}
	_, _ = requestWithTimeout(r, game)
}

// handlePut обрабатывает PUT-запросы для игры.
func handlePut(game id, w http.ResponseWriter, r *http.Request) {
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := request{player: getPlayer(data)}

	req.kind = identifyRequest(data)
	if req.kind == 0 {
		http.Error(w, "can't parse the request", http.StatusBadRequest)
		return
	}

	switch req.kind {
	case makeMove:
		req.move, err = parseUCI(data.Get("move"))
		if err != nil {
			http.Error(w, "can't parse move", http.StatusBadRequest)
			return
		}
	case forfeit:
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}

	res, err := requestWithTimeout(req, game)
	if err != nil {
		if errors.Is(err, errGameNotFound) {
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
			return
		}
		if errors.Is(err, errGameRequestTimeout) {
			http.Error(w, "503 Status Unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, "Unknown Server Error", http.StatusInternalServerError)
		log.Printf("unexpected error: %v", err)
		return
	}
	if res.err != nil {
		switch {
		case errors.Is(res.err, errGameNotFound):
			http.Error(w, "404 Game Not Found", http.StatusNotFound)
			return
		case errors.Is(res.err, errWrongTurn) || errors.Is(err, errGameOver):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(res.err, errInvalidMove):
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	serveState(w, res.state)
}

func getPlayer(data url.Values) player {
	switch data.Get("player") {
	case "white":
		return white
	case "black":
		return black
	default:
		return 0
	}
}

// identifyRequest возвращает тип запроса.
func identifyRequest(data url.Values) kindOfRequest {
	var have []kindOfRequest

	keys := map[string]kindOfRequest{
		"move":      makeMove,
		"takeback":  takeback,
		"drawoffer": draw,
		"adjourn":   adjourn,
		"forfeit":   forfeit,
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
