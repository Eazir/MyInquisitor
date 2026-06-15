package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appAcc "github.com/myinquisitor/backend/internal/application/accounting"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type AccountingHandler struct {
	recordTxUC    *appAcc.RecordTransactionUseCase
	listTxUC      *appAcc.ListTransactionsUseCase
	balanceUC     *appAcc.GetMonthlyBalanceUseCase
	cashFlowUC    *appAcc.GetCashFlowUseCase
	projectionsUC *appAcc.GetProjectionsUseCase
	createCatUC   *appAcc.CreateCategoryUseCase
	listCatUC     *appAcc.ListCategoriesUseCase
	deleteCatUC   *appAcc.DeleteCategoryUseCase
}

func NewAccountingHandler(
	recordTxUC *appAcc.RecordTransactionUseCase,
	listTxUC *appAcc.ListTransactionsUseCase,
	balanceUC *appAcc.GetMonthlyBalanceUseCase,
	cashFlowUC *appAcc.GetCashFlowUseCase,
	projectionsUC *appAcc.GetProjectionsUseCase,
	createCatUC *appAcc.CreateCategoryUseCase,
	listCatUC *appAcc.ListCategoriesUseCase,
	deleteCatUC *appAcc.DeleteCategoryUseCase,
) *AccountingHandler {
	return &AccountingHandler{
		recordTxUC:    recordTxUC,
		listTxUC:      listTxUC,
		balanceUC:     balanceUC,
		cashFlowUC:    cashFlowUC,
		projectionsUC: projectionsUC,
		createCatUC:   createCatUC,
		listCatUC:     listCatUC,
		deleteCatUC:   deleteCatUC,
	}
}

func (h *AccountingHandler) RecordTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input dto.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.recordTxUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *AccountingHandler) ListTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	now := time.Now()
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))

	result, err := h.listTxUC.Execute(c.Request.Context(), userID, year, month)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AccountingHandler) MonthlyBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))

	result, err := h.balanceUC.Execute(c.Request.Context(), userID, year, month)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AccountingHandler) CashFlow(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	start := c.DefaultQuery("start", "2026-01-01")
	end := c.DefaultQuery("end", "2026-12-31")

	result, err := h.cashFlowUC.Execute(c.Request.Context(), userID, start, end)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AccountingHandler) Projections(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	months, _ := strconv.Atoi(c.DefaultQuery("months", "6"))
	if months < 1 || months > 24 {
		months = 6
	}

	historyMonths, _ := strconv.Atoi(c.DefaultQuery("history_months", "12"))
	if historyMonths < 1 || historyMonths > 36 {
		historyMonths = 12
	}

	result, err := h.projectionsUC.Execute(c.Request.Context(), userID, months, historyMonths)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AccountingHandler) CreateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input dto.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.createCatUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *AccountingHandler) ListCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	categoryType := c.Query("type")

	result, err := h.listCatUC.Execute(c.Request.Context(), userID, categoryType)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AccountingHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid category id")
		return
	}

	if err := h.deleteCatUC.Execute(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}
