package service

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository"
	"github.com/google/uuid"
)

type ArtistService struct {
	artists *repository.ArtistRepository
}

func NewArtistService(artists *repository.ArtistRepository) *ArtistService {
	return &ArtistService{artists: artists}
}

func (s *ArtistService) CreateArtist(ctx context.Context, name, officialSiteURL string) (domain.Artist, error) {
	return s.artists.Create(ctx, uuid.New(), name, officialSiteURL)
}

func (s *ArtistService) GetArtist(ctx context.Context, id uuid.UUID) (domain.Artist, error) {
	return s.artists.Get(ctx, id)
}

func (s *ArtistService) ListArtists(ctx context.Context, query string, limit, offset int32) ([]domain.Artist, int64, error) {
	return s.artists.List(ctx, query, limit, offset)
}

// UpdateArtist applies a partial update: nil fields keep their current value.
func (s *ArtistService) UpdateArtist(ctx context.Context, id uuid.UUID, name, officialSiteURL *string) (domain.Artist, error) {
	current, err := s.artists.Get(ctx, id)
	if err != nil {
		return domain.Artist{}, err
	}
	newName := current.Name
	if name != nil {
		newName = *name
	}
	newURL := current.OfficialSiteURL
	if officialSiteURL != nil {
		newURL = *officialSiteURL
	}
	return s.artists.Update(ctx, id, newName, newURL)
}

func (s *ArtistService) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	return s.artists.Delete(ctx, id)
}
