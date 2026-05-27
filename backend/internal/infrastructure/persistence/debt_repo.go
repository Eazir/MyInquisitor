package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)

type DebtRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewDebtRepository(db *PostgresDB) *DebtRepository {
	return &DebtRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCDebtParams(d *entity.Debt) sqlc.CreateDebtParams {
	return sqlc.CreateDebtParams{
		UserID:             d.UserID,
		CategoryID:         toPGUUID(d.CategoryID),
		Name:               d.Name,
		Description:        toPGText(d.Description),
		TotalAmount:        toPGNumeric(d.TotalAmount),
		RemainingAmount:    toPGNumeric(d.RemainingAmount),
		InterestRate:       toPGNumeric(d.InterestRate),
		TotalInstallments:  d.TotalInstallments,
		CurrentInstallment: d.CurrentInstallment,
		Status:             d.Status,
		StartDate:          toPGDate(d.StartDate),
		EndDate:            toPGDatePtr(d.EndDate),
	}
}

func fromSQLCDebt(s sqlc.Debt) entity.Debt {
	return entity.Debt{
		ID:                 s.ID,
		UserID:             s.UserID,
		CategoryID:         fromPGUUID(s.CategoryID),
		Name:               s.Name,
		Description:        fromPGText(s.Description),
		TotalAmount:        fromPGNumeric(s.TotalAmount),
		RemainingAmount:    fromPGNumeric(s.RemainingAmount),
		InterestRate:       fromPGNumeric(s.InterestRate),
		TotalInstallments:  s.TotalInstallments,
		CurrentInstallment: s.CurrentInstallment,
		Status:             s.Status,
		StartDate:          fromPGDate(s.StartDate),
		EndDate:            fromPGDatePtr(s.EndDate),
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
}

func (r *DebtRepository) Create(ctx context.Context, d *entity.Debt) error {
	created, err := r.q.CreateDebt(ctx, toSQLCDebtParams(d))
	if err != nil {
		return err
	}
	*d = fromSQLCDebt(created)
	return nil
}

func (r *DebtRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Debt, error) {
	s, err := r.q.GetDebtByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	d := fromSQLCDebt(s)
	return &d, nil
}

func (r *DebtRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Debt, error) {
	rows, err := r.q.ListDebtsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Debt, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCDebt(row)
	}
	return result, nil
}

func (r *DebtRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Debt, error) {
	rows, err := r.q.ListActiveDebtsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Debt, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCDebt(row)
	}
	return result, nil
}

func (r *DebtRepository) Update(ctx context.Context, d *entity.Debt) error {
	updated, err := r.q.UpdateDebt(ctx, sqlc.UpdateDebtParams{
		ID:                 d.ID,
		Name:               d.Name,
		Description:        toPGText(d.Description),
		TotalAmount:        toPGNumeric(d.TotalAmount),
		RemainingAmount:    toPGNumeric(d.RemainingAmount),
		InterestRate:       toPGNumeric(d.InterestRate),
		TotalInstallments:  d.TotalInstallments,
		CurrentInstallment: d.CurrentInstallment,
		Status:             d.Status,
		EndDate:            toPGDatePtr(d.EndDate),
	})
	if err != nil {
		return err
	}
	*d = fromSQLCDebt(updated)
	return nil
}

func (r *DebtRepository) UpdateCurrentInstallment(ctx context.Context, id uuid.UUID, installment int32) error {
	_, err := r.q.UpdateDebtCurrentInstallment(ctx, sqlc.UpdateDebtCurrentInstallmentParams{
		ID:                 id,
		CurrentInstallment: installment,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *DebtRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteDebt(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

type DebtMonthlyStatusRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewDebtMonthlyStatusRepository(db *PostgresDB) *DebtMonthlyStatusRepository {
	return &DebtMonthlyStatusRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCDebtMonthlyStatusParams(s *entity.DebtMonthlyStatus) sqlc.CreateDebtMonthlyStatusParams {
	return sqlc.CreateDebtMonthlyStatusParams{
		DebtID:            s.DebtID,
		Month:             toPGDate(s.Month),
		InstallmentNum:    s.InstallmentNum,
		TotalInstallments: s.TotalInstallments,
		AmountDue:         toPGNumeric(s.AmountDue),
		InterestAmount:    toPGNumeric(s.InterestAmount),
		PrincipalAmount:   toPGNumeric(s.PrincipalAmount),
	}
}

func fromSQLCDebtMonthlyStatus(s sqlc.DebtMonthlyStatus) entity.DebtMonthlyStatus {
	return entity.DebtMonthlyStatus{
		ID:                s.ID,
		DebtID:            s.DebtID,
		Month:             fromPGDate(s.Month),
		InstallmentNum:    s.InstallmentNum,
		TotalInstallments: s.TotalInstallments,
		AmountDue:         fromPGNumeric(s.AmountDue),
		InterestAmount:    fromPGNumeric(s.InterestAmount),
		PrincipalAmount:   fromPGNumeric(s.PrincipalAmount),
		AmountPaid:        fromPGNumeric(s.AmountPaid),
		Paid:              s.Paid,
		PaidAt:            fromPGTimestamptz(s.PaidAt),
		Notes:             fromPGText(s.Notes),
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func (r *DebtMonthlyStatusRepository) Create(ctx context.Context, s *entity.DebtMonthlyStatus) error {
	created, err := r.q.CreateDebtMonthlyStatus(ctx, toSQLCDebtMonthlyStatusParams(s))
	if err != nil {
		return err
	}
	*s = fromSQLCDebtMonthlyStatus(created)
	return nil
}

func (r *DebtMonthlyStatusRepository) GetByDebtIDAndMonth(ctx context.Context, debtID uuid.UUID, month string) (*entity.DebtMonthlyStatus, error) {
	var t pgtype.Date
	if err := t.Scan(month); err != nil {
		return nil, err
	}
	s, err := r.q.GetDebtMonthlyStatus(ctx, sqlc.GetDebtMonthlyStatusParams{
		DebtID: debtID,
		Month:  t,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	d := fromSQLCDebtMonthlyStatus(s)
	return &d, nil
}

func (r *DebtMonthlyStatusRepository) ListByDebtID(ctx context.Context, debtID uuid.UUID) ([]entity.DebtMonthlyStatus, error) {
	rows, err := r.q.ListDebtMonthlyStatusByDebtID(ctx, debtID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.DebtMonthlyStatus, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCDebtMonthlyStatus(row)
	}
	return result, nil
}

func (r *DebtMonthlyStatusRepository) MarkAsPaid(ctx context.Context, debtID uuid.UUID, month string, amountPaid float64, notes *string) error {
	var t pgtype.Date
	if err := t.Scan(month); err != nil {
		return err
	}
	_, err := r.q.MarkDebtMonthAsPaid(ctx, sqlc.MarkDebtMonthAsPaidParams{
		DebtID:     debtID,
		Month:      t,
		AmountPaid: toPGNumeric(amountPaid),
		Notes:      toPGText(notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *DebtMonthlyStatusRepository) Update(ctx context.Context, s *entity.DebtMonthlyStatus) error {
	updated, err := r.q.UpdateDebtMonthlyStatus(ctx, sqlc.UpdateDebtMonthlyStatusParams{
		DebtID:     s.DebtID,
		Month:      toPGDate(s.Month),
		AmountPaid: toPGNumeric(s.AmountPaid),
		Paid:       s.Paid,
		Notes:      toPGText(s.Notes),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	*s = fromSQLCDebtMonthlyStatus(updated)
	return nil
}

func (r *DebtMonthlyStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteDebtMonthlyStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
