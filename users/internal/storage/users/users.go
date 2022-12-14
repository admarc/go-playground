package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/admarc/users/internal/models"
	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{db: db}
}

func (s Storage) Create(ctx context.Context, name string) (models.User, error) {
	id := uuid.NewString()
	_, err := s.db.ExecContext(ctx, "INSERT into users (id, name) values (?,?)", id, name)
	if err != nil {
		return models.User{}, fmt.Errorf("Failed to execute insert %w", err)
	}
	return models.User{ID: id, Name: name}, nil
}

func (s Storage) Get(ctx context.Context, id string) (models.User, error) {
	var name string
	fmt.Println(id)
	err := s.db.QueryRowContext(ctx, "select u.name from users as u where u.id = :id;", sql.Named("id", id)).Scan(&name)
	if err != nil {
		return models.User{}, fmt.Errorf("Failed to fetch user %w", err)
	}
	return models.User{ID: id, Name: name}, nil
}
