package store

import (
	"github.com/sadmadrus/chessBox/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) Create(u *model.User) (*model.User, error) {

	return nil, nil
}

func (ur *UserRepository) FindByEmail(email string) *model.User {
	return nil
}
