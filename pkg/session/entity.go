package session

import (
	"log"
	"time"
)

type SessionServiceInterface interface {
	CreateNew() error
	UpdateData() error
	Find() ([]Session, error)
}

type SessionService struct {
	Data Session
}

type Session struct {
	Uid         uint64
	StartDate   time.Time
	EndDate     time.Time
	BlackPlayer uint64
	WhitePlayer uint64
	Description string
	State       uint8
	Notation    string
}

func (s *SessionService) CreateNew() error {
	log.Printf("Write to database: %+v", s.Data)
	return nil
}

func (s *SessionService) UpdateData() error {
	log.Printf("Update data in database: %+v", s.Data)
	return nil
}

func (s *SessionService) Find() ([]Session, error) {
	result := make([]Session, 0)
	log.Printf("Write to database: %+v", s.Data)
	return result, nil
}

func GetSessionService() SessionServiceInterface {
	return &SessionService{}
}
