// Пакет game реализует микросервис игры в шахматы.
package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/google/uuid"
	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/validation"
)

const errCantParse = "failed to parse"

// games содержит игры под управлением данного микросервиса.
var games sync.Map

// game представляет игру в шахматы.
type game struct {
	id      id
	board   board.Board
	manager string // менеджер игровых сессий
	white   string // гейт, играющий за белых
	black   string // гейт, играющий за чёрных
	state   state
	moves   []fullMove // основная вариация
}

// id — идентификатор игры.
type id string

// state — состояние игры (продолжается, закончена)
type state int

const (
	ongoing state = iota
	drawn
	whiteWon
	blackWon
)

// kindOfRequest показывает, что хочет игрок.
type kindOfRequest int

const (
	makeMove kindOfRequest = iota + 1
	takeback
	draw
	adjourn
	forfeit
)

// gameState — модель для ответа о состоянии игры.
// TODO: остальные поля
type gameState struct {
	FEN string `json:"fen"`
}

// new создаёт новую игру.
func new(manager, white, black string) (*game, error) {
	// TODO тут будет проверка, отвечают ли
	return &game{
		id:      newId(),
		manager: manager,
		white:   white,
		black:   black,
		state:   ongoing,
		board:   *board.Classical(),
	}, nil
}

// start запускает игру.
func (g *game) start() error {
	if _, ok := games.Load(g.id); ok {
		return fmt.Errorf("game already registered")
	}
	games.Store(g.id, g)
	log.Printf("Started serving game: %s", g.id.string())
	return nil
}

// handler — "ручка" для конкретной игры.
func (g *game) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			g.serveCurrentState(w)
		case http.MethodPut:
			g.handlePut(w, r)
		case http.MethodDelete:
			g.stop()
		case http.MethodOptions:
			w.Header().Set("Allow", "GET, PUT, DELETE, OPTIONS")
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE, OPTIONS")
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		}
	}
}

// serveCurrentState пишет в ответ текущее состояние игры.
func (g *game) serveCurrentState(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(g.currentState()); err != nil {
		http.Error(w, "Failed to encode the current game state.", http.StatusInternalServerError)
		log.Printf("error encoding game state: %v", err)
		return
	}
}

// currentState возвращает состояние игры.
func (g *game) currentState() gameState {
	return gameState{
		FEN: g.board.FEN(),
	}
}

// handlePut обрабатывает PUT-запросы для игры.
func (g *game) handlePut(w http.ResponseWriter, r *http.Request) {
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player := data.Get("player")
	if player != "white" && player != "black" {
		http.Error(w, "no valid player specified", http.StatusBadRequest)
		return
	}

	req := identifyRequest(data)
	if req == 0 {
		http.Error(w, "can't parse the request", http.StatusBadRequest)
		return
	}

	switch req {
	case makeMove:
		g.processMoveRequest(data, w)
	default:
		http.Error(w, "not (yet) implemented", http.StatusNotImplemented)
	}
	g.serveCurrentState(w)
}

// stop удаляет игру из сервиса.
func (g *game) stop() {
	log.Printf("Removing game: %s", g.id.string())
	games.Delete(g.id)
}

// processMoveRequest обрабатывает запрос на совершение хода
func (g *game) processMoveRequest(data url.Values, w http.ResponseWriter) {
	if !moveIsInTurn(data.Get("player"), g.board.NextToMove()) {
		http.Error(w, "wrong turn", http.StatusConflict)
	}

	m, err := parseUCI(data.Get("move"))
	if err != nil {
		http.Error(w, "can't parse move", http.StatusBadRequest)
		return
	}

	if err := g.move(m); err != nil {
		http.Error(w, fmt.Sprintf("move not allowed: %v", err), http.StatusForbidden)
	}
}

// move совершает ход. Если возвращена ошибка, состояние игры не изменилось.
func (g *game) move(m halfMove) error {
	var promoteTo board.Piece
	if p, ok := m.(promotion); ok {
		promoteTo = p.toPiece()
	}

	err := validation.CanMove(g.board, m.fromSquare(), m.toSquare(), promoteTo)
	if err != nil {
		return err
	}

	switch v := m.(type) {
	case simpleMove:
		err = g.board.Move(v.fromSquare(), v.toSquare())
	case promotion:
		err = g.board.Promote(v.fromSquare(), v.toSquare(), v.toPiece())
	case castling:
		err = g.board.Castle(board.Castling(v))
	default:
		err = fmt.Errorf("unknown move type")
	}

	if err != nil {
		return err
	}

	if g.board.NextToMove() {
		g.moves[len(g.moves)-1].black = m
	} else {
		g.moves = append(g.moves, fullMove{white: m})
	}

	return nil
}

// moveIsInTurn возвращает true, если ход этого игрока.
func moveIsInTurn(player string, whiteToMove bool) bool {
	if player == "white" {
		return whiteToMove
	}
	return !whiteToMove
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

// newId генерирует уникальный id для игры.
func newId() id {
	err := fmt.Errorf("not nil")
	var u uuid.UUID
	for err != nil {
		u, err = uuid.NewRandom()
	}
	return id(u.String())
}

func (i id) string() string {
	return string(i)
}
