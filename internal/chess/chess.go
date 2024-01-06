// Пакет chess реализует игру в шахматы.
package chess

import (
	"errors"
	"fmt"

	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/board/moves"
)

var (
	ErrCantParse          = errors.New("failed to parse")
	ErrGameOver           = errors.New("game is already over")
	ErrGameRequestTimeout = errors.New("game request timed out")
	ErrInvalidMove        = errors.New("move is invalid")
	ErrNoPlayerSpecified  = errors.New("player is required, but no valid player is provided")
	ErrWrongTurn          = errors.New("wrong turn to move")
)

// Game представляет игру в шахматы.
type Game struct {
	StartingPosition board.Board
	Moves            []Move // основная вариация
	State            State

	current *board.Board
}

// CurrentPosition возвращает текущую позицию в игре.
// Если в игре есть невалидные ходы, вернётся пустая доска.
func (g *Game) CurrentPosition() board.Board {
	if g.current != nil {
		return *g.current
	}

	g.updateCurrent()
	return *g.current
}

func (g *Game) updateCurrent() {
	pos := g.StartingPosition
	g.current = &pos

	for _, m := range g.Moves {
		err := move(g.current, m)
		if err != nil {
			g.current = &board.Board{}
			return
		}
	}
}

// MakeMove производит ход в игре. Если возвращена ошибка, состояние не
// изменилось.
func (g *Game) MakeMove(m Move, p Player) error {
	if g.State != Ongoing {
		return ErrGameOver
	}
	pos := g.CurrentPosition()
	if !moveIsInTurn(p, pos.NextToMove()) {
		return ErrWrongTurn
	}
	err := move(&pos, m)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidMove, err)
	}
	g.Moves = append(g.Moves, m)
	g.updateCurrent()

	// TODO: тут надо обработать возникшую позицию: мат, пат и т.п.

	return nil
}

// State — состояние игры (продолжается, закончена)
type State int

const (
	Ongoing State = iota
	Drawn
	WhiteWon
	BlackWon
)

type Player int

const (
	White Player = iota + 1
	Black
)

func (p Player) String() string {
	if p == White {
		return "white"
	}
	return "black"
}

// move совершает ход на доске. Если возвращена ошибка, состояние на доске не
// изменилось.
func move(b *board.Board, m Move) error {
	if m == nil {
		return ErrInvalidMove
	}

	var promoteTo board.Piece
	if p, ok := m.(promotion); ok {
		promoteTo = p.toPiece()
	}

	ok, err := moves.IsValid(*b, m.FromSquare(), m.ToSquare(), promoteTo)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidMove
	}

	switch v := m.(type) {
	case simpleMove:
		err = b.Move(v.FromSquare(), v.ToSquare())
	case promotion:
		err = b.Promote(v.FromSquare(), v.ToSquare(), v.toPiece())
	case castling:
		err = b.Castle(board.Castling(v))
	default:
		err = fmt.Errorf("unknown move type")
	}

	return err
}

// Forfeit обрабатывает сдачу в игре.
func (g *Game) Forfeit(p Player) error {
	if g.State != Ongoing {
		return ErrGameOver
	}
	if p != White && p != Black {
		return ErrNoPlayerSpecified
	}
	if p == White {
		g.win(Black)
	} else {
		g.win(White)
	}
	return nil
}

// win переводит игру в состояние выигрыша указанного игрока.
func (g *Game) win(p Player) {
	switch p {
	case White:
		g.State = WhiteWon
	case Black:
		g.State = BlackWon
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

func (st State) String() string {
	switch st {
	case Drawn:
		return "1/2-1/2"
	case WhiteWon:
		return "1-0"
	case BlackWon:
		return "0-1"
	default:
		return ""
	}
}
