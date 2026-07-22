package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ErrInvalidScrapeJobStatus is returned when a status transition is
// attempted with a value outside ScrapeJobStatus's known set.
var ErrInvalidScrapeJobStatus = errors.New("domain: invalid scrape job status")

// ScrapeJobStatus is the lifecycle state of a ScrapeJob.
type ScrapeJobStatus string

const (
	ScrapeJobStatusPending   ScrapeJobStatus = "pending"
	ScrapeJobStatusRunning   ScrapeJobStatus = "running"
	ScrapeJobStatusSucceeded ScrapeJobStatus = "succeeded"
	ScrapeJobStatusFailed    ScrapeJobStatus = "failed"
)

// Valid reports whether s is one of the known ScrapeJobStatus values.
func (s ScrapeJobStatus) Valid() bool {
	switch s {
	case ScrapeJobStatusPending, ScrapeJobStatusRunning, ScrapeJobStatusSucceeded, ScrapeJobStatusFailed:
		return true
	default:
		return false
	}
}

// ScrapeJob tracks a single scraping run for an artist.
type ScrapeJob struct {
	ID           uuid.UUID
	ArtistID     uuid.UUID
	ArtistName   string
	Status       ScrapeJobStatus
	StartedAt    *time.Time
	FinishedAt   *time.Time
	ErrorMessage *string
}

// NewScrapeJob creates a pending ScrapeJob for the given artist.
func NewScrapeJob(id, artistID uuid.UUID) ScrapeJob {
	return ScrapeJob{
		ID:       id,
		ArtistID: artistID,
		Status:   ScrapeJobStatusPending,
	}
}

// SetStatus updates the status, rejecting unknown values.
func (j *ScrapeJob) SetStatus(status ScrapeJobStatus) error {
	if !status.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidScrapeJobStatus, status)
	}
	j.Status = status
	return nil
}
