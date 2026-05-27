package debt

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type MarkPaidUseCase struct {
	monthlyRepo repository.DebtMonthlyStatusRepository
	debtRepo    repository.DebtRepository
}

func NewMarkPaidUseCase(monthlyRepo repository.DebtMonthlyStatusRepository, debtRepo repository.DebtRepository) *MarkPaidUseCase {
	return &MarkPaidUseCase{monthlyRepo: monthlyRepo, debtRepo: debtRepo}
}

func (uc *MarkPaidUseCase) Execute(ctx context.Context, debtID uuid.UUID, year, month int, input dto.MarkDebtPaidInput) (*dto.DebtMonthlyStatusOutput, error) {
	monthStr := fmt.Sprintf("%04d-%02d-01", year, month)

	status, err := uc.monthlyRepo.GetByDebtIDAndMonth(ctx, debtID, monthStr)
	if err != nil {
		return nil, fmt.Errorf("get monthly status: %w", err)
	}

	now := time.Now()
	status.Paid = true
	status.PaidAt = &now
	status.AmountPaid = input.AmountPaid
	status.Notes = input.Notes

	if err := uc.monthlyRepo.Update(ctx, status); err != nil {
		return nil, fmt.Errorf("mark as paid: %w", err)
	}

	debt, err := uc.debtRepo.GetByID(ctx, debtID)
	if err != nil {
		return nil, fmt.Errorf("get debt: %w", err)
	}

	newRemaining := debt.RemainingAmount - input.AmountPaid
	if newRemaining < 0 {
		newRemaining = 0
	}
	debt.RemainingAmount = newRemaining
	debt.CurrentInstallment = status.InstallmentNum

	if newRemaining == 0 {
		debt.Status = "paid"
	}

	if err := uc.debtRepo.Update(ctx, debt); err != nil {
		return nil, fmt.Errorf("update debt after payment: %w", err)
	}

	_ = entity.DebtMonthlyStatus{}
	return monthlyToOutput(status), nil
}

type GetMonthlyStatusUseCase struct {
	monthlyRepo repository.DebtMonthlyStatusRepository
}

func NewGetMonthlyStatusUseCase(monthlyRepo repository.DebtMonthlyStatusRepository) *GetMonthlyStatusUseCase {
	return &GetMonthlyStatusUseCase{monthlyRepo: monthlyRepo}
}

func (uc *GetMonthlyStatusUseCase) Execute(ctx context.Context, debtID uuid.UUID) ([]dto.DebtMonthlyStatusOutput, error) {
	statuses, err := uc.monthlyRepo.ListByDebtID(ctx, debtID)
	if err != nil {
		return nil, fmt.Errorf("list monthly status: %w", err)
	}

	output := make([]dto.DebtMonthlyStatusOutput, len(statuses))
	for i, s := range statuses {
		out := monthlyToOutput(&s)
		output[i] = *out
	}

	return output, nil
}
