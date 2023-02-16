package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	cfg         *Config
	database    *sql.DB
	sessionRepo *SessionRepository
	userRepo    *UserRepository
}

func CreateNewStore(c *Config) *Store {
	return &Store{
		cfg: c,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("sqlite3", s.cfg.DatabaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		return err
	}

	s.database = db
	return nil
}

func (s *Store) Close() {
	s.database.Close()
}

func (s *Store) User() *UserRepository {
	if s.userRepo == nil {
		s.userRepo = &UserRepository{store: s}
	}
	return s.userRepo
}

func (s *Store) Session() *SessionRepository {
	if s.sessionRepo == nil {
		s.sessionRepo = &SessionRepository{store: s}
	}
	return s.sessionRepo
}
