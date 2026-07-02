package repository

import (
	"context"

	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/domain"
	"github.com/SakamotoHiroya/live-concert-notifier/backend/internal/repository/sqlcgen"
	"github.com/google/uuid"
)

type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(db sqlcgen.DBTX) *UserRepository {
	return &UserRepository{q: sqlcgen.New(db)}
}

func (r *UserRepository) Create(ctx context.Context, id uuid.UUID, email string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, sqlcgen.CreateUserParams{ID: toUUID(id), Email: email})
	if err != nil {
		return domain.User{}, classifyErr(err)
	}
	return userFromRow(row), nil
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := r.q.GetUser(ctx, toUUID(id))
	if err != nil {
		return domain.User{}, classifyErr(err)
	}
	return userFromRow(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, classifyErr(err)
	}
	return userFromRow(row), nil
}

func userFromRow(row sqlcgen.User) domain.User {
	return domain.User{
		ID:        fromUUID(row.ID),
		Email:     row.Email,
		CreatedAt: fromTimestamptz(row.CreatedAt),
	}
}
