package repository

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository/sqlcgen"
	"github.com/google/uuid"
)

type ArtistRepository struct {
	q *sqlcgen.Queries
}

func NewArtistRepository(db sqlcgen.DBTX) *ArtistRepository {
	return &ArtistRepository{q: sqlcgen.New(db)}
}

func (r *ArtistRepository) Create(ctx context.Context, id uuid.UUID, name, officialSiteURL string) (domain.Artist, error) {
	row, err := r.q.CreateArtist(ctx, sqlcgen.CreateArtistParams{
		ID:              toUUID(id),
		Name:            name,
		OfficialSiteUrl: officialSiteURL,
	})
	if err != nil {
		return domain.Artist{}, classifyErr(err)
	}
	return artistFromRow(row), nil
}

func (r *ArtistRepository) Get(ctx context.Context, id uuid.UUID) (domain.Artist, error) {
	row, err := r.q.GetArtist(ctx, toUUID(id))
	if err != nil {
		return domain.Artist{}, classifyErr(err)
	}
	return artistFromRow(row), nil
}

func (r *ArtistRepository) GetByURL(ctx context.Context, officialSiteURL string) (domain.Artist, error) {
	row, err := r.q.GetArtistByURL(ctx, officialSiteURL)
	if err != nil {
		return domain.Artist{}, classifyErr(err)
	}
	return artistFromRow(row), nil
}

// List returns artists whose name matches the query (empty query = all),
// paginated by limit/offset, along with the total match count.
func (r *ArtistRepository) List(ctx context.Context, query string, limit, offset int32) ([]domain.Artist, int64, error) {
	var q *string
	if query != "" {
		q = &query
	}

	rows, err := r.q.ListArtists(ctx, sqlcgen.ListArtistsParams{Query: q, Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, classifyErr(err)
	}
	total, err := r.q.CountArtists(ctx, q)
	if err != nil {
		return nil, 0, classifyErr(err)
	}

	artists := make([]domain.Artist, 0, len(rows))
	for _, row := range rows {
		artists = append(artists, artistFromRow(row))
	}
	return artists, total, nil
}

// ListAll returns every artist, e.g. for the scraper batch to iterate over.
func (r *ArtistRepository) ListAll(ctx context.Context) ([]domain.Artist, error) {
	rows, err := r.q.ListAllArtists(ctx)
	if err != nil {
		return nil, classifyErr(err)
	}
	artists := make([]domain.Artist, 0, len(rows))
	for _, row := range rows {
		artists = append(artists, artistFromRow(row))
	}
	return artists, nil
}

func (r *ArtistRepository) Update(ctx context.Context, id uuid.UUID, name, officialSiteURL string) (domain.Artist, error) {
	row, err := r.q.UpdateArtist(ctx, sqlcgen.UpdateArtistParams{
		ID:              toUUID(id),
		Name:            name,
		OfficialSiteUrl: officialSiteURL,
	})
	if err != nil {
		return domain.Artist{}, classifyErr(err)
	}
	return artistFromRow(row), nil
}

func (r *ArtistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	n, err := r.q.DeleteArtist(ctx, toUUID(id))
	if err != nil {
		return classifyErr(err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func artistFromRow(row sqlcgen.Artist) domain.Artist {
	return domain.Artist{
		ID:              fromUUID(row.ID),
		Name:            row.Name,
		OfficialSiteURL: row.OfficialSiteUrl,
		CreatedAt:       fromTimestamptz(row.CreatedAt),
	}
}
