package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// requirePool skips the test unless TEST_DATABASE_URL points at a reachable
// PostgreSQL instance with the migrations in backend/migrations applied.
func requirePool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping repository integration test")
	}
	pool, err := repository.Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func TestRepositories_CRUD(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()

	users := repository.NewUserRepository(pool)
	artists := repository.NewArtistRepository(pool)
	follows := repository.NewFollowRepository(pool)
	concerts := repository.NewConcertRepository(pool)
	jobs := repository.NewScrapeJobRepository(pool)

	user, err := users.Create(ctx, uuid.New(), uuid.NewString()+"@example.com")
	if err != nil {
		t.Fatalf("users.Create: %v", err)
	}
	if _, err := users.Create(ctx, uuid.New(), user.Email); err == nil {
		t.Fatal("users.Create with duplicate email did not error")
	} else if err != repository.ErrConflict {
		t.Fatalf("users.Create duplicate email err = %v, want ErrConflict", err)
	}

	artist, err := artists.Create(ctx, uuid.New(), "Test Artist", "https://example.com/"+uuid.NewString())
	if err != nil {
		t.Fatalf("artists.Create: %v", err)
	}

	if err := follows.Follow(ctx, user.ID, artist.ID); err != nil {
		t.Fatalf("follows.Follow: %v", err)
	}
	following, err := follows.IsFollowing(ctx, user.ID, artist.ID)
	if err != nil || !following {
		t.Fatalf("follows.IsFollowing = %v, %v, want true, nil", following, err)
	}

	date := time.Now().AddDate(0, 0, 30).Truncate(24 * time.Hour)
	newConcert := repository.NewConcert{
		ID:            uuid.New(),
		ArtistID:      artist.ID,
		Title:         "ARENA TOUR 2026",
		VenueName:     "さいたまスーパーアリーナ",
		VenueLocation: "埼玉県",
		Date:          date,
		CoPerformers:  []string{"Vaundy"},
		IsFestival:    false,
		SourceURL:     "https://example.com/schedule",
	}
	concert, inserted, err := concerts.Create(ctx, newConcert)
	if err != nil || !inserted {
		t.Fatalf("concerts.Create = %v, %v, %v, want inserted", concert, inserted, err)
	}

	// Duplicate (artist_id, date, venue_name) must be skipped, not error.
	dup, inserted, err := concerts.Create(ctx, newConcert)
	if err != nil {
		t.Fatalf("concerts.Create duplicate returned error: %v", err)
	}
	if inserted {
		t.Fatalf("concerts.Create duplicate reported inserted=true, want false")
	}
	if dup.ID != uuid.Nil {
		t.Fatalf("concerts.Create duplicate returned non-zero concert: %+v", dup)
	}

	got, err := concerts.Get(ctx, concert.ID)
	if err != nil {
		t.Fatalf("concerts.Get: %v", err)
	}
	if got.ArtistName != artist.Name {
		t.Fatalf("concerts.Get artist_name = %q, want %q", got.ArtistName, artist.Name)
	}

	dashboard, total, err := concerts.ListForDashboard(ctx, user.ID, 20, 0)
	if err != nil {
		t.Fatalf("concerts.ListForDashboard: %v", err)
	}
	if total != 1 || len(dashboard) != 1 {
		t.Fatalf("concerts.ListForDashboard = %d items (total %d), want 1", len(dashboard), total)
	}

	job, err := jobs.Create(ctx, uuid.New(), artist.ID)
	if err != nil {
		t.Fatalf("jobs.Create: %v", err)
	}
	if job.Status != domain.ScrapeJobStatusPending {
		t.Fatalf("jobs.Create status = %q, want pending", job.Status)
	}
	now := time.Now()
	updated, err := jobs.UpdateStatus(ctx, job.ID, domain.ScrapeJobStatusSucceeded, &now, &now, nil)
	if err != nil {
		t.Fatalf("jobs.UpdateStatus: %v", err)
	}
	if updated.Status != domain.ScrapeJobStatusSucceeded {
		t.Fatalf("jobs.UpdateStatus status = %q, want succeeded", updated.Status)
	}

	if err := follows.Unfollow(ctx, user.ID, artist.ID); err != nil {
		t.Fatalf("follows.Unfollow: %v", err)
	}
	if err := artists.Delete(ctx, artist.ID); err != nil {
		t.Fatalf("artists.Delete: %v", err)
	}
}
