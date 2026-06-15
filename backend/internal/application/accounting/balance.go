package accounting

import (
	"context"
	"fmt"
	"time"

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

	// Always compute from raw transactions for real-time accuracy
	start := monthStr
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	end := fmt.Sprintf("%04d-%02d-01", nextYear, nextMonth)

	txs, err := uc.txRepo.ListByUserIDAndMonth(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list transactions for balance: %w", err)
	}

	var totalIncome, totalExpenses float64
	for _, t := range txs {
		switch t.Type {
		case "income":
			totalIncome += t.Amount
		case "expense":
			totalExpenses += t.Amount
		}
	}

	netBalance := totalIncome - totalExpenses

	return &dto.MonthlyBalanceOutput{
		Month:            monthStr,
		TotalIncome:      totalIncome,
		TotalExpenses:    totalExpenses,
		NetBalance:       netBalance,
		TotalObligations: totalExpenses,
	}, nil
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
	txRepo          repository.TransactionRepository
	expenseRepo     repository.RecurringExpenseRepository
	debtMonthlyRepo repository.DebtMonthlyStatusRepository
}

func NewGetProjectionsUseCase(txRepo repository.TransactionRepository, expenseRepo repository.RecurringExpenseRepository, debtMonthlyRepo repository.DebtMonthlyStatusRepository) *GetProjectionsUseCase {
	return &GetProjectionsUseCase{txRepo: txRepo, expenseRepo: expenseRepo, debtMonthlyRepo: debtMonthlyRepo}
}

func (uc *GetProjectionsUseCase) Execute(ctx context.Context, userID uuid.UUID, months int, historyMonths int) ([]dto.ProjectionOutput, error) {
	now := time.Now()

	startProj := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endProj := startProj.AddDate(0, months, 0)

	// 1. Historical averages from transactions
	startHist := startProj.AddDate(0, -historyMonths, 0)

	histTx, err := uc.txRepo.ListByDateRange(ctx, userID, startHist.Format("2006-01-02"), startProj.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("list historical transactions: %w", err)
	}

	incomeByMonth := make(map[string]float64)
	expenseByMonth := make(map[string]float64)
	for _, tx := range histTx {
		monthKey := tx.ReferenceDate.Format("2006-01")
		switch tx.Type {
		case "income":
			incomeByMonth[monthKey] += tx.Amount
		case "expense":
			if tx.RecurringExpenseID == nil {
				expenseByMonth[monthKey] += tx.Amount
			}
		}
	}

	avgIncome := avgFromMap(incomeByMonth)
	avgExtra := avgFromMap(expenseByMonth)

	// 2. Fixed expenses from recurring expenses
	expenses, err := uc.expenseRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list active expenses: %w", err)
	}

	// Build itemized list per month
	fixedItemsByMonth := make(map[string][]dto.FixedExpenseItem)
	for _, e := range expenses {
		for i := 0; i < months; i++ {
			m := startProj.AddDate(0, i, 0)
			monthKey := m.Format("2006-01")
			var amount float64
			var apply bool
			switch e.Frequency {
			case "monthly":
				apply = true
				amount = e.Amount
			case "yearly":
				billMonth := e.StartDate.Month()
				if e.BillingMonth != nil {
					billMonth = time.Month(*e.BillingMonth)
				}
				apply = m.Month() == billMonth
				amount = e.Amount
			case "weekly":
				apply = true
				amount = e.Amount * 52 / 12
			case "biweekly":
				apply = true
				amount = e.Amount * 26 / 12
			case "once":
				apply = m.Year() == e.StartDate.Year() && m.Month() == e.StartDate.Month()
				amount = e.Amount
			}
			if apply {
				fixedItemsByMonth[monthKey] = append(fixedItemsByMonth[monthKey], dto.FixedExpenseItem{
					ID:        e.ID.String(),
					Name:      e.Name,
					Amount:    amount,
					Frequency: e.Frequency,
					DueDay:    e.DueDay,
				})
			}
		}
	}

	// 3. Debt payments within projection range
	unpaidDebts, err := uc.debtMonthlyRepo.ListUnpaidByUserIDAndDateRange(ctx, userID, startProj.Format("2006-01-02"), endProj.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("list unpaid debts: %w", err)
	}

	debtByMonth := make(map[string]float64)
	for _, d := range unpaidDebts {
		monthKey := d.Month.Format("2006-01")
		debtByMonth[monthKey] += d.AmountDue
	}

	// 4. Build projections
	projections := make([]dto.ProjectionOutput, months)
	monthNames := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	for i := 0; i < months; i++ {
		m := startProj.AddDate(0, i, 0)
		monthKey := m.Format("2006-01")
		monthLabel := fmt.Sprintf("%s %d", monthNames[m.Month()-1], m.Year())

		items := fixedItemsByMonth[monthKey]
		var fixed float64
		for _, item := range items {
			fixed += item.Amount
		}

		debtPmt := debtByMonth[monthKey]
		totalExp := fixed + avgExtra + debtPmt

		flist := make([]dto.FixedExpenseItem, len(items))
		copy(flist, items)

		projections[i] = dto.ProjectionOutput{
			Month:              monthKey,
			MonthLabel:         monthLabel,
			BaseIncome:         avgIncome,
			IncomeModifier:     0,
			ProjectedIncome:    avgIncome,
			FixedExpenses:      fixed,
			FixedExpensesList:  flist,
			ExtraBudgetaryAvg:  avgExtra,
			ExtraExpensesTotal: 0,
			ExtraExpensesList:  []dto.ExtraExpenseItem{},
			DebtPayments:       debtPmt,
			TotalExpenses:      totalExp,
			ProjectedBalance:   avgIncome - totalExp,
		}
	}

	return projections, nil
}

func avgFromMap(m map[string]float64) float64 {
	if len(m) == 0 {
		return 0
	}
	var total float64
	for _, v := range m {
		total += v
	}
	return total / float64(len(m))
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
