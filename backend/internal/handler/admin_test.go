package handler_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/google/uuid"
)

func TestAPIHandler_AdminScrape(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	siteURL, _ := url.Parse("https://example.com/" + uuid.NewString())
	artistRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Scrape Target", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost: %v", err)
	}
	artist := artistRes.(*oas.Artist)

	triggerRes, err := h.AdminScrapePost(ctx, oas.OptTriggerScrapeRequest{
		Value: oas.TriggerScrapeRequest{ArtistIds: []uuid.UUID{artist.ID}},
		Set:   true,
	})
	if err != nil {
		t.Fatalf("AdminScrapePost: %v", err)
	}
	if len(triggerRes.JobIds) != 1 {
		t.Fatalf("AdminScrapePost job_ids = %v, want 1 entry", triggerRes.JobIds)
	}
	jobID := triggerRes.JobIds[0]

	getRes, err := h.AdminScrapeJobsJobIdGet(ctx, oas.AdminScrapeJobsJobIdGetParams{JobId: jobID})
	if err != nil {
		t.Fatalf("AdminScrapeJobsJobIdGet: %v", err)
	}
	job, ok := getRes.(*oas.ScrapeJob)
	if !ok {
		t.Fatalf("AdminScrapeJobsJobIdGet = %T, want *oas.ScrapeJob", getRes)
	}
	if job.Status != oas.ScrapeJobStatusPending || job.ArtistID != artist.ID {
		t.Fatalf("AdminScrapeJobsJobIdGet = %+v, want pending status for artist %s", job, artist.ID)
	}

	listRes, err := h.AdminScrapeJobsGet(ctx, oas.AdminScrapeJobsGetParams{
		ArtistID: oas.OptUUID{Value: artist.ID, Set: true},
	})
	if err != nil {
		t.Fatalf("AdminScrapeJobsGet: %v", err)
	}
	if listRes.Total != 1 {
		t.Fatalf("AdminScrapeJobsGet total = %d, want 1", listRes.Total)
	}

	// Omitting the request body (or artist_ids) targets every artist.
	allRes, err := h.AdminScrapePost(ctx, oas.OptTriggerScrapeRequest{})
	if err != nil {
		t.Fatalf("AdminScrapePost (no body): %v", err)
	}
	if len(allRes.JobIds) < 1 {
		t.Fatalf("AdminScrapePost (no body) job_ids = %v, want at least 1", allRes.JobIds)
	}

	notFoundRes, err := h.AdminScrapeJobsJobIdGet(ctx, oas.AdminScrapeJobsJobIdGetParams{JobId: uuid.New()})
	if err != nil {
		t.Fatalf("AdminScrapeJobsJobIdGet (missing) returned error instead of NotFound: %v", err)
	}
	if _, ok := notFoundRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("AdminScrapeJobsJobIdGet (missing) = %T, want *oas.ErrorResponse", notFoundRes)
	}
}
