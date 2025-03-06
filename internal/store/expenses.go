package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	UserId      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseStore struct {
	db *sql.DB
}

func (s *ExpenseStore) Create(ctx context.Context, expense *Expense) error {
	query := `INSERT INTO expenses (id, title, amount, description, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, uuid.New(), expense.Title, expense.Amount, expense.Description, expense.UserId).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *ExpenseStore) List(ctx context.Context) ([]Expense, error) {
	query := `SELECT id, title, amount, user_id, description, created_at, updated_at FROM expenses`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.UserId, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}
