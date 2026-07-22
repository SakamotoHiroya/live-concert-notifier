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

func TestAPIHandler_Dashboard(t *testing.T) {
	h := newTestAPIHandler(t)
	ctx := context.Background()

	userRes, err := h.UsersPost(ctx, &oas.CreateUserRequest{Email: uuid.NewString() + "@example.com"})
	if err != nil {
		t.Fatalf("UsersPost: %v", err)
	}
	user := userRes.(*oas.User)

	siteURL, _ := url.Parse("https://example.com/" + uuid.NewString())
	artistRes, err := h.ArtistsPost(ctx, &oas.CreateArtistRequest{Name: "Dashboard Artist", OfficialSiteURL: *siteURL})
	if err != nil {
		t.Fatalf("ArtistsPost: %v", err)
	}
	artist := artistRes.(*oas.Artist)

	if _, err := h.UsersUserIdFollowsPost(ctx, &oas.FollowArtistRequest{ArtistID: artist.ID}, oas.UsersUserIdFollowsPostParams{UserId: user.ID}); err != nil {
		t.Fatalf("UsersUserIdFollowsPost: %v", err)
	}

	pool, err := repository.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	concertRepo := repository.NewConcertRepository(pool)

	upcoming := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)
	if _, inserted, err := concertRepo.Create(ctx, repository.NewConcert{
		ID:            uuid.New(),
		ArtistID:      artist.ID,
		VenueName:     "Dashboard Hall",
		VenueLocation: "Tokyo",
		Date:          upcoming,
		SourceURL:     "https://example.com/schedule",
	}); err != nil || !inserted {
		t.Fatalf("seed concert: inserted=%v err=%v", inserted, err)
	}

	res, err := h.UsersUserIdDashboardGet(ctx, oas.UsersUserIdDashboardGetParams{UserId: user.ID})
	if err != nil {
		t.Fatalf("UsersUserIdDashboardGet: %v", err)
	}
	dashboard, ok := res.(*oas.DashboardResponse)
	if !ok {
		t.Fatalf("UsersUserIdDashboardGet = %T, want *oas.DashboardResponse", res)
	}
	if dashboard.Total != 1 || len(dashboard.Items) != 1 {
		t.Fatalf("dashboard = %d items (total %d), want 1", len(dashboard.Items), dashboard.Total)
	}
	if !dashboard.Items[0].IsNew {
		t.Fatalf("dashboard item IsNew = false, want true (just discovered)")
	}

	notFoundRes, err := h.UsersUserIdDashboardGet(ctx, oas.UsersUserIdDashboardGetParams{UserId: uuid.New()})
	if err != nil {
		t.Fatalf("UsersUserIdDashboardGet (missing user) returned error instead of NotFound: %v", err)
	}
	if _, ok := notFoundRes.(*oas.ErrorResponse); !ok {
		t.Fatalf("UsersUserIdDashboardGet (missing user) = %T, want *oas.ErrorResponse", notFoundRes)
	}
}
