package accounting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type CreateCategoryUseCase struct {
	catRepo repository.CategoryRepository
}

func NewCreateCategoryUseCase(catRepo repository.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{catRepo: catRepo}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.CreateCategoryInput) (*dto.CategoryOutput, error) {
	cat := &entity.Category{
		UserID: userID,
		Name:   input.Name,
		Type:   input.Type,
		Icon:   input.Icon,
		Color:  input.Color,
	}

	if err := uc.catRepo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}

	return categoryToOutput(cat), nil
}

type ListCategoriesUseCase struct {
	catRepo repository.CategoryRepository
}

func NewListCategoriesUseCase(catRepo repository.CategoryRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{catRepo: catRepo}
}

func (uc *ListCategoriesUseCase) Execute(ctx context.Context, userID uuid.UUID, categoryType string) ([]dto.CategoryOutput, error) {
	var cats []entity.Category
	var err error

	if categoryType != "" {
		cats, err = uc.catRepo.ListByType(ctx, userID, categoryType)
	} else {
		cats, err = uc.catRepo.ListByUserID(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	output := make([]dto.CategoryOutput, len(cats))
	for i, c := range cats {
		output[i] = *categoryToOutput(&c)
	}

	return output, nil
}

type DeleteCategoryUseCase struct {
	catRepo repository.CategoryRepository
}

func NewDeleteCategoryUseCase(catRepo repository.CategoryRepository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{catRepo: catRepo}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	return uc.catRepo.Delete(ctx, id)
}

func categoryToOutput(c *entity.Category) *dto.CategoryOutput {
	return &dto.CategoryOutput{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Type:      c.Type,
		Icon:      c.Icon,
		Color:     c.Color,
		CreatedAt: c.CreatedAt,
	}
}
