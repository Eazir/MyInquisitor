package debt

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
)

var ErrDebtNotEditable = errors.New("debt cannot be edited in its current status")

func generateInstallments(
	ctx context.Context,
	monthlyRepo repository.DebtMonthlyStatusRepository,
	debtID uuid.UUID,
	remainingAmount float64,
	interestRate float64,
	totalInstallments int32,
	startNum int32,
	startMonth time.Time,
) error {
	remainingCount := totalInstallments - startNum + 1
	if remainingCount <= 0 {
		return nil
	}

	principalPerInstallment := remainingAmount / float64(remainingCount)
	monthlyRate := interestRate / 100.0

	for j := int32(0); j < remainingCount; j++ {
		num := startNum + j
		remainingBefore := remainingAmount - principalPerInstallment*float64(j)
		interestAmt := math.Round(remainingBefore*monthlyRate*100) / 100
		amountDue := math.Round((principalPerInstallment+interestAmt)*100) / 100

		monthDate := startMonth.AddDate(0, int(j), 0)

		monthly := &entity.DebtMonthlyStatus{
			DebtID:            debtID,
			Month:             monthDate,
			InstallmentNum:    num,
			TotalInstallments: totalInstallments,
			AmountDue:         amountDue,
			InterestAmount:    interestAmt,
			PrincipalAmount:   principalPerInstallment,
		}

		if err := monthlyRepo.Create(ctx, monthly); err != nil {
			return fmt.Errorf("create monthly status for installment %d: %w", num, err)
		}
	}
	return nil
}
