package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)

type TransactionRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewTransactionRepository(db *PostgresDB) *TransactionRepository {
	return &TransactionRepository{db: db, q: sqlc.New(db.Pool)}
}

func toSQLCTransactionParams(t *entity.Transaction) sqlc.CreateTransactionParams {
	return sqlc.CreateTransactionParams{
		UserID:             t.UserID,
		CategoryID:         toPGUUID(t.CategoryID),
		Type:               t.Type,
		Amount:             toPGNumeric(t.Amount),
		Description:        toPGText(t.Description),
		Source:             toPGText(t.Source),
		ReferenceDate:      toPGDate(t.ReferenceDate),
		IsRecurring:        t.IsRecurring,
		RecurringExpenseID: toPGUUID(t.RecurringExpenseID),
	}
}

func fromSQLCTransaction(s sqlc.Transaction) entity.Transaction {
	return entity.Transaction{
		ID:                 s.ID,
		UserID:             s.UserID,
		CategoryID:         fromPGUUID(s.CategoryID),
		Type:               s.Type,
		Amount:             fromPGNumeric(s.Amount),
		Description:        fromPGText(s.Description),
		Source:             fromPGText(s.Source),
		ReferenceDate:      fromPGDate(s.ReferenceDate),
		IsRecurring:        s.IsRecurring,
		RecurringExpenseID: fromPGUUID(s.RecurringExpenseID),
		CreatedAt:          s.CreatedAt,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, t *entity.Transaction) error {
	created, err := r.q.CreateTransaction(ctx, toSQLCTransactionParams(t))
	if err != nil {
		return err
	}
	*t = fromSQLCTransaction(created)
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	s, err := r.q.GetTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	t := fromSQLCTransaction(s)
	return &t, nil
}

func (r *TransactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Transaction, error) {
	rows, err := r.q.ListTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Transaction, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCTransaction(row)
	}
	return result, nil
}

func (r *TransactionRepository) ListByUserIDAndMonth(ctx context.Context, userID uuid.UUID, startMonth, endMonth string) ([]entity.Transaction, error) {
	sm, em := pgtype.Date{}, pgtype.Date{}
	if err := sm.Scan(startMonth); err != nil {
		return nil, err
	}
	if err := em.Scan(endMonth); err != nil {
		return nil, err
	}
	rows, err := r.q.ListTransactionsByUserIDAndMonth(ctx, sqlc.ListTransactionsByUserIDAndMonthParams{
		UserID:        userID,
		ReferenceDate: sm,
		ReferenceDate_2: em,
	})
	if err != nil {
		return nil, err
	}
	result := make([]entity.Transaction, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCTransaction(row)
	}
	return result, nil
}

