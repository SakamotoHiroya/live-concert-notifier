package handler

import (
	"context"
	"errors"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
)

// UsersPost implements POST /users.
func (h *APIHandler) UsersPost(ctx context.Context, req *oas.CreateUserRequest) (oas.UsersPostRes, error) {
	user, err := h.users.CreateUser(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return &oas.UsersPostConflict{Code: "CONFLICT", Message: "email already registered"}, nil
		}
		return nil, err
	}
	return &oas.User{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt}, nil
}

// UsersUserIdGet implements GET /users/{userId}.
func (h *APIHandler) UsersUserIdGet(ctx context.Context, params oas.UsersUserIdGetParams) (oas.UsersUserIdGetRes, error) {
	user, err := h.users.GetUser(ctx, params.UserId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "user not found"}, nil
		}
		return nil, err
	}
	return &oas.User{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt}, nil
}
