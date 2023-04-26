// Пакет game реализует микросервис игры в шахматы.
package game

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/validation"
)

const gameRequestTimeout = time.Second * 3

var (
	errGameNotFound       = errors.New("game not found")
	errGameOver           = errors.New("game is already over")
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

// start создаёт новую игру.  g.id будет проигнорирован.
func start(manager, white, black string, g *game) (id, error) {
	// TODO тут будет проверка, отвечают ли
	reqCh := make(chan request)

	var gameId id
	loaded := true
	for loaded {
		gameId = newId()
		_, loaded = games.LoadOrStore(gameId, reqCh)
	}

	if g == nil {
		g = &game{board: *board.Classical()}
	}
	g.id = gameId

	go func(in <-chan request) {
		for {
			req := <-in
			switch req.kind {
			case showState:
				g.returnState(req.replyTo)
			case stopGame:
				return
			default:
				g.processRequest(req)
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

// processRequest обрабатывает запрос от игрока.
func (g *game) processRequest(r request) {
	var err error
	switch r.kind {
	case makeMove:
		err = g.processMoveRequest(r)
	case forfeit:
		err = g.forfeit(r)
	}
	res := response{
		err:   err,
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
	if g.status != ongoing {
		return errGameOver
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

// forfeit обрабатывает сдачу в игре.
func (g *game) forfeit(r request) error {
	if g.status != ongoing {
		return errGameOver
	}
	if r.player != white && r.player != black {
		return errNoPlayerSpecified
	}
	if r.player == white {
		g.win(black)
	}
	g.win(white)
	return nil
}

// win переводит игру в состояние выигрыша указанного игрока.
func (g *game) win(p player) {
	switch p {
	case white:
		g.status = whiteWon
	case black:
		g.status = blackWon
	}
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

// newId генерирует id для игры (уникальность здесь не проверяется).
func newId() id {
	return id(fmt.Sprint(time.Now().UTC().Format("002150405"), rand.Int()))
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

	select {
	case ch.(chan request) <- r:
	case <-time.After(gameRequestTimeout):
		log.Printf("sending request to game %v timed out", game.string())
		return response{}, errGameRequestTimeout
	}

	select {
	case res := <-r.replyTo:
		return res, nil
	case <-time.After(gameRequestTimeout):
		log.Printf("waiting for response from game %v timed out", game.string())
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
