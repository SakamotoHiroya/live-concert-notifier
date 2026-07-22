package repository

import (
	"context"
	"time"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository/sqlcgen"
	"github.com/google/uuid"
)

type ScrapeJobRepository struct {
	q *sqlcgen.Queries
}

func NewScrapeJobRepository(db sqlcgen.DBTX) *ScrapeJobRepository {
	return &ScrapeJobRepository{q: sqlcgen.New(db)}
}

func (r *ScrapeJobRepository) Create(ctx context.Context, id, artistID uuid.UUID) (domain.ScrapeJob, error) {
	row, err := r.q.CreateScrapeJob(ctx, sqlcgen.CreateScrapeJobParams{
		ID:       toUUID(id),
		ArtistID: toUUID(artistID),
		Status:   string(domain.ScrapeJobStatusPending),
	})
	if err != nil {
		return domain.ScrapeJob{}, classifyErr(err)
	}
	return domain.ScrapeJob{
		ID:       fromUUID(row.ID),
		ArtistID: fromUUID(row.ArtistID),
		Status:   domain.ScrapeJobStatus(row.Status),
	}, nil
}

func (r *ScrapeJobRepository) Get(ctx context.Context, id uuid.UUID) (domain.ScrapeJob, error) {
	row, err := r.q.GetScrapeJob(ctx, toUUID(id))
	if err != nil {
		return domain.ScrapeJob{}, classifyErr(err)
	}
	return domain.ScrapeJob{
		ID:           fromUUID(row.ID),
		ArtistID:     fromUUID(row.ArtistID),
		ArtistName:   row.ArtistName,
		Status:       domain.ScrapeJobStatus(row.Status),
		StartedAt:    fromTimestamptzPtr(row.StartedAt),
		FinishedAt:   fromTimestamptzPtr(row.FinishedAt),
		ErrorMessage: row.ErrorMessage,
	}, nil
}

// ScrapeJobFilter narrows ScrapeJobRepository.List; nil fields are unfiltered.
type ScrapeJobFilter struct {
	ArtistID *uuid.UUID
	Status   *domain.ScrapeJobStatus
}

func (r *ScrapeJobRepository) List(ctx context.Context, filter ScrapeJobFilter, limit, offset int32) ([]domain.ScrapeJob, int64, error) {
	artistID := filterUUID(filter.ArtistID)
	var status *string
	if filter.Status != nil {
		s := string(*filter.Status)
		status = &s
	}

	rows, err := r.q.ListScrapeJobs(ctx, sqlcgen.ListScrapeJobsParams{
		Limit:    limit,
		Offset:   offset,
		ArtistID: artistID,
		Status:   status,
	})
	if err != nil {
		return nil, 0, classifyErr(err)
	}
	total, err := r.q.CountScrapeJobs(ctx, sqlcgen.CountScrapeJobsParams{ArtistID: artistID, Status: status})
	if err != nil {
		return nil, 0, classifyErr(err)
	}

	jobs := make([]domain.ScrapeJob, 0, len(rows))
	for _, row := range rows {
		jobs = append(jobs, domain.ScrapeJob{
			ID:           fromUUID(row.ID),
			ArtistID:     fromUUID(row.ArtistID),
			ArtistName:   row.ArtistName,
			Status:       domain.ScrapeJobStatus(row.Status),
			StartedAt:    fromTimestamptzPtr(row.StartedAt),
			FinishedAt:   fromTimestamptzPtr(row.FinishedAt),
			ErrorMessage: row.ErrorMessage,
		})
	}
	return jobs, total, nil
}

// UpdateStatus transitions a job's status, rejecting invalid status values
// via domain.ScrapeJobStatus.Valid before writing.
func (r *ScrapeJobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ScrapeJobStatus, startedAt, finishedAt *time.Time, errMsg *string) (domain.ScrapeJob, error) {
	if !status.Valid() {
		return domain.ScrapeJob{}, domain.ErrInvalidScrapeJobStatus
	}
	row, err := r.q.UpdateScrapeJobStatus(ctx, sqlcgen.UpdateScrapeJobStatusParams{
		ID:           toUUID(id),
		Status:       string(status),
		StartedAt:    toTimestamptzPtr(startedAt),
		FinishedAt:   toTimestamptzPtr(finishedAt),
		ErrorMessage: errMsg,
	})
	if err != nil {
		return domain.ScrapeJob{}, classifyErr(err)
	}
	return domain.ScrapeJob{
		ID:           fromUUID(row.ID),
		ArtistID:     fromUUID(row.ArtistID),
		Status:       domain.ScrapeJobStatus(row.Status),
		StartedAt:    fromTimestamptzPtr(row.StartedAt),
		FinishedAt:   fromTimestamptzPtr(row.FinishedAt),
		ErrorMessage: row.ErrorMessage,
	}, nil
}
