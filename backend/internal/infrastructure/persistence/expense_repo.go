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

type RecurringExpenseRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewRecurringExpenseRepository(db *PostgresDB) *RecurringExpenseRepository {
	return &RecurringExpenseRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCRecurringExpenseParams(e *entity.RecurringExpense) sqlc.CreateRecurringExpenseParams {
	dueDay := pgtype.Int4{Valid: false}
	if e.DueDay != nil {
		dueDay = pgtype.Int4{Int32: *e.DueDay, Valid: true}
	}
	billingMonth := pgtype.Int4{Valid: false}
	if e.BillingMonth != nil {
		billingMonth = pgtype.Int4{Int32: *e.BillingMonth, Valid: true}
	}
	return sqlc.CreateRecurringExpenseParams{
		UserID:       e.UserID,
		CategoryID:   toPGUUID(e.CategoryID),
		Name:         e.Name,
		Description:  toPGText(e.Description),
		Amount:       toPGNumeric(e.Amount),
		Frequency:    e.Frequency,
		DueDay:       dueDay,
		BillingMonth: billingMonth,
		Status:       e.Status,
		StartDate:    toPGDate(e.StartDate),
		EndDate:      toPGDatePtr(e.EndDate),
	}
}

func fromSQLCRecurringExpense(s sqlc.RecurringExpense) entity.RecurringExpense {
	var dueDay *int32
	if s.DueDay.Valid {
		dueDay = &s.DueDay.Int32
	}
	var billingMonth *int32
	if s.BillingMonth.Valid {
		billingMonth = &s.BillingMonth.Int32
	}
	return entity.RecurringExpense{
		ID:           s.ID,
		UserID:       s.UserID,
		CategoryID:   fromPGUUID(s.CategoryID),
		Name:         s.Name,
		Description:  fromPGText(s.Description),
		Amount:       fromPGNumeric(s.Amount),
		Frequency:    s.Frequency,
		DueDay:       dueDay,
		BillingMonth: billingMonth,
		Status:       s.Status,
		StartDate:    fromPGDate(s.StartDate),
		EndDate:      fromPGDatePtr(s.EndDate),
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

func (r *RecurringExpenseRepository) Create(ctx context.Context, e *entity.RecurringExpense) error {
	created, err := r.q.CreateRecurringExpense(ctx, toSQLCRecurringExpenseParams(e))
	if err != nil {
		return err
	}
	*e = fromSQLCRecurringExpense(created)
	return nil
}

func (r *RecurringExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RecurringExpense, error) {
	s, err := r.q.GetRecurringExpenseByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	e := fromSQLCRecurringExpense(s)
	return &e, nil
}

func (r *RecurringExpenseRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.RecurringExpense, error) {
	rows, err := r.q.ListRecurringExpensesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.RecurringExpense, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCRecurringExpense(row)
	}
	return result, nil
}

func (r *RecurringExpenseRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.RecurringExpense, error) {
	rows, err := r.q.ListActiveRecurringExpensesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.RecurringExpense, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCRecurringExpense(row)
	}
	return result, nil
}

func (r *RecurringExpenseRepository) Update(ctx context.Context, e *entity.RecurringExpense) error {
	dueDay := pgtype.Int4{Valid: false}
	if e.DueDay != nil {
		dueDay = pgtype.Int4{Int32: *e.DueDay, Valid: true}
	}
	billingMonth := pgtype.Int4{Valid: false}
	if e.BillingMonth != nil {
		billingMonth = pgtype.Int4{Int32: *e.BillingMonth, Valid: true}
	}
	updated, err := r.q.UpdateRecurringExpense(ctx, sqlc.UpdateRecurringExpenseParams{
		ID:           e.ID,
		Name:         e.Name,
		Description:  toPGText(e.Description),
		Amount:       toPGNumeric(e.Amount),
		Frequency:    e.Frequency,
		DueDay:       dueDay,
		BillingMonth: billingMonth,
		Status:       e.Status,
		EndDate:      toPGDatePtr(e.EndDate),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	*e = fromSQLCRecurringExpense(updated)
	return nil
}

func (r *RecurringExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteRecurringExpense(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

type ExpenseMonthlyStatusRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewExpenseMonthlyStatusRepository(db *PostgresDB) *ExpenseMonthlyStatusRepository {
	return &ExpenseMonthlyStatusRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCExpenseMonthlyStatusParams(s *entity.ExpenseMonthlyStatus) sqlc.CreateExpenseMonthlyStatusParams {
	amountPaid := pgtype.Numeric{Valid: false}
	if s.AmountPaid != nil {
		amountPaid = toPGNumeric(*s.AmountPaid)
	}
	return sqlc.CreateExpenseMonthlyStatusParams{
		ExpenseID:  s.ExpenseID,
		Month:      toPGDate(s.Month),
		AmountPaid: amountPaid,
		Notes:      toPGText(s.Notes),
	}
}

func fromSQLCExpenseMonthlyStatus(s sqlc.ExpenseMonthlyStatus) entity.ExpenseMonthlyStatus {
	var amountPaid *float64
	if s.AmountPaid.Valid {
		v := fromPGNumeric(s.AmountPaid)
		amountPaid = &v
	}
	return entity.ExpenseMonthlyStatus{
		ID:        s.ID,
		ExpenseID: s.ExpenseID,
		Month:     fromPGDate(s.Month),
		Paid:      s.Paid,
		PaidAt:    fromPGTimestamptz(s.PaidAt),
		AmountPaid: amountPaid,
		Notes:     fromPGText(s.Notes),
		CreatedAt: s.CreatedAt,
	}
}

func (r *ExpenseMonthlyStatusRepository) Create(ctx context.Context, s *entity.ExpenseMonthlyStatus) error {
	created, err := r.q.CreateExpenseMonthlyStatus(ctx, toSQLCExpenseMonthlyStatusParams(s))
	if err != nil {
		return err
	}
	*s = fromSQLCExpenseMonthlyStatus(created)
	return nil
}

func (r *ExpenseMonthlyStatusRepository) GetByExpenseIDAndMonth(ctx context.Context, expenseID uuid.UUID, month string) (*entity.ExpenseMonthlyStatus, error) {
	var t pgtype.Date
	if err := t.Scan(month); err != nil {
		return nil, err
	}
	s, err := r.q.GetExpenseMonthlyStatus(ctx, sqlc.GetExpenseMonthlyStatusParams{
		ExpenseID: expenseID,
		Month:     t,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	d := fromSQLCExpenseMonthlyStatus(s)
	return &d, nil
}

func (r *ExpenseMonthlyStatusRepository) ListByExpenseID(ctx context.Context, expenseID uuid.UUID) ([]entity.ExpenseMonthlyStatus, error) {
	rows, err := r.q.ListExpenseMonthlyStatusByExpenseID(ctx, expenseID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.ExpenseMonthlyStatus, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCExpenseMonthlyStatus(row)
	}
	return result, nil
}

func (r *ExpenseMonthlyStatusRepository) Upsert(ctx context.Context, s *entity.ExpenseMonthlyStatus) error {
	amountPaid := pgtype.Numeric{Valid: false}
	if s.AmountPaid != nil {
		amountPaid = toPGNumeric(*s.AmountPaid)
	}
	updated, err := r.q.UpsertExpenseMonthlyStatus(ctx, sqlc.UpsertExpenseMonthlyStatusParams{
		ExpenseID:  s.ExpenseID,
		Month:      toPGDate(s.Month),
		Paid:       s.Paid,
		AmountPaid: amountPaid,
		Notes:      toPGText(s.Notes),
	})
	if err != nil {
		return err
	}
	*s = fromSQLCExpenseMonthlyStatus(updated)
	return nil
}

func (r *ExpenseMonthlyStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteExpenseMonthlyStatus(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
