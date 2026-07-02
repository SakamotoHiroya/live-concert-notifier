package service

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

type ConcertService struct {
	concerts *repository.ConcertRepository
}

func NewConcertService(concerts *repository.ConcertRepository) *ConcertService {
	return &ConcertService{concerts: concerts}
}

func (s *ConcertService) GetConcert(ctx context.Context, id uuid.UUID) (domain.Concert, error) {
	return s.concerts.Get(ctx, id)
}

func (s *ConcertService) ListConcerts(ctx context.Context, filter repository.ConcertFilter, limit, offset int32) ([]domain.Concert, int64, error) {
	return s.concerts.List(ctx, filter, limit, offset)
}
