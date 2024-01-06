package gameserver

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sadmadrus/chessBox/internal/board"
	"github.com/sadmadrus/chessBox/internal/chess"
)

var (
	ErrInvalidPosition = errors.New("the supplied position is invalid")
)

// game представляет игровую сессию.
type Game struct {
	chess.Game
	ID          ID
	Subscribers []Subscriber
}

// newGame возвращает новую игру.
// TODO: добавить возможность задавать начальную позицию, историю ходов.
func newGame(st Storage) Game {
	g := Game{
		Game: chess.Game{
			StartingPosition: *board.Classical(),
		},
	}

	err := errors.New("error")
	for err != nil {
		g.ID = newID()
		err = st.StoreGame(g)
	}

	return g
}

// Subscriber — сущность (URL), подписанная на уведомления об игре.
type Subscriber struct {
	// URL, на который будут посланы уведомления.
	URL string

	// Время, когда подписка была подтверждена в последний раз. Если с тех пор
	// прошло больше, чем предельная продолжительность, подписка может быть
	// отменена.
	LastConfirmed time.Time
}

// ID — идентификатор игры.
type ID string

// newId генерирует id для игры (уникальность здесь не проверяется).
func newID() ID {
	err := errors.New("error")
	var u uuid.UUID
	for err != nil {
		u, err = uuid.NewRandom()
	}

	return ID(u.String())
}

func (id ID) String() string {
	return string(id)
}
