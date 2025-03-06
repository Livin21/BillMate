package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Expenses interface {
		Create(context.Context, *Expense) error
		List(context.Context) ([]Expense, error)
	}
	Users interface {
		Create(context.Context, *User) error
		List(context.Context) ([]User, error)
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, string) (*User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Expenses: &ExpenseStore{db},
		Users:    &UserStore{db},
	}
}
