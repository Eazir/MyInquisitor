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

type ListUseCase struct {
	debtRepo repository.DebtRepository
}

func NewListUseCase(debtRepo repository.DebtRepository) *ListUseCase {
	return &ListUseCase{debtRepo: debtRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]dto.DebtOutput, error) {
	var domainDebts []entity.Debt

	var err error
	if activeOnly {
		domainDebts, err = uc.debtRepo.ListActiveByUserID(ctx, userID)
	} else {
		domainDebts, err = uc.debtRepo.ListByUserID(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("list debts: %w", err)
	}

	output := make([]dto.DebtOutput, len(domainDebts))
	for i, d := range domainDebts {
		out := debtToOutput(&d)
		output[i] = *out
	}

	return output, nil
}

type GetByIDUseCase struct {
	debtRepo repository.DebtRepository
}

func NewGetByIDUseCase(debtRepo repository.DebtRepository) *GetByIDUseCase {
	return &GetByIDUseCase{debtRepo: debtRepo}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id uuid.UUID) (*dto.DebtOutput, error) {
	debt, err := uc.debtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get debt: %w", err)
	}

	return debtToOutput(debt), nil
}

type UpdateUseCase struct {
	debtRepo    repository.DebtRepository
	monthlyRepo repository.DebtMonthlyStatusRepository
}

func NewUpdateUseCase(debtRepo repository.DebtRepository, monthlyRepo repository.DebtMonthlyStatusRepository) *UpdateUseCase {
	return &UpdateUseCase{debtRepo: debtRepo, monthlyRepo: monthlyRepo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id uuid.UUID, input dto.UpdateDebtInput) (*dto.DebtOutput, error) {
	debt, err := uc.debtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get debt: %w", err)
	}

	if debt.Status == "paused" || debt.Status == "paid" {
		if input.Status == nil {
			return nil, ErrDebtNotEditable
		}
		input.Name = nil
		input.Description = nil
		input.TotalAmount = nil
		input.RemainingAmount = nil
		input.InterestRate = nil
		input.TotalInstallments = nil
		input.StartDate = nil
		input.DueDay = nil
	}

	financialChanged := input.TotalAmount != nil || input.InterestRate != nil || input.TotalInstallments != nil
	reactivating := debt.Status == "paused" && input.Status != nil && *input.Status == "active"

	// If total_amount changes but remaining_amount is not explicit, apply the same delta
	adjustedRemainingAmount := debt.RemainingAmount
	if input.RemainingAmount != nil {
		adjustedRemainingAmount = *input.RemainingAmount
	} else if input.TotalAmount != nil {
		delta := *input.TotalAmount - debt.TotalAmount
		adjustedRemainingAmount = debt.RemainingAmount + delta
		if adjustedRemainingAmount < 0 {
			adjustedRemainingAmount = 0
		}
	}

	if financialChanged || reactivating {
		allStatuses, err := uc.monthlyRepo.ListByDebtID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("list monthly statuses: %w", err)
		}

		var paidCount int32
		for _, s := range allStatuses {
			if s.Paid {
				paidCount++
			}
		}

		newTotalInstallments := debt.TotalInstallments
		if input.TotalInstallments != nil {
			newTotalInstallments = *input.TotalInstallments
		}

		remainingCount := newTotalInstallments - paidCount
		if remainingCount > 0 && paidCount < newTotalInstallments {
			newInterestRate := debt.InterestRate
			if input.InterestRate != nil {
				newInterestRate = *input.InterestRate
			}

			startNum := paidCount + 1

			var firstUnpaidMonth time.Time
			for _, s := range allStatuses {
				if !s.Paid && s.InstallmentNum == startNum {
					firstUnpaidMonth = s.Month
					break
				}
			}

			for _, s := range allStatuses {
				if !s.Paid && s.InstallmentNum <= newTotalInstallments {
					if err := uc.monthlyRepo.Delete(ctx, s.ID); err != nil {
						return nil, fmt.Errorf("delete unpaid installment %d: %w", s.InstallmentNum, err)
					}
				}
			}

			startMonth := firstUnpaidMonth
			if startMonth.IsZero() {
				location := debt.StartDate.Location()
				if location == nil {
					location = time.UTC
				}
				startMonth = time.Date(debt.StartDate.Year(), debt.StartDate.Month()+time.Month(paidCount), 1, 0, 0, 0, 0, location)
			}

			shift := 0
			if reactivating && !firstUnpaidMonth.IsZero() {
				now := time.Now()
				shift = int(now.Year()-firstUnpaidMonth.Year())*12 + int(now.Month()-firstUnpaidMonth.Month())
			}
			startMonth = startMonth.AddDate(0, shift, 0)

			if err := generateInstallments(ctx, uc.monthlyRepo, id, adjustedRemainingAmount, newInterestRate, newTotalInstallments, startNum, startMonth); err != nil {
				return nil, err
			}
		}
	}

	if input.Name != nil {
		debt.Name = *input.Name
	}
	if input.Description != nil {
		debt.Description = input.Description
	}
	if input.TotalAmount != nil {
		debt.TotalAmount = *input.TotalAmount
	}
	if input.RemainingAmount != nil {
		debt.RemainingAmount = *input.RemainingAmount
	} else if input.TotalAmount != nil {
		debt.RemainingAmount = adjustedRemainingAmount
	}
	if input.InterestRate != nil {
		debt.InterestRate = *input.InterestRate
	}
	if input.TotalInstallments != nil {
		debt.TotalInstallments = *input.TotalInstallments
	}
	if input.Status != nil {
		debt.Status = *input.Status
	}
	if input.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *input.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date: %w", err)
		}
		debt.StartDate = parsed
	}
	if input.DueDay != nil {
		debt.DueDay = input.DueDay
	}

	if err := uc.debtRepo.Update(ctx, debt); err != nil {
		return nil, fmt.Errorf("update debt: %w", err)
	}

	return debtToOutput(debt), nil
}

type DeleteUseCase struct {
	debtRepo repository.DebtRepository
}

func NewDeleteUseCase(debtRepo repository.DebtRepository) *DeleteUseCase {
	return &DeleteUseCase{debtRepo: debtRepo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.debtRepo.Delete(ctx, id)
}
