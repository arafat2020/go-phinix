package executor

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arafat2020/go-phoenix/runtime/dberr"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening sqlmock: %v", err)
	}
	defer db.Close()

	exec := &Executor{db: db, Querier: db}
	ctx := context.Background()

	u := &User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(`INSERT INTO users \(email, name\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("test@example.com", "Test User").
		WillReturnRows(rows)

	err = exec.CreateUser(ctx, u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.ID != 1 {
		t.Errorf("expected user ID 1, got %d", u.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening sqlmock: %v", err)
	}
	defer db.Close()

	exec := &Executor{db: db, Querier: db}
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "email", "name"}).
		AddRow(1, "test@example.com", "Test User")

	mock.ExpectQuery(`SELECT id, email, name FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	u, err := exec.GetUser(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.ID != 1 || u.Email != "test@example.com" || u.Name != "Test User" {
		t.Errorf("unexpected user result: %+v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening sqlmock: %v", err)
	}
	defer db.Close()

	exec := &Executor{db: db, Querier: db}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, email, name FROM users WHERE id = \$1`).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	u, err := exec.GetUser(ctx, 99)
	if !errors.Is(err, dberr.ErrNotFound) {
		t.Fatalf("expected dberr.ErrNotFound, got: %v", err)
	}
	if u != nil {
		t.Errorf("expected user to be nil, got: %+v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
