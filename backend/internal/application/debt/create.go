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

type CreateUseCase struct {
	debtRepo      repository.DebtRepository
	monthlyRepo   repository.DebtMonthlyStatusRepository
}

func NewCreateUseCase(debtRepo repository.DebtRepository, monthlyRepo repository.DebtMonthlyStatusRepository) *CreateUseCase {
	return &CreateUseCase{debtRepo: debtRepo, monthlyRepo: monthlyRepo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.CreateDebtInput) (*dto.DebtOutput, error) {
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

	debt := &entity.Debt{
		UserID:             userID,
		CategoryID:         input.CategoryID,
		Name:               input.Name,
		Description:        input.Description,
		TotalAmount:        input.TotalAmount,
		RemainingAmount:    input.TotalAmount,
		InterestRate:       input.InterestRate,
		TotalInstallments:  input.TotalInstallments,
		CurrentInstallment: 0,
		Status:             "active",
		StartDate:          startDate,
		EndDate:            endDate,
		DueDay:             input.DueDay,
	}

	if err := uc.debtRepo.Create(ctx, debt); err != nil {
		return nil, fmt.Errorf("create debt: %w", err)
	}

	startMonth := time.Date(debt.StartDate.Year(), debt.StartDate.Month(), 1, 0, 0, 0, 0, debt.StartDate.Location())
	if err := generateInstallments(ctx, uc.monthlyRepo, debt.ID, debt.TotalAmount, debt.InterestRate, debt.TotalInstallments, 1, startMonth); err != nil {
		return nil, err
	}

	return debtToOutput(debt), nil
}

func debtToOutput(d *entity.Debt) *dto.DebtOutput {
	var endDate *string
	if d.EndDate != nil {
		s := d.EndDate.Format("2006-01-02")
		endDate = &s
	}

	return &dto.DebtOutput{
		ID:                 d.ID,
		UserID:             d.UserID,
		CategoryID:         d.CategoryID,
		Name:               d.Name,
		Description:        d.Description,
		TotalAmount:        d.TotalAmount,
		RemainingAmount:    d.RemainingAmount,
		InterestRate:       d.InterestRate,
		TotalInstallments:  d.TotalInstallments,
		CurrentInstallment: d.CurrentInstallment,
		Status:             d.Status,
		StartDate:          d.StartDate.Format("2006-01-02"),
		EndDate:            endDate,
		DueDay:             d.DueDay,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

func monthlyToOutput(m *entity.DebtMonthlyStatus) *dto.DebtMonthlyStatusOutput {
	return &dto.DebtMonthlyStatusOutput{
		ID:              m.ID,
		DebtID:          m.DebtID,
		Month:           m.Month.Format("2006-01-02"),
		InstallmentNum:  m.InstallmentNum,
		AmountDue:       m.AmountDue,
		InterestAmount:  m.InterestAmount,
		PrincipalAmount: m.PrincipalAmount,
		AmountPaid:      m.AmountPaid,
		Paid:            m.Paid,
		PaidAt:          m.PaidAt,
		Notes:           m.Notes,
	}
}
