// Пакет game реализует микросервис игры в шахматы.
package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sadmadrus/chessBox/internal/board"
)

// active содержит игры под управлением данного микросервиса.
var active map[id]*game

// game представляет игру в шахматы.
type game struct {
	id      id
	board   board.Board
	manager string // менеджер игровых сессий
	white   string // гейт, играющий за белых
	black   string // гейт, играющий за чёрных
	state   state
	moves   []move // основная вариация
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

// move содержит информацию о ходе.
type move struct {
	white, black validMove
}

// validMove содержит информацию о ходе одного игрока.
type validMove interface {
	aMove()
}

type halfmove struct {
	from, to board.Square
}

func (h halfmove) aMove() {
}

// promotion описывает ход с проведением пешки.
type promotion struct {
	halfmove
	promoteTo board.Piece
}

func (p promotion) aMove() {
}

type castling board.Castling

func (c castling) aMove() {
}

// gameState — модель для ответа о состоянии игры.
// TODO: остальные поля
type gameState struct {
	FEN string `json:"fen"`
}

func init() {
	active = make(map[id]*game)
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

// registerAndServe запускает игру.
func (g *game) registerAndServe() error {
	if _, ok := active[g.id]; ok {
		return fmt.Errorf("game already registered")
	}
	http.HandleFunc("/"+g.id.string(), g.handler())
	active[g.id] = g
	log.Printf("Started serving game: %s", g.id.string())
	return nil
}

// handler — "ручка" для конкретной игры.
func (g *game) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			g.serveCurrentState(w)
		case http.MethodPut, http.MethodDelete:
			http.Error(w, "Not (yet) implemented.", http.StatusNotImplemented)
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
