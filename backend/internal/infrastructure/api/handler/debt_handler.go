package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appDebt "github.com/myinquisitor/backend/internal/application/debt"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type DebtHandler struct {
	createUC         *appDebt.CreateUseCase
	listUC           *appDebt.ListUseCase
	getByIDUC        *appDebt.GetByIDUseCase
	updateUC         *appDebt.UpdateUseCase
	deleteUC         *appDebt.DeleteUseCase
	markPaidUC       *appDebt.MarkPaidUseCase
	getMonthlyStatus *appDebt.GetMonthlyStatusUseCase
}

func NewDebtHandler(
	createUC *appDebt.CreateUseCase,
	listUC *appDebt.ListUseCase,
	getByIDUC *appDebt.GetByIDUseCase,
	updateUC *appDebt.UpdateUseCase,
	deleteUC *appDebt.DeleteUseCase,
	markPaidUC *appDebt.MarkPaidUseCase,
	getMonthlyStatus *appDebt.GetMonthlyStatusUseCase,
) *DebtHandler {
	return &DebtHandler{
		createUC:         createUC,
		listUC:           listUC,
		getByIDUC:        getByIDUC,
		updateUC:         updateUC,
		deleteUC:         deleteUC,
		markPaidUC:       markPaidUC,
		getMonthlyStatus: getMonthlyStatus,
	}
}

func (h *DebtHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var input dto.CreateDebtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.createUC.Execute(c.Request.Context(), userID, input)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *DebtHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	activeOnly := c.Query("status") == "active"

	result, err := h.listUC.Execute(c.Request.Context(), userID, activeOnly)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *DebtHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid debt id")
		return
	}

	result, err := h.getByIDUC.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "debt not found")
			return
		}
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *DebtHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid debt id")
		return
	}

	var input dto.UpdateDebtInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.updateUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "debt not found")
			return
		}
		if errors.Is(err, appDebt.ErrDebtNotEditable) {
			response.Error(c, http.StatusBadRequest, "DEBT_NOT_EDITABLE", "debt is paused or paid")
			return
		}
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *DebtHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid debt id")
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), id); err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusNoContent, nil)
}

func (h *DebtHandler) MarkAsPaid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid debt id")
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

	var input dto.MarkDebtPaidInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.markPaidUC.Execute(c.Request.Context(), id, year, month, input)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *DebtHandler) GetMonthlyStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid debt id")
		return
	}

	result, err := h.getMonthlyStatus.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}
