package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type CreateUseCase struct {
	expenseRepo repository.RecurringExpenseRepository
}

func NewCreateUseCase(expenseRepo repository.RecurringExpenseRepository) *CreateUseCase {
	return &CreateUseCase{expenseRepo: expenseRepo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.CreateExpenseInput) (*dto.ExpenseOutput, error) {
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate *time.Time
	if input.EndDate != nil {
		t, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		endDate = &t
	}

	if input.Frequency == "once" && endDate == nil {
		endDate = &startDate
	}

	expense := &entity.RecurringExpense{
		UserID:       userID,
		CategoryID:   input.CategoryID,
		Name:         input.Name,
		Description:  input.Description,
		Amount:       input.Amount,
		Frequency:    input.Frequency,
		DueDay:       input.DueDay,
		BillingMonth: input.BillingMonth,
		Status:       "active",
		StartDate:    startDate,
		EndDate:      endDate,
	}

	if err := uc.expenseRepo.Create(ctx, expense); err != nil {
		return nil, fmt.Errorf("create expense: %w", err)
	}

	return expenseToOutput(expense), nil
}

func expenseToOutput(e *entity.RecurringExpense) *dto.ExpenseOutput {
	var endDate *string
	if e.EndDate != nil {
		s := e.EndDate.Format("2006-01-02")
		endDate = &s
	}

	return &dto.ExpenseOutput{
		ID:           e.ID,
		UserID:       e.UserID,
		CategoryID:   e.CategoryID,
		Name:         e.Name,
		Description:  e.Description,
		Amount:       e.Amount,
		Frequency:    e.Frequency,
		DueDay:       e.DueDay,
		BillingMonth: e.BillingMonth,
		Status:       e.Status,
		StartDate:    e.StartDate.Format("2006-01-02"),
		EndDate:      endDate,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func monthlyStatusToOutput(m *entity.ExpenseMonthlyStatus) *dto.ExpenseMonthlyStatusOutput {
	return &dto.ExpenseMonthlyStatusOutput{
		ID:        m.ID,
		ExpenseID: m.ExpenseID,
		Month:     m.Month.Format("2006-01-02"),
		Paid:      m.Paid,
		PaidAt:    m.PaidAt,
		AmountPaid: m.AmountPaid,
		Notes:     m.Notes,
	}
}
