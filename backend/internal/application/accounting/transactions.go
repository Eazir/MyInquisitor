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

type RecordTransactionUseCase struct {
	txRepo repository.TransactionRepository
}

func NewRecordTransactionUseCase(txRepo repository.TransactionRepository) *RecordTransactionUseCase {
	return &RecordTransactionUseCase{txRepo: txRepo}
}

func (uc *RecordTransactionUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.CreateTransactionInput) (*dto.TransactionOutput, error) {
	refDate, err := time.Parse("2006-01-02", input.ReferenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid reference_date: %w", err)
	}

	tx := &entity.Transaction{
		UserID:             userID,
		CategoryID:         input.CategoryID,
		Type:               input.Type,
		Amount:             input.Amount,
		Description:        input.Description,
		Source:             input.Source,
		ReferenceDate:      refDate,
		RecurringExpenseID: input.RecurringExpenseID,
	}

	if err := uc.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	return transactionToOutput(tx), nil
}

type ListTransactionsUseCase struct {
	txRepo repository.TransactionRepository
}

func NewListTransactionsUseCase(txRepo repository.TransactionRepository) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{txRepo: txRepo}
}

func (uc *ListTransactionsUseCase) Execute(ctx context.Context, userID uuid.UUID, year, month int) ([]dto.TransactionOutput, error) {
	start := fmt.Sprintf("%04d-%02d-01", year, month)
	end := fmt.Sprintf("%04d-%02d-01", year+1, 1)
	if month < 12 {
		end = fmt.Sprintf("%04d-%02d-01", year, month+1)
	}

	txs, err := uc.txRepo.ListByUserIDAndMonth(ctx, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}

	output := make([]dto.TransactionOutput, len(txs))
	for i, t := range txs {
		output[i] = *transactionToOutput(&t)
	}

	return output, nil
}

func transactionToOutput(t *entity.Transaction) *dto.TransactionOutput {
	return &dto.TransactionOutput{
		ID:                 t.ID,
		UserID:             t.UserID,
		CategoryID:         t.CategoryID,
		Type:               t.Type,
		Amount:             t.Amount,
		Description:        t.Description,
		Source:             t.Source,
		ReferenceDate:      t.ReferenceDate.Format("2006-01-02"),
		IsRecurring:        t.IsRecurring,
		RecurringExpenseID: t.RecurringExpenseID,
		CreatedAt:          t.CreatedAt,
	}
}
