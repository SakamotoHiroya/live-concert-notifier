package repository

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository/sqlcgen"
	"github.com/google/uuid"
)

type FollowRepository struct {
	q *sqlcgen.Queries
}

func NewFollowRepository(db sqlcgen.DBTX) *FollowRepository {
	return &FollowRepository{q: sqlcgen.New(db)}
}

func (r *FollowRepository) Follow(ctx context.Context, userID, artistID uuid.UUID) error {
	err := r.q.FollowArtist(ctx, sqlcgen.FollowArtistParams{UserID: toUUID(userID), ArtistID: toUUID(artistID)})
	return classifyErr(err)
}

func (r *FollowRepository) Unfollow(ctx context.Context, userID, artistID uuid.UUID) error {
	n, err := r.q.UnfollowArtist(ctx, sqlcgen.UnfollowArtistParams{UserID: toUUID(userID), ArtistID: toUUID(artistID)})
	if err != nil {
		return classifyErr(err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *FollowRepository) IsFollowing(ctx context.Context, userID, artistID uuid.UUID) (bool, error) {
	ok, err := r.q.IsFollowing(ctx, sqlcgen.IsFollowingParams{UserID: toUUID(userID), ArtistID: toUUID(artistID)})
	if err != nil {
		return false, classifyErr(err)
	}
	return ok, nil
}

func (r *FollowRepository) ListFollowed(ctx context.Context, userID uuid.UUID) ([]domain.Artist, error) {
	rows, err := r.q.ListFollowedArtists(ctx, toUUID(userID))
	if err != nil {
		return nil, classifyErr(err)
	}
	artists := make([]domain.Artist, 0, len(rows))
	for _, row := range rows {
		artists = append(artists, artistFromRow(row))
	}
	return artists, nil
}
