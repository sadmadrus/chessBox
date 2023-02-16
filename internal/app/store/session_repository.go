package store

import "github.com/sadmadrus/chessBox/internal/app/model"

type SessionRepository struct {
	store *Store
}

func (sr *SessionRepository) Create(s *model.Session) (*model.Session, error) {
	return nil, nil
}

func (sr *SessionRepository) EndGame(s *model.Session) error {
	return nil
}
