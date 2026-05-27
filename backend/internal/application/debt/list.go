package debt

import (
	"context"
	"fmt"

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
	debtRepo repository.DebtRepository
}

func NewUpdateUseCase(debtRepo repository.DebtRepository) *UpdateUseCase {
	return &UpdateUseCase{debtRepo: debtRepo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id uuid.UUID, input dto.UpdateDebtInput) (*dto.DebtOutput, error) {
	debt, err := uc.debtRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get debt: %w", err)
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
