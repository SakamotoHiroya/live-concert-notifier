package service

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// ScraperTrigger starts the scraper batch (Cloud Run Job in production) for
// the given artists. Implementations are swapped per environment: a real
// Cloud Run Admin API client in production, a no-op/log stub locally.
type ScraperTrigger interface {
	Trigger(ctx context.Context, artistIDs []uuid.UUID) error
}

// LogScraperTrigger logs the trigger instead of calling any external system.
// Used for local development until the Cloud Run Admin API client exists.
type LogScraperTrigger struct{}

func (LogScraperTrigger) Trigger(ctx context.Context, artistIDs []uuid.UUID) error {
	log.Printf("service: scraper trigger requested for %d artist(s): %v", len(artistIDs), artistIDs)
	return nil
}
