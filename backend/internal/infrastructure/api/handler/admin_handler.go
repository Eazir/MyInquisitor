package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appAdmin "github.com/myinquisitor/backend/internal/application/admin"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type AdminHandler struct {
	listUsersUC       *appAdmin.ListUsersUseCase
	createUserUC      *appAdmin.CreateUserUseCase
	updateUserUC      *appAdmin.UpdateUserUseCase
	deactivateUserUC  *appAdmin.DeactivateUserUseCase
	generateInviteUC  *appAdmin.GenerateInviteUseCase
}

func NewAdminHandler(
	listUsersUC *appAdmin.ListUsersUseCase,
	createUserUC *appAdmin.CreateUserUseCase,
	updateUserUC *appAdmin.UpdateUserUseCase,
	deactivateUserUC *appAdmin.DeactivateUserUseCase,
	generateInviteUC *appAdmin.GenerateInviteUseCase,
) *AdminHandler {
	return &AdminHandler{
		listUsersUC:       listUsersUC,
		createUserUC:      createUserUC,
		updateUserUC:      updateUserUC,
		deactivateUserUC:  deactivateUserUC,
		generateInviteUC:  generateInviteUC,
	}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.listUsersUC.Execute(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, users, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: int64(total),
	})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var input dto.AdminCreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.createUserUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	adminID, _ := c.Get("user_id")
	aID, _ := adminID.(uuid.UUID)

	var input dto.AdminUpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.updateUserUC.Execute(c.Request.Context(), id, aID, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AdminHandler) SetActive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	active, err := strconv.ParseBool(c.Param("active"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PARAM", "active must be true or false")
		return
	}

	result, err := h.deactivateUserUC.Execute(c.Request.Context(), id, active)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AdminHandler) GenerateInvite(c *gin.Context) {
	adminID, _ := c.Get("user_id")
	id, _ := adminID.(uuid.UUID)

	token, err := h.generateInviteUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"token": token,
		"url":   "/register/" + token,
	})
}
