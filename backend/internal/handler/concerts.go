package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
)

func toOASConcert(c domain.Concert) (*oas.Concert, error) {
	u, err := url.Parse(c.SourceURL)
	if err != nil {
		return nil, err
	}
	return &oas.Concert{
		ID:            c.ID,
		ArtistID:      c.ArtistID,
		ArtistName:    c.ArtistName,
		Title:         oas.OptString{Value: c.Title, Set: true},
		VenueName:     c.VenueName,
		VenueLocation: c.VenueLocation,
		Date:          c.Date,
		CoPerformers:  c.CoPerformers,
		IsFestival:    c.IsFestival,
		SourceURL:     *u,
		DiscoveredAt:  c.DiscoveredAt,
		CreatedAt:     c.CreatedAt,
	}, nil
}

func concertFilterFromParams(artistID oas.OptUUID, from, to oas.OptDate, isFestival oas.OptBool) repository.ConcertFilter {
	var filter repository.ConcertFilter
	if artistID.IsSet() {
		v := artistID.Value
		filter.ArtistID = &v
	}
	if from.IsSet() {
		v := from.Value
		filter.From = &v
	}
	if to.IsSet() {
		v := to.Value
		filter.To = &v
	}
	if isFestival.IsSet() {
		v := isFestival.Value
		filter.IsFestival = &v
	}
	return filter
}

// ConcertsGet implements GET /concerts.
func (h *APIHandler) ConcertsGet(ctx context.Context, params oas.ConcertsGetParams) (*oas.ConcertList, error) {
	filter := concertFilterFromParams(params.ArtistID, params.From, params.To, params.IsFestival)
	limit, offset := int32(defaultLimit), int32(0)
	if params.Limit.IsSet() {
		limit = int32(params.Limit.Value)
	}
	if params.Offset.IsSet() {
		offset = int32(params.Offset.Value)
	}

	concerts, total, err := h.concerts.ListConcerts(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	items := make([]oas.Concert, 0, len(concerts))
	for _, c := range concerts {
		oc, err := toOASConcert(c)
		if err != nil {
			return nil, err
		}
		items = append(items, *oc)
	}
	return &oas.ConcertList{Items: items, Total: int(total)}, nil
}

// ConcertsConcertIdGet implements GET /concerts/{concertId}.
func (h *APIHandler) ConcertsConcertIdGet(ctx context.Context, params oas.ConcertsConcertIdGetParams) (oas.ConcertsConcertIdGetRes, error) {
	concert, err := h.concerts.GetConcert(ctx, params.ConcertId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "concert not found"}, nil
		}
		return nil, err
	}
	return toOASConcert(concert)
}
