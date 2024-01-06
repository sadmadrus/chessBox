package gameserver

import (
	"errors"
	"sync"

	"github.com/sadmadrus/chessBox/internal/board/validation"
	"github.com/sadmadrus/chessBox/internal/chess"
)

// Storage предоставляет хранилище для ходов игры.
type Storage interface {
	LoadGame(ID) (Game, error)           // Загрузить состояние игры из хранилища.
	TakebackMove(ID, int) error          // Отменить ход(ы) так, чтобы номер последнего полухода (0-based) совпал с заданным.
	StoreMove(ID, int, chess.Move) error // Сохранить очередной полуход с заданным номером (0-based).
	StoreGame(Game) error                // Сохранить новую игру в заданном состоянии; должно вернуть ErrGameExists, если игра с этим ID существует.
	StoreResult(ID, chess.State) error   // Сохранить состояние игры (если она закончена).
}

var (
	ErrGameExists     = errors.New("game with this ID exists")
	ErrGameNotFound   = errors.New("game not found")
	ErrMoveNoMismatch = errors.New("this move number is not the next move for this game")
	ErrNoSuchMove     = errors.New("no such move in this game")
)

// MemoryStorage — хранилище игр в памяти для локального инстанса. Для тестов и
// как proof-of-concept.
type MemoryStorage struct {
	m     sync.Mutex
	games map[ID]Game
}

func (ms *MemoryStorage) StoreGame(g Game) error {
	if !validation.IsLegal(g.StartingPosition) {
		return ErrInvalidPosition
	}

	ms.m.Lock()
	defer ms.m.Unlock()

	if _, exists := ms.games[g.ID]; exists {
		return ErrGameExists
	}

	ms.games[g.ID] = g
	return nil
}

func (ms *MemoryStorage) StoreMove(id ID, moveNo int, move chess.Move) error {
	ms.m.Lock()
	defer ms.m.Unlock()

	g, ok := ms.games[id]
	if !ok {
		return ErrGameNotFound
	}

	if len(g.Moves) != moveNo {
		return ErrMoveNoMismatch
	}

	g.Moves = append(g.Moves, move)
	ms.games[id] = g
	return nil
}

func (ms *MemoryStorage) TakebackMove(id ID, moveNo int) error {
	ms.m.Lock()
	defer ms.m.Unlock()

	g, ok := ms.games[id]
	if !ok {
		return ErrGameNotFound
	}

	if len(g.Moves) < moveNo+1 {
		return ErrNoSuchMove
	}

	g.Moves = g.Moves[:moveNo+1]
	ms.games[id] = g
	return nil
}

func (ms *MemoryStorage) LoadGame(id ID) (Game, error) {
	ms.m.Lock()
	defer ms.m.Unlock()

	var err error
	g, ok := ms.games[id]
	if !ok {
		err = ErrGameNotFound
	}

	return g, err
}

func (ms *MemoryStorage) StoreResult(id ID, s chess.State) error {
	ms.m.Lock()
	defer ms.m.Unlock()

	g, ok := ms.games[id]
	if !ok {
		return ErrGameNotFound
	}

	g.State = s
	ms.games[id] = g
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	games := make(map[ID]Game)
	ms := MemoryStorage{games: games}
	return &ms
}
