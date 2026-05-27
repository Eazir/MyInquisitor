package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appExpense "github.com/myinquisitor/backend/internal/application/expense"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type ExpenseHandler struct {
	createUC    *appExpense.CreateUseCase
	listUC      *appExpense.ListUseCase
	getByIDUC   *appExpense.GetByIDUseCase
	updateUC    *appExpense.UpdateUseCase
	deleteUC    *appExpense.DeleteUseCase
	togglePaidUC *appExpense.TogglePaidUseCase
}

func NewExpenseHandler(
	createUC *appExpense.CreateUseCase,
	listUC *appExpense.ListUseCase,
	getByIDUC *appExpense.GetByIDUseCase,
	updateUC *appExpense.UpdateUseCase,
	deleteUC *appExpense.DeleteUseCase,
	togglePaidUC *appExpense.TogglePaidUseCase,
) *ExpenseHandler {
	return &ExpenseHandler{
		createUC:     createUC,
		listUC:       listUC,
		getByIDUC:    getByIDUC,
		updateUC:     updateUC,
		deleteUC:     deleteUC,
		togglePaidUC: togglePaidUC,
	}
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input dto.CreateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.createUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *ExpenseHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	activeOnly := c.Query("status") == "active"

	result, err := h.listUC.Execute(c.Request.Context(), userID, activeOnly)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid expense id")
		return
	}

	result, err := h.getByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "expense not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid expense id")
		return
	}

	var input dto.UpdateExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.updateUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "expense not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid expense id")
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}

func (h *ExpenseHandler) TogglePaid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid expense id")
		return
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PARAM", "invalid year")
		return
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PARAM", "invalid month")
		return
	}

	var input dto.ToggleExpensePaidInput
	if err := c.ShouldBindJSON(&input); err == nil {
		_ = input
	}

	result, err := h.togglePaidUC.Execute(c.Request.Context(), id, year, month, input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}
