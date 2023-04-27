// Пакет game реализует игры в шахматы, идущие одновременно на множестве досок.
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
	ErrCantParse          = errors.New("failed to parse")
	ErrGameNotFound       = errors.New("game not found")
	ErrGameOver           = errors.New("game is already over")
	ErrGameRequestTimeout = errors.New("game request timed out")
	ErrInvalidMove        = errors.New("move is invalid")
	ErrNoPlayerSpecified  = errors.New("player is required, but no valid player is provided")
	ErrWrongTurn          = errors.New("wrong turn to move")
)

// games содержит игры под управлением данного микросервиса.
// Использует ключ типа id и значение типа chan request.
var games sync.Map

// game представляет игру в шахматы.
type game struct {
	id     ID
	board  board.Board
	status state
	moves  []fullMove // основная вариация
}

// ID — идентификатор игры.
type ID string

// state — состояние игры (продолжается, закончена)
type state int

const (
	ongoing state = iota
	drawn
	whiteWon
	blackWon
)

// Request представляет собой запрос на изменение состояния игры.
type Request struct {
	Player  Player
	Kind    RequestType
	Move    Move
	ReplyTo chan response // канал будет закрыт после отсылки ответа
}

// response представляет собой ответ игры на запрос об изменении состояния.
type response struct {
	err   error
	state State
}

type Player int

const (
	White Player = iota + 1
	Black
)

// RequestType показывает, что нужно сделать с игрой.
type RequestType int

const (
	MakeMove RequestType = iota + 1
	TakebackMove
	OfferDraw
	OfferAdjourn
	Forfeit
	ShowState
	Delete
)

// State — модель для ответа о состоянии игры.
// TODO: остальные поля
type State struct {
	Status string `json:"status"`
	FEN    string `json:"fen"`
}

// New создаёт новую игру из начальной позиции.
func New(manager, white, black string) (ID, error) {
	return start(manager, white, black, nil)
}

// start создаёт новую игру.  g.id будет проигнорирован.
func start(manager, white, black string, g *game) (ID, error) {
	// TODO тут будет проверка, отвечают ли
	reqCh := make(chan Request)

	var gameId ID
	loaded := true
	for loaded {
		gameId = newId()
		_, loaded = games.LoadOrStore(gameId, reqCh)
	}

	if g == nil {
		g = &game{board: *board.Classical()}
	}
	g.id = gameId

	go func(in <-chan Request) {
		for {
			req := <-in
			switch req.Kind {
			case ShowState:
				g.returnState(req.ReplyTo)
			case Delete:
				games.Delete(gameId)
				return
			default:
				g.processRequest(req)
			}
		}
	}(reqCh)

	log.Printf("Started serving game: %s", g.id.String())

	return gameId, nil
}

// retrace восстанавливает состояние игры из последовательности ходов.
func retrace(mm []fullMove) (*game, error) {
	g := &game{board: *board.Classical()}
	for i, fm := range mm {
		err := g.move(fm.white)
		if err != nil {
			return nil, fmt.Errorf("on move %v white error: %w", i+1, err)
		}
		if i == len(mm)-1 && fm.black == nil {
			continue
		}
		err = g.move(fm.black)
		if err != nil {
			return nil, fmt.Errorf("on move %v black error: %w", i+1, err)
		}
	}
	return g, nil
}

// returnState возвращает состояние игры.
func (g *game) returnState(out chan<- response) {
	out <- response{nil, g.state()}
	close(out)
}

func (g *game) state() State {
	return State{
		Status: g.status.string(),
		FEN:    g.board.FEN(),
	}
}

// processRequest обрабатывает запрос от игрока.
func (g *game) processRequest(r Request) {
	var err error
	switch r.Kind {
	case MakeMove:
		err = g.processMoveRequest(r)
	case Forfeit:
		err = g.forfeit(r)
	}
	res := response{
		err:   err,
		state: g.state(),
	}
	r.ReplyTo <- res
	close(r.ReplyTo)
}

// processMoveRequest обрабатывает запрос хода и совершает ход, если он легален;
// если возвращена ошибка, состояние игры не изменилось.
func (g *game) processMoveRequest(r Request) error {
	if r.Player == 0 {
		return ErrNoPlayerSpecified
	}
	if g.status != ongoing {
		return ErrGameOver
	}
	if !moveIsInTurn(r.Player, g.board.NextToMove()) {
		return ErrWrongTurn
	}
	err := g.move(r.Move)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMove, err)
	}
	return nil
}

// move совершает ход. Если возвращена ошибка, состояние игры не изменилось.
func (g *game) move(m Move) error {
	var promoteTo board.Piece
	if p, ok := m.(promotion); ok {
		promoteTo = p.toPiece()
	}

	err := validation.CanMove(g.board, m.FromSquare(), m.ToSquare(), promoteTo)
	if err != nil {
		return err
	}

	switch v := m.(type) {
	case simpleMove:
		err = g.board.Move(v.FromSquare(), v.ToSquare())
	case promotion:
		err = g.board.Promote(v.FromSquare(), v.ToSquare(), v.toPiece())
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
func (g *game) forfeit(r Request) error {
	if g.status != ongoing {
		return ErrGameOver
	}
	if r.Player != White && r.Player != Black {
		return ErrNoPlayerSpecified
	}
	if r.Player == White {
		g.win(Black)
	}
	g.win(White)
	return nil
}

// win переводит игру в состояние выигрыша указанного игрока.
func (g *game) win(p Player) {
	switch p {
	case White:
		g.status = whiteWon
	case Black:
		g.status = blackWon
	}
}

// moveIsInTurn возвращает true, если ход этого игрока.
func moveIsInTurn(p Player, whiteToMove bool) bool {
	if p == White {
		return whiteToMove
	}
	if p == Black {
		return !whiteToMove
	}
	return false
}

// Exists сообщает, существует ли игра.
func (id ID) Exists() bool {
	_, ok := games.Load(id)
	return ok
}

// newId генерирует id для игры (уникальность здесь не проверяется).
func newId() ID {
	return ID(fmt.Sprint(time.Now().UTC().Format("002150405"), rand.Int()))
}

func (id ID) String() string {
	return string(id)
}

func (id ID) Do(r Request) (State, error) {
	ch, ok := games.Load(id)
	if !ok {
		return State{}, ErrGameNotFound
	}

	if r.ReplyTo == nil {
		r.ReplyTo = make(chan response)
	}

	select {
	case ch.(chan Request) <- r:
	case <-time.After(gameRequestTimeout):
		log.Printf("sending request to game %v timed out", id.String())
		return State{}, ErrGameRequestTimeout
	}

	select {
	case res := <-r.ReplyTo:
		err := res.err
		return res.state, err
	case <-time.After(gameRequestTimeout):
		log.Printf("waiting for response from game %v timed out", id.String())
		return State{}, ErrGameRequestTimeout
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
