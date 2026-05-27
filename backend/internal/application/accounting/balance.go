package accounting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type GetMonthlyBalanceUseCase struct {
	summaryRepo repository.MonthlySummaryRepository
	txRepo      repository.TransactionRepository
}

func NewGetMonthlyBalanceUseCase(summaryRepo repository.MonthlySummaryRepository, txRepo repository.TransactionRepository) *GetMonthlyBalanceUseCase {
	return &GetMonthlyBalanceUseCase{summaryRepo: summaryRepo, txRepo: txRepo}
}

func (uc *GetMonthlyBalanceUseCase) Execute(ctx context.Context, userID uuid.UUID, year, month int) (*dto.MonthlyBalanceOutput, error) {
	monthStr := fmt.Sprintf("%04d-%02d-01", year, month)

	summary, err := uc.summaryRepo.GetByUserIDAndMonth(ctx, userID, monthStr)
	if err != nil {
		return nil, fmt.Errorf("get monthly summary: %w", err)
	}

	return summaryToOutput(summary), nil
}

type GetCashFlowUseCase struct {
	txRepo repository.TransactionRepository
}

func NewGetCashFlowUseCase(txRepo repository.TransactionRepository) *GetCashFlowUseCase {
	return &GetCashFlowUseCase{txRepo: txRepo}
}

func (uc *GetCashFlowUseCase) Execute(ctx context.Context, userID uuid.UUID, startDate, endDate string) (*dto.CashFlowOutput, error) {
	txs, err := uc.txRepo.ListByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("list transactions by date range: %w", err)
	}

	var totalIn, totalOut float64
	for _, t := range txs {
		if t.Type == "income" {
			totalIn += t.Amount
		} else {
			totalOut += t.Amount
		}
	}

	return &dto.CashFlowOutput{
		Period:   fmt.Sprintf("%s_to_%s", startDate, endDate),
		Entries:  []dto.CashFlowEntry{},
		TotalIn:  totalIn,
		TotalOut: totalOut,
		Balance:  totalIn - totalOut,
	}, nil
}

type GetProjectionsUseCase struct {
	summaryRepo   repository.MonthlySummaryRepository
	expenseRepo   repository.RecurringExpenseRepository
	debtRepo      repository.DebtRepository
}

func NewGetProjectionsUseCase(summaryRepo repository.MonthlySummaryRepository, expenseRepo repository.RecurringExpenseRepository, debtRepo repository.DebtRepository) *GetProjectionsUseCase {
	return &GetProjectionsUseCase{summaryRepo: summaryRepo, expenseRepo: expenseRepo, debtRepo: debtRepo}
}

func (uc *GetProjectionsUseCase) Execute(ctx context.Context, userID uuid.UUID, months int) ([]dto.ProjectionOutput, error) {
	projections := make([]dto.ProjectionOutput, months)

	expenses, err := uc.expenseRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list active expenses: %w", err)
	}

	_, err = uc.debtRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list active debts: %w", err)
	}

	lastSummary, err := uc.summaryRepo.ListByYear(ctx, userID, "2020-01-01", "2030-01-01")
	if err == nil && len(lastSummary) > 0 {
		_ = lastSummary
	}

	var monthlyExpenses float64
	for _, e := range expenses {
		if e.Frequency == "monthly" {
			monthlyExpenses += e.Amount
		} else if e.Frequency == "yearly" {
			monthlyExpenses += e.Amount / 12
		}
	}

	for i := 0; i < months; i++ {
		projections[i] = dto.ProjectionOutput{
			Month:              fmt.Sprintf("month_%d", i+1),
			ProjectedIncome:    0,
			ProjectedExpenses:  monthlyExpenses,
			ProjectedDebts:     0,
			ProjectedBalance:   -monthlyExpenses,
		}
	}

	return projections, nil
}

func summaryToOutput(s *entity.MonthlySummary) *dto.MonthlyBalanceOutput {
	return &dto.MonthlyBalanceOutput{
		Month:             s.Month.Format("2006-01-02"),
		TotalIncome:       s.TotalIncome,
		TotalExpenses:     s.TotalExpenses,
		TotalDebtPayments: s.TotalDebtPayments,
		TotalObligations:  s.TotalObligations,
		NetBalance:        s.NetBalance,
		ProjectedIncome:   s.ProjectedIncome,
		IncomeBreakdown:   s.IncomeBreakdown,
		ExpenseBreakdown:  s.ExpenseBreakdown,
	}
}
