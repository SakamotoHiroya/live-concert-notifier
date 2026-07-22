package service

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

type FollowService struct {
	follows *repository.FollowRepository
	users   *repository.UserRepository
	artists *repository.ArtistRepository
}

func NewFollowService(follows *repository.FollowRepository, users *repository.UserRepository, artists *repository.ArtistRepository) *FollowService {
	return &FollowService{follows: follows, users: users, artists: artists}
}

// Follow makes userID follow artistID, returning repository.ErrNotFound if
// either does not exist and repository.ErrConflict if already following.
func (s *FollowService) Follow(ctx context.Context, userID, artistID uuid.UUID) error {
	if _, err := s.users.Get(ctx, userID); err != nil {
		return err
	}
	if _, err := s.artists.Get(ctx, artistID); err != nil {
		return err
	}
	return s.follows.Follow(ctx, userID, artistID)
}

// Unfollow removes the follow relationship, returning repository.ErrNotFound
// if it did not exist.
func (s *FollowService) Unfollow(ctx context.Context, userID, artistID uuid.UUID) error {
	return s.follows.Unfollow(ctx, userID, artistID)
}

// ListFollowed returns userID's followed artists, or repository.ErrNotFound
// if userID does not exist.
func (s *FollowService) ListFollowed(ctx context.Context, userID uuid.UUID) ([]domain.Artist, error) {
	if _, err := s.users.Get(ctx, userID); err != nil {
		return nil, err
	}
	return s.follows.ListFollowed(ctx, userID)
}
