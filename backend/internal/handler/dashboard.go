package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/service"
)

func toOASDashboardConcert(dc service.DashboardConcert) (*oas.DashboardConcert, error) {
	u, err := url.Parse(dc.SourceURL)
	if err != nil {
		return nil, err
	}
	return &oas.DashboardConcert{
		ID:            dc.ID,
		ArtistID:      dc.ArtistID,
		ArtistName:    dc.ArtistName,
		Title:         oas.OptString{Value: dc.Title, Set: true},
		VenueName:     dc.VenueName,
		VenueLocation: dc.VenueLocation,
		Date:          dc.Date,
		CoPerformers:  dc.CoPerformers,
		IsFestival:    dc.IsFestival,
		SourceURL:     *u,
		DiscoveredAt:  dc.DiscoveredAt,
		CreatedAt:     dc.CreatedAt,
		IsNew:         dc.IsNew,
	}, nil
}

// UsersUserIdDashboardGet implements GET /users/{userId}/dashboard.
func (h *APIHandler) UsersUserIdDashboardGet(ctx context.Context, params oas.UsersUserIdDashboardGetParams) (oas.UsersUserIdDashboardGetRes, error) {
	limit, offset := int32(defaultLimit), int32(0)
	if params.Limit.IsSet() {
		limit = int32(params.Limit.Value)
	}
	if params.Offset.IsSet() {
		offset = int32(params.Offset.Value)
	}

	concerts, total, err := h.dashboard.List(ctx, params.UserId, limit, offset)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "user not found"}, nil
		}
		return nil, err
	}

	items := make([]oas.DashboardConcert, 0, len(concerts))
	for _, c := range concerts {
		oc, err := toOASDashboardConcert(c)
		if err != nil {
			return nil, err
		}
		items = append(items, *oc)
	}
	return &oas.DashboardResponse{Items: items, Total: int(total)}, nil
}
