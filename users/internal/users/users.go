package users

import (
	"context"
	"fmt"

	"github.com/admarc/users/internal/models"
)

//go:generate moq -rm -out repository_mock.go . Repository
type Repository interface {
	Create(ctx context.Context, name string) (models.User, error)
	Get(ctx context.Context, id string) (models.User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Create(ctx context.Context, name string) (models.User, error) {
	if name == "" {
		return models.User{}, fmt.Errorf("invalid name argument: %w", models.UserCreateParamInvalidNameErr)
	}

	usr, err := s.repo.Create(ctx, name)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return usr, nil
}

func (s Service) Get(ctx context.Context, id string) (models.User, error) {

	usr, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return usr, nil
}
