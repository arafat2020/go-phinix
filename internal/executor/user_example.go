package executor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arafat2020/go-phoenix/runtime/dberr"
)

// User is a placeholder model — real generated models come in Phase 4.
type User struct {
	ID    int64
	Email string
	Name  string
}

func (e *Executor) CreateUser(ctx context.Context, u *User) error {
	const q = `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`
	return e.QueryRowContext(ctx, q, u.Email, u.Name).Scan(&u.ID)
}

func (e *Executor) GetUser(ctx context.Context, id int64) (*User, error) {
	const q = `SELECT id, email, name FROM users WHERE id = $1`
	var u User
	err := e.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email, &u.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dberr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return &u, nil
}
