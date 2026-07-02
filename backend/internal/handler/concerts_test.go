package handler_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

func TestAPIHandler_Concerts(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	siteURL, _ := url.Parse("https://example.com/" + uuid.NewString())
	artistRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Concert Artist", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost: %v", err)
	}
	artist := artistRes.(*oas.Artist)

	// Concert creation isn't exposed over HTTP yet (that's the scraper's job,
	// #12), so seed directly through the repository.
	pool, err := repository.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	concertRepo := repository.NewConcertRepository(pool)

	date := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
	created, inserted, err := concertRepo.Create(ctx, repository.NewConcert{
		ID:            uuid.New(),
		ArtistID:      artist.ID,
		Title:         "TOUR 2026",
		VenueName:     "Test Hall",
		VenueLocation: "Tokyo",
		Date:          date,
		CoPerformers:  []string{"Opener"},
		IsFestival:    true,
		SourceURL:     "https://example.com/schedule",
	})
	if err != nil || !inserted {
		t.Fatalf("seed concert: %v, %v, %v", created, inserted, err)
	}

	getRes, err := h.ConcertsConcertIdGet(ctx, oas.ConcertsConcertIdGetParams{ConcertId: created.ID})
	if err != nil {
		t.Fatalf("ConcertsConcertIdGet: %v", err)
	}
	got, ok := getRes.(*oas.Concert)
	if !ok || got.ArtistName != "Concert Artist" {
		t.Fatalf("ConcertsConcertIdGet = %+v (%T), want artist_name Concert Artist", getRes, getRes)
	}

	listRes, err := h.ConcertsGet(ctx, oas.ConcertsGetParams{
		ArtistID:   oas.OptUUID{Value: artist.ID, Set: true},
		IsFestival: oas.OptBool{Value: true, Set: true},
	})
	if err != nil {
		t.Fatalf("ConcertsGet: %v", err)
	}
	if listRes.Total != 1 {
		t.Fatalf("ConcertsGet total = %d, want 1", listRes.Total)
	}

	emptyRes, err := h.ConcertsGet(ctx, oas.ConcertsGetParams{
		ArtistID:   oas.OptUUID{Value: artist.ID, Set: true},
		IsFestival: oas.OptBool{Value: false, Set: true},
	})
	if err != nil {
		t.Fatalf("ConcertsGet (is_festival=false): %v", err)
	}
	if emptyRes.Total != 0 {
		t.Fatalf("ConcertsGet (is_festival=false) total = %d, want 0", emptyRes.Total)
	}

	notFoundRes, err := h.ConcertsConcertIdGet(ctx, oas.ConcertsConcertIdGetParams{ConcertId: uuid.New()})
	if err != nil {
		t.Fatalf("ConcertsConcertIdGet (missing) returned error instead of NotFound: %v", err)
	}
	if _, ok := notFoundRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("ConcertsConcertIdGet (missing) = %T, want *oas.ErrorResponse", notFoundRes)
	}
}
