package expense

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type ListUseCase struct {
	expenseRepo repository.RecurringExpenseRepository
}

func NewListUseCase(expenseRepo repository.RecurringExpenseRepository) *ListUseCase {
	return &ListUseCase{expenseRepo: expenseRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]dto.ExpenseOutput, error) {
	var domainExpenses []entity.RecurringExpense

	var err error
	if activeOnly {
		domainExpenses, err = uc.expenseRepo.ListActiveByUserID(ctx, userID)
	} else {
		domainExpenses, err = uc.expenseRepo.ListByUserID(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}

	output := make([]dto.ExpenseOutput, len(domainExpenses))
	for i, e := range domainExpenses {
		out := expenseToOutput(&e)
		output[i] = *out
	}

	return output, nil
}

type GetByIDUseCase struct {
	expenseRepo repository.RecurringExpenseRepository
}

func NewGetByIDUseCase(expenseRepo repository.RecurringExpenseRepository) *GetByIDUseCase {
	return &GetByIDUseCase{expenseRepo: expenseRepo}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id uuid.UUID) (*dto.ExpenseOutput, error) {
	expense, err := uc.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get expense: %w", err)
	}

	return expenseToOutput(expense), nil
}

type UpdateUseCase struct {
	expenseRepo repository.RecurringExpenseRepository
}

func NewUpdateUseCase(expenseRepo repository.RecurringExpenseRepository) *UpdateUseCase {
	return &UpdateUseCase{expenseRepo: expenseRepo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id uuid.UUID, input dto.UpdateExpenseInput) (*dto.ExpenseOutput, error) {
	expense, err := uc.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get expense: %w", err)
	}

	if input.Name != nil {
		expense.Name = *input.Name
	}
	if input.Description != nil {
		expense.Description = input.Description
	}
	if input.Amount != nil {
		expense.Amount = *input.Amount
	}
	if input.Frequency != nil {
		expense.Frequency = *input.Frequency
	}
	if input.DueDay != nil {
		expense.DueDay = input.DueDay
	}
	if input.BillingMonth != nil {
		expense.BillingMonth = input.BillingMonth
	}
	if input.Status != nil {
		expense.Status = *input.Status
	}

	if err := uc.expenseRepo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("update expense: %w", err)
	}

	return expenseToOutput(expense), nil
}

type DeleteUseCase struct {
	expenseRepo repository.RecurringExpenseRepository
}

func NewDeleteUseCase(expenseRepo repository.RecurringExpenseRepository) *DeleteUseCase {
	return &DeleteUseCase{expenseRepo: expenseRepo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.expenseRepo.Delete(ctx, id)
}

type GetMonthlyStatusUseCase struct {
	monthlyRepo repository.ExpenseMonthlyStatusRepository
}

func NewGetMonthlyStatusUseCase(monthlyRepo repository.ExpenseMonthlyStatusRepository) *GetMonthlyStatusUseCase {
	return &GetMonthlyStatusUseCase{monthlyRepo: monthlyRepo}
}

func (uc *GetMonthlyStatusUseCase) Execute(ctx context.Context, expenseID uuid.UUID, year, month int) (*dto.ExpenseMonthlyStatusOutput, error) {
	monthStr := fmt.Sprintf("%04d-%02d-01", year, month)
	status, err := uc.monthlyRepo.GetByExpenseIDAndMonth(ctx, expenseID, monthStr)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get expense monthly status: %w", err)
	}
	return monthlyStatusToOutput(status), nil
}

type TogglePaidUseCase struct {
	monthlyRepo repository.ExpenseMonthlyStatusRepository
	expenseRepo repository.RecurringExpenseRepository
	txRepo      repository.TransactionRepository
}

func NewTogglePaidUseCase(monthlyRepo repository.ExpenseMonthlyStatusRepository, expenseRepo repository.RecurringExpenseRepository, txRepo repository.TransactionRepository) *TogglePaidUseCase {
	return &TogglePaidUseCase{monthlyRepo: monthlyRepo, expenseRepo: expenseRepo, txRepo: txRepo}
}

func (uc *TogglePaidUseCase) Execute(ctx context.Context, expenseID uuid.UUID, year, month int, input dto.ToggleExpensePaidInput) (*dto.ExpenseMonthlyStatusOutput, error) {
	monthStr := fmt.Sprintf("%04d-%02d-01", year, month)

	status, err := uc.monthlyRepo.GetByExpenseIDAndMonth(ctx, expenseID, monthStr)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			monthDate, _ := time.Parse("2006-01-02", monthStr)
			status = &entity.ExpenseMonthlyStatus{
				ExpenseID: expenseID,
				Month:     monthDate,
				Paid:      false,
			}
		} else {
			return nil, fmt.Errorf("get expense monthly status: %w", err)
		}
	}

	status.Paid = !status.Paid
	if status.Paid {
		now := time.Now()
		status.PaidAt = &now
		if input.AmountPaid != nil {
			status.AmountPaid = input.AmountPaid
		}
		status.Notes = input.Notes
	} else {
		status.PaidAt = nil
		status.AmountPaid = nil
	}

	if err := uc.monthlyRepo.Upsert(ctx, status); err != nil {
		return nil, fmt.Errorf("toggle expense paid: %w", err)
	}

	if status.Paid {
		expense, err := uc.expenseRepo.GetByID(ctx, expenseID)
		if err != nil {
			return nil, fmt.Errorf("get expense: %w", err)
		}
		amount := expense.Amount
		if status.AmountPaid != nil {
			amount = *status.AmountPaid
		}
		desc := fmt.Sprintf("Gasto: %s - %s", expense.Name, monthStr)
		tx := &entity.Transaction{
			UserID:        expense.UserID,
			Type:          "expense",
			Amount:        amount,
			Description:   &desc,
			ReferenceDate: time.Now(),
		}
		if err := uc.txRepo.Create(ctx, tx); err != nil {
			return nil, fmt.Errorf("create transaction: %w", err)
		}

		if expense.Frequency == "once" {
			expense.Status = "cancelled"
			if err := uc.expenseRepo.Update(ctx, expense); err != nil {
				return nil, fmt.Errorf("cancel once expense: %w", err)
			}
		}
	}

	return monthlyStatusToOutput(status), nil
}
