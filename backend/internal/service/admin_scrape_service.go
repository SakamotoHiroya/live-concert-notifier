package service

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

type AdminScrapeService struct {
	scrapeJobs *repository.ScrapeJobRepository
	artists    *repository.ArtistRepository
	trigger    ScraperTrigger
}

func NewAdminScrapeService(scrapeJobs *repository.ScrapeJobRepository, artists *repository.ArtistRepository, trigger ScraperTrigger) *AdminScrapeService {
	return &AdminScrapeService{scrapeJobs: scrapeJobs, artists: artists, trigger: trigger}
}

// TriggerScrape creates a pending ScrapeJob per target artist (all artists if
// artistIDs is empty) and asks the scraper to start, per docs/ai/architecture.md 4-2.
func (s *AdminScrapeService) TriggerScrape(ctx context.Context, artistIDs []uuid.UUID) ([]uuid.UUID, error) {
	targets := artistIDs
	if len(targets) == 0 {
		all, err := s.artists.ListAll(ctx)
		if err != nil {
			return nil, err
		}
		targets = make([]uuid.UUID, 0, len(all))
		for _, a := range all {
			targets = append(targets, a.ID)
		}
	}

	jobIDs := make([]uuid.UUID, 0, len(targets))
	for _, artistID := range targets {
		job, err := s.scrapeJobs.Create(ctx, uuid.New(), artistID)
		if err != nil {
			return nil, err
		}
		jobIDs = append(jobIDs, job.ID)
	}

	if err := s.trigger.Trigger(ctx, targets); err != nil {
		return nil, err
	}
	return jobIDs, nil
}

func (s *AdminScrapeService) GetScrapeJob(ctx context.Context, id uuid.UUID) (domain.ScrapeJob, error) {
	return s.scrapeJobs.Get(ctx, id)
}

func (s *AdminScrapeService) ListScrapeJobs(ctx context.Context, filter repository.ScrapeJobFilter, limit, offset int32) ([]domain.ScrapeJob, int64, error) {
	return s.scrapeJobs.List(ctx, filter, limit, offset)
}
