// Пакет game реализует микросервис игры в шахматы.
package game

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/validation"
)

const gameRequestTimeout = time.Second * 3

var (
	errGameNotFound       = errors.New("game not found")
	errGameRequestTimeout = errors.New("game request timed out")
	errInvalidMove        = errors.New("move is invalid")
	errNoPlayerSpecified  = errors.New("player is required, but no valid player is provided")
	errWrongTurn          = errors.New("wrong turn to move")
)

const errCantParse = "failed to parse"

// games содержит игры под управлением данного микросервиса.
// Использует ключ типа id и значение типа chan request.
var games sync.Map

// game представляет игру в шахматы.
type game struct {
	id     id
	board  board.Board
	status state
	moves  []fullMove // основная вариация
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

// request представляет собой запрос на изменение состояния игры.
type request struct {
	player  player
	kind    kindOfRequest
	move    halfMove
	replyTo chan response // канал будет закрыт после отсылки ответа
}

// response представляет собой ответ игры на запрос об изменении состояния.
type response struct {
	err   error
	state gameState
}

// player указывает на игрока, отправившего запрос.
type player int

const (
	white player = iota + 1
	black
)

// kindOfRequest показывает, что нужно сделать с игрой.
type kindOfRequest int

const (
	makeMove kindOfRequest = iota + 1
	takeback
	draw
	adjourn
	forfeit
	showState
	stopGame
)

// gameState — модель для ответа о состоянии игры.
// TODO: остальные поля
type gameState struct {
	Status string `json:"status"`
	FEN    string `json:"fen"`
}

// start создаёт новую игру.
func start(manager, white, black string) (id, error) {
	// TODO тут будет проверка, отвечают ли
	reqCh := make(chan request)

	var gameId id
	loaded := true
	for loaded {
		gameId = newId()
		_, loaded = games.LoadOrStore(gameId, reqCh)
	}

	g := &game{id: gameId, board: *board.Classical()}

	go func(in <-chan request) {
		for {
			req := <-in
			switch req.kind {
			case showState:
				g.returnState(req.replyTo)
			case stopGame:
				return
			case makeMove:
				g.processMove(req)
			}
		}
	}(reqCh)

	log.Printf("Started serving game: %s", g.id.string())

	return gameId, nil
}

// returnState возвращает состояние игры.
func (g *game) returnState(out chan<- response) {
	out <- response{nil, g.state()}
	close(out)
}

func (g *game) state() gameState {
	return gameState{
		Status: g.status.string(),
		FEN:    g.board.FEN(),
	}
}

// processMove обрабатывает запрос хода.
func (g *game) processMove(r request) {
	res := response{
		err:   g.processMoveRequest(r),
		state: g.state(),
	}
	r.replyTo <- res
	close(r.replyTo)
}

// processMoveRequest обрабатывает запрос хода и совершает ход, если он легален;
// если возвращена ошибка, состояние игры не изменилось.
func (g *game) processMoveRequest(r request) error {
	if r.player == 0 {
		return errNoPlayerSpecified
	}
	if !moveIsInTurn(r.player, g.board.NextToMove()) {
		return errWrongTurn
	}
	err := g.move(r.move)
	if err != nil {
		return fmt.Errorf("%w: %v", errInvalidMove, err)
	}
	return nil
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
func moveIsInTurn(p player, whiteToMove bool) bool {
	if p == white {
		return whiteToMove
	}
	if p == black {
		return !whiteToMove
	}
	return false
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

func requestWithTimeout(r request, game id) (response, error) {
	ch, ok := games.Load(game)
	if !ok {
		return response{}, errGameNotFound
	}

	if r.replyTo == nil {
		r.replyTo = make(chan response)
	}

	go func() { ch.(chan request) <- r }()

	select {
	case res := <-r.replyTo:
		return res, nil
	case <-time.After(gameRequestTimeout):
		return response{}, errGameRequestTimeout
	}
}

func (st state) string() string {
	switch st {
	case ongoing:
		return "ongoing"
	case drawn:
		return "1/2-1/2"
	case whiteWon:
		return "1-0"
	case blackWon:
		return "0-1"
	default:
		return ""
	}
}
