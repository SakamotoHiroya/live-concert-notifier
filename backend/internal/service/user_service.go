package service

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	users *repository.UserRepository
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) CreateUser(ctx context.Context, email string) (domain.User, error) {
	return s.users.Create(ctx, uuid.New(), email)
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.users.Get(ctx, id)
}
