package repository

import (
	"context"
	"errors"
	"time"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository/sqlcgen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ConcertRepository struct {
	q *sqlcgen.Queries
}

func NewConcertRepository(db sqlcgen.DBTX) *ConcertRepository {
	return &ConcertRepository{q: sqlcgen.New(db)}
}

// NewConcert is the input for creating a Concert.
type NewConcert struct {
	ID            uuid.UUID
	ArtistID      uuid.UUID
	Title         string
	VenueName     string
	VenueLocation string
	Date          time.Time
	CoPerformers  []string
	IsFestival    bool
	SourceURL     string
	RawText       string
}

// Create inserts a concert, skipping silently (inserted=false) if a concert
// with the same (artist_id, date, venue_name) already exists.
func (r *ConcertRepository) Create(ctx context.Context, c NewConcert) (concert domain.Concert, inserted bool, err error) {
	coPerformers := c.CoPerformers
	if coPerformers == nil {
		coPerformers = []string{}
	}
	row, err := r.q.CreateConcert(ctx, sqlcgen.CreateConcertParams{
		ID:            toUUID(c.ID),
		ArtistID:      toUUID(c.ArtistID),
		Title:         c.Title,
		VenueName:     c.VenueName,
		VenueLocation: c.VenueLocation,
		Date:          toDate(c.Date),
		CoPerformers:  coPerformers,
		IsFestival:    c.IsFestival,
		SourceUrl:     c.SourceURL,
		RawText:       c.RawText,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Concert{}, false, nil
		}
		return domain.Concert{}, false, classifyErr(err)
	}
	return domain.Concert{
		ID:            fromUUID(row.ID),
		ArtistID:      fromUUID(row.ArtistID),
		Title:         row.Title,
		VenueName:     row.VenueName,
		VenueLocation: row.VenueLocation,
		Date:          fromDate(row.Date),
		CoPerformers:  row.CoPerformers,
		IsFestival:    row.IsFestival,
		SourceURL:     row.SourceUrl,
		RawText:       row.RawText,
		DiscoveredAt:  fromTimestamptz(row.DiscoveredAt),
		CreatedAt:     fromTimestamptz(row.CreatedAt),
	}, true, nil
}

func (r *ConcertRepository) Get(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	row, err := r.q.GetConcert(ctx, toUUID(id))
	if err != nil {
		return domain.Concert{}, classifyErr(err)
	}
	return domain.Concert{
		ID:            fromUUID(row.ID),
		ArtistID:      fromUUID(row.ArtistID),
		ArtistName:    row.ArtistName,
		Title:         row.Title,
		VenueName:     row.VenueName,
		VenueLocation: row.VenueLocation,
		Date:          fromDate(row.Date),
		CoPerformers:  row.CoPerformers,
		IsFestival:    row.IsFestival,
		SourceURL:     row.SourceUrl,
		RawText:       row.RawText,
		DiscoveredAt:  fromTimestamptz(row.DiscoveredAt),
		CreatedAt:     fromTimestamptz(row.CreatedAt),
	}, nil
}

// ConcertFilter narrows ConcertRepository.List; nil fields are unfiltered.
type ConcertFilter struct {
	ArtistID   *uuid.UUID
	From       *time.Time
	To         *time.Time
	IsFestival *bool
}

func (r *ConcertRepository) List(ctx context.Context, filter ConcertFilter, limit, offset int32) ([]domain.Concert, int64, error) {
	artistID, from, to := filterUUID(filter.ArtistID), filterDate(filter.From), filterDate(filter.To)

	rows, err := r.q.ListConcerts(ctx, sqlcgen.ListConcertsParams{
		Limit:      limit,
		Offset:     offset,
		ArtistID:   artistID,
		FromDate:   from,
		ToDate:     to,
		IsFestival: filter.IsFestival,
	})
	if err != nil {
		return nil, 0, classifyErr(err)
	}
	total, err := r.q.CountConcerts(ctx, sqlcgen.CountConcertsParams{
		ArtistID:   artistID,
		FromDate:   from,
		ToDate:     to,
		IsFestival: filter.IsFestival,
	})
	if err != nil {
		return nil, 0, classifyErr(err)
	}

	concerts := make([]domain.Concert, 0, len(rows))
	for _, row := range rows {
		concerts = append(concerts, domain.Concert{
			ID:            fromUUID(row.ID),
			ArtistID:      fromUUID(row.ArtistID),
			ArtistName:    row.ArtistName,
			Title:         row.Title,
			VenueName:     row.VenueName,
			VenueLocation: row.VenueLocation,
			Date:          fromDate(row.Date),
			CoPerformers:  row.CoPerformers,
			IsFestival:    row.IsFestival,
			SourceURL:     row.SourceUrl,
			RawText:       row.RawText,
			DiscoveredAt:  fromTimestamptz(row.DiscoveredAt),
			CreatedAt:     fromTimestamptz(row.CreatedAt),
		})
	}
	return concerts, total, nil
}

// ListForDashboard returns upcoming concerts for artists userID follows.
func (r *ConcertRepository) ListForDashboard(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Concert, int64, error) {
	rows, err := r.q.ListDashboardConcerts(ctx, sqlcgen.ListDashboardConcertsParams{
		UserID: toUUID(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, classifyErr(err)
	}
	total, err := r.q.CountDashboardConcerts(ctx, toUUID(userID))
	if err != nil {
		return nil, 0, classifyErr(err)
	}

	concerts := make([]domain.Concert, 0, len(rows))
	for _, row := range rows {
		concerts = append(concerts, domain.Concert{
			ID:            fromUUID(row.ID),
			ArtistID:      fromUUID(row.ArtistID),
			ArtistName:    row.ArtistName,
			Title:         row.Title,
			VenueName:     row.VenueName,
			VenueLocation: row.VenueLocation,
			Date:          fromDate(row.Date),
			CoPerformers:  row.CoPerformers,
			IsFestival:    row.IsFestival,
			SourceURL:     row.SourceUrl,
			RawText:       row.RawText,
			DiscoveredAt:  fromTimestamptz(row.DiscoveredAt),
			CreatedAt:     fromTimestamptz(row.CreatedAt),
		})
	}
	return concerts, total, nil
}

func filterUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return toUUID(*id)
}

func filterDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return toDate(*t)
}
