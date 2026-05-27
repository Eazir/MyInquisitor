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

type CategoryRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewCategoryRepository(db *PostgresDB) *CategoryRepository {
	return &CategoryRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCCategoryParams(c *entity.Category) sqlc.CreateCategoryParams {
	return sqlc.CreateCategoryParams{
		UserID: c.UserID,
		Name:   c.Name,
		Type:   c.Type,
		Icon:   toPGText(c.Icon),
		Color:  toPGText(c.Color),
	}
}

func fromSQLCCategory(s sqlc.Category) entity.Category {
	return entity.Category{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		Type:      s.Type,
		Icon:      fromPGText(s.Icon),
		Color:     fromPGText(s.Color),
		CreatedAt: s.CreatedAt,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, c *entity.Category) error {
	created, err := r.q.CreateCategory(ctx, toSQLCCategoryParams(c))
	if err != nil {
		return err
	}
	*c = fromSQLCCategory(created)
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	s, err := r.q.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	c := fromSQLCCategory(s)
	return &c, nil
}

func (r *CategoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Category, error) {
	rows, err := r.q.ListCategoriesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Category, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCCategory(row)
	}
	return result, nil
}

func (r *CategoryRepository) ListByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]entity.Category, error) {
	rows, err := r.q.ListCategoriesByType(ctx, sqlc.ListCategoriesByTypeParams{
		UserID: userID,
		Type:   categoryType,
	})
	if err != nil {
		return nil, err
	}
	result := make([]entity.Category, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCCategory(row)
	}
	return result, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