func (r *TransactionRepository) ListByDateRange(ctx context.Context, userID uuid.UUID, start, end string) ([]entity.Transaction, error) {
	sd, ed := pgtype.Date{}, pgtype.Date{}
	if err := sd.Scan(start); err != nil {
		return nil, err
	}
	if err := ed.Scan(end); err != nil {
		return nil, err
	}
	rows, err := r.q.ListTransactionsByDateRange(ctx, sqlc.ListTransactionsByDateRangeParams{
		UserID:          userID,
		ReferenceDate:   sd,
		ReferenceDate_2: ed,
	})
	if err != nil {
		return nil, err
	}
	result := make([]entity.Transaction, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCTransaction(row)
	}
	return result, nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *entity.Transaction) error {
	updated, err := r.q.UpdateTransaction(ctx, sqlc.UpdateTransactionParams{
		ID:            t.ID,
		UserID:        t.UserID,
		CategoryID:    toPGUUID(t.CategoryID),
		Amount:        toPGNumeric(t.Amount),
		Description:   toPGText(t.Description),
		Source:        toPGText(t.Source),
		ReferenceDate: toPGDate(t.ReferenceDate),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	*t = fromSQLCTransaction(updated)
	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	err := r.q.DeleteTransaction(ctx, sqlc.DeleteTransactionParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

type MonthlySummaryRepository struct {
	db *PostgresDB
	q  *sqlc.Queries
}

func NewMonthlySummaryRepository(db *PostgresDB) *MonthlySummaryRepository {
	return &MonthlySummaryRepository{db: db, q: sqlc.New(db.Pool)}
}

func marshalBreakdown(m map[string]float64) []byte {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	return b
}

func unmarshalBreakdown(b []byte) map[string]float64 {
	if b == nil {
		return nil
	}
	var m map[string]float64
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

func fromSQLCMonthlySummary(s sqlc.MonthlySummary) entity.MonthlySummary {
	var projectedIncome *float64
	if s.ProjectedIncome.Valid {
		v := fromPGNumeric(s.ProjectedIncome)
		projectedIncome = &v
	}
	return entity.MonthlySummary{
		ID:                s.ID,
		UserID:            s.UserID,
		Month:             fromPGDate(s.Month),
		TotalIncome:       fromPGNumeric(s.TotalIncome),
		IncomeBreakdown:   unmarshalBreakdown(s.IncomeBreakdown),
		TotalExpenses:     fromPGNumeric(s.TotalExpenses),
		ExpenseBreakdown:  unmarshalBreakdown(s.ExpenseBreakdown),
		TotalDebtPayments: fromPGNumeric(s.TotalDebtPayments),
		DebtBreakdown:     unmarshalBreakdown(s.DebtBreakdown),
		TotalObligations:  fromPGNumeric(s.TotalObligations),
		NetBalance:        fromPGNumeric(s.NetBalance),
		ProjectedIncome:   projectedIncome,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

func (r *MonthlySummaryRepository) Upsert(ctx context.Context, s *entity.MonthlySummary) error {
	projectedIncome := pgtype.Numeric{Valid: false}
	if s.ProjectedIncome != nil {
		projectedIncome = toPGNumeric(*s.ProjectedIncome)
	}
	updated, err := r.q.UpsertMonthlySummary(ctx, sqlc.UpsertMonthlySummaryParams{
		UserID:            s.UserID,
		Month:             toPGDate(s.Month),
		TotalIncome:       toPGNumeric(s.TotalIncome),
		IncomeBreakdown:   marshalBreakdown(s.IncomeBreakdown),
		TotalExpenses:     toPGNumeric(s.TotalExpenses),
		ExpenseBreakdown:  marshalBreakdown(s.ExpenseBreakdown),
		TotalDebtPayments: toPGNumeric(s.TotalDebtPayments),
		DebtBreakdown:     marshalBreakdown(s.DebtBreakdown),
		TotalObligations:  toPGNumeric(s.TotalObligations),
		NetBalance:        toPGNumeric(s.NetBalance),
		ProjectedIncome:   projectedIncome,
	})
	if err != nil {
		return err
	}
	*s = fromSQLCMonthlySummary(updated)
	return nil
}

func (r *MonthlySummaryRepository) GetByUserIDAndMonth(ctx context.Context, userID uuid.UUID, month string) (*entity.MonthlySummary, error) {
	var m pgtype.Date
	if err := m.Scan(month); err != nil {
		return nil, err
	}
	s, err := r.q.GetMonthlySummary(ctx, sqlc.GetMonthlySummaryParams{
		UserID: userID,
		Month:  m,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	d := fromSQLCMonthlySummary(s)
	return &d, nil
}

func (r *MonthlySummaryRepository) ListByYear(ctx context.Context, userID uuid.UUID, startYear, endYear string) ([]entity.MonthlySummary, error) {
	sy, ey := pgtype.Date{}, pgtype.Date{}
	if err := sy.Scan(startYear); err != nil {
		return nil, err
	}
	if err := ey.Scan(endYear); err != nil {
		return nil, err
	}
	rows, err := r.q.ListMonthlySummariesByYear(ctx, sqlc.ListMonthlySummariesByYearParams{
		UserID:  userID,
		Month:   sy,
		Month_2: ey,
	})
	if err != nil {
		return nil, err
	}
	result := make([]entity.MonthlySummary, len(rows))
	for i, row := range rows {
		result[i] = fromSQLCMonthlySummary(row)
	}
	return result, nil
}

func (r *MonthlySummaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteMonthlySummary(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
