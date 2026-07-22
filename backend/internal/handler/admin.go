package handler

import (
	"context"
	"errors"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/oas"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

func toOASScrapeJob(job domain.ScrapeJob) *oas.ScrapeJob {
	out := &oas.ScrapeJob{
		ID:         job.ID,
		ArtistID:   job.ArtistID,
		ArtistName: job.ArtistName,
		Status:     oas.ScrapeJobStatus(job.Status),
	}
	if job.StartedAt != nil {
		out.StartedAt = oas.OptNilDateTime{Value: *job.StartedAt, Set: true}
	}
	if job.FinishedAt != nil {
		out.FinishedAt = oas.OptNilDateTime{Value: *job.FinishedAt, Set: true}
	}
	if job.ErrorMessage != nil {
		out.ErrorMessage = oas.OptNilString{Value: *job.ErrorMessage, Set: true}
	}
	return out
}

// AdminScrapePost implements POST /admin/scrape.
func (h *APIHandler) AdminScrapePost(ctx context.Context, req oas.OptTriggerScrapeRequest) (*oas.TriggerScrapeResponse, error) {
	var artistIDs []uuid.UUID
	if req.IsSet() {
		artistIDs = req.Value.ArtistIds
	}

	jobIDs, err := h.adminScrape.TriggerScrape(ctx, artistIDs)
	if err != nil {
		return nil, err
	}
	return &oas.TriggerScrapeResponse{JobIds: jobIDs}, nil
}

// AdminScrapeJobsGet implements GET /admin/scrape-jobs.
func (h *APIHandler) AdminScrapeJobsGet(ctx context.Context, params oas.AdminScrapeJobsGetParams) (*oas.ScrapeJobList, error) {
	var filter repository.ScrapeJobFilter
	if params.ArtistID.IsSet() {
		v := params.ArtistID.Value
		filter.ArtistID = &v
	}
	if params.Status.IsSet() {
		v := domain.ScrapeJobStatus(params.Status.Value)
		filter.Status = &v
	}
	limit, offset := int32(defaultLimit), int32(0)
	if params.Limit.IsSet() {
		limit = int32(params.Limit.Value)
	}
	if params.Offset.IsSet() {
		offset = int32(params.Offset.Value)
	}

	jobs, total, err := h.adminScrape.ListScrapeJobs(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}
	items := make([]oas.ScrapeJob, 0, len(jobs))
	for _, j := range jobs {
		items = append(items, *toOASScrapeJob(j))
	}
	return &oas.ScrapeJobList{Items: items, Total: int(total)}, nil
}

// AdminScrapeJobsJobIdGet implements GET /admin/scrape-jobs/{jobId}.
func (h *APIHandler) AdminScrapeJobsJobIdGet(ctx context.Context, params oas.AdminScrapeJobsJobIdGetParams) (oas.AdminScrapeJobsJobIdGetRes, error) {
	job, err := h.adminScrape.GetScrapeJob(ctx, params.JobId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &oas.ErrorResponse{Code: "NOT_FOUND", Message: "scrape job not found"}, nil
		}
		return nil, err
	}
	return toOASScrapeJob(job), nil
}
