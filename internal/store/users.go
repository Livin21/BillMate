package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Password  string    `json:"_"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, password, phone, email) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	id := uuid.New()
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	queryErr := s.db.QueryRowContext(ctx, query, id, user.Name, password, user.Phone, user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if queryErr != nil {
		return queryErr
	}
	return nil
}

func (s *UserStore) List(ctx context.Context) ([]User, error) {
	query := `SELECT id, name, phone, email, created_at, updated_at FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}