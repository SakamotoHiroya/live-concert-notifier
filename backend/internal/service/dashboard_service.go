package service

import (
	"context"
	"time"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

// NewBadgeWindow is how recently a concert must have been discovered to be
// flagged "is_new" on the dashboard, per docs/spec.md 4-3.
const NewBadgeWindow = 7 * 24 * time.Hour

type DashboardService struct {
	concerts *repository.ConcertRepository
	users    *repository.UserRepository
}

func NewDashboardService(concerts *repository.ConcertRepository, users *repository.UserRepository) *DashboardService {
	return &DashboardService{concerts: concerts, users: users}
}

// DashboardConcert pairs a Concert with its computed "NEW" badge state.
type DashboardConcert struct {
	domain.Concert
	IsNew bool
}

// List returns userID's upcoming concerts (across followed artists) with the
// "NEW" badge computed, or repository.ErrNotFound if userID does not exist.
func (s *DashboardService) List(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]DashboardConcert, int64, error) {
	if _, err := s.users.Get(ctx, userID); err != nil {
		return nil, 0, err
	}

	concerts, total, err := s.concerts.ListForDashboard(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	items := make([]DashboardConcert, 0, len(concerts))
	for _, c := range concerts {
		items = append(items, DashboardConcert{Concert: c, IsNew: c.IsNew(now, NewBadgeWindow)})
	}
	return items, total, nil
}
