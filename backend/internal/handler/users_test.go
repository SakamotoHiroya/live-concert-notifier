package handler_test

import (
	"context"
	"os"
	"testing"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/handler"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/service"
	"github.com/google/uuid"
)

func newTestAPIHandler(t *testing.T) *handler.APIHandler {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping handler integration test")
	}
	pool, err := repository.Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)

	userService := service.NewUserService(repository.NewUserRepository(pool))
	artistService := service.NewArtistService(repository.NewArtistRepository(pool))
	return handler.NewAPIHandler(userService, artistService)
}

func TestAPIHandler_Users(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	res, err := h.UsersPost(ctx, &oas.CreateUserRequest{Email: uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("UsersPost: %v", err)
	}
	created, ok := res.(*oas.User)
	if !ok {
		t.Fatalf("UsersPost returned %T, want *oas.User", res)
	}

	// Duplicate email must be rejected with 409, not a raw error.
	dupRes, err := h.UsersPost(ctx, &oas.CreateUserRequest{Email: created.Email})
	if err != nil {
		t.Fatalf("UsersPost (duplicate) returned error instead of a Conflict response: %v", err)
	}
	if _, ok := dupRes.(*oas.UsersPostConflict); !ok {
		t.Fatalf("UsersPost (duplicate) = %T, want *oas.UsersPostConflict", dupRes)
	}

	getRes, err := h.UsersUserIdGet(ctx, oas.UsersUserIdGetParams{UserId: created.ID})
	if err != nil {
		t.Fatalf("UsersUserIdGet: %v", err)
	}
	got, ok := getRes.(*oas.User)
	if !ok || got.Email != created.Email {
		t.Fatalf("UsersUserIdGet = %+v (%T), want user with email %q", getRes, getRes, created.Email)
	}

	notFoundRes, err := h.UsersUserIdGet(ctx, oas.UsersUserIdGetParams{UserId: uuid.New()})
	if err != nil {
		t.Fatalf("UsersUserIdGet (missing) returned error instead of a NotFound response: %v", err)
	}
	if _, ok := notFoundRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("UsersUserIdGet (missing) = %T, want *oas.ErrorResponse", notFoundRes)
	}
}
