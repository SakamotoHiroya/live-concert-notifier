package handler

import (
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/service"
)

// APIHandler implements oas.Handler. Methods are added incrementally,
// one file per resource (users.go, artists.go, ...); unimplemented
// operations fall back to oas.UnimplementedHandler.
type APIHandler struct {
	oas.UnimplementedHandler

	users   *service.UserService
	artists *service.ArtistService
}

func NewAPIHandler(users *service.UserService, artists *service.ArtistService) *APIHandler {
	return &APIHandler{users: users, artists: artists}
}
