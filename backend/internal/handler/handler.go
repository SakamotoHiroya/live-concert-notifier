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

	users     *service.UserService
	artists   *service.ArtistService
	follows   *service.FollowService
	concerts  *service.ConcertService
	dashboard *service.DashboardService
}

func NewAPIHandler(users *service.UserService, artists *service.ArtistService, follows *service.FollowService, concerts *service.ConcertService, dashboard *service.DashboardService) *APIHandler {
	return &APIHandler{users: users, artists: artists, follows: follows, concerts: concerts, dashboard: dashboard}
}
