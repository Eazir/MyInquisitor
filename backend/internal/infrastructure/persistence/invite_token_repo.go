package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)

type InviteTokenRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewInviteTokenRepository(db *PostgresDB) *InviteTokenRepository {
	return &InviteTokenRepository{db: db, q: sqlc.New(db.Pool)}
}

func fromSQLCInviteToken(s sqlc.InviteToken) entity.InviteToken {
	return entity.InviteToken{
		ID:        s.ID,
		Token:     s.Token,
		CreatedBy: s.CreatedBy,
		Used:      s.Used,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}
}

func (r *InviteTokenRepository) Create(ctx context.Context, t *entity.InviteToken) error {
	created, err := r.q.CreateInviteToken(ctx, sqlc.CreateInviteTokenParams{
		Token:     t.Token,
		CreatedBy: t.CreatedBy,
		ExpiresAt: t.ExpiresAt,
	})
	if err != nil {
		return err
	}
	*t = fromSQLCInviteToken(created)
	return nil
}

func (r *InviteTokenRepository) GetByToken(ctx context.Context, token string) (*entity.InviteToken, error) {
	s, err := r.q.GetInviteToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	t := fromSQLCInviteToken(s)
	return &t, nil
}

func (r *InviteTokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.q.MarkInviteTokenAsUsed(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
