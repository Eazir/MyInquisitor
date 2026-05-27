package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appProfile "github.com/myinquisitor/backend/internal/application/profile"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type ProfileHandler struct {
	updateProfileUC     *appProfile.UpdateProfileUseCase
	changePasswordUC    *appProfile.ChangePasswordUseCase
}

func NewProfileHandler(
	updateProfileUC *appProfile.UpdateProfileUseCase,
	changePasswordUC *appProfile.ChangePasswordUseCase,
) *ProfileHandler {
	return &ProfileHandler{
		updateProfileUC:  updateProfileUC,
		changePasswordUC: changePasswordUC,
	}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uuid.UUID)

	var input dto.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.updateProfileUC.Execute(c.Request.Context(), uid, input)
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

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uuid.UUID)

	var input dto.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	err := h.changePasswordUC.Execute(c.Request.Context(), uid, input)
	if err != nil {
		if errors.Is(err, appProfile.ErrInvalidPassword) {
			response.Error(c, http.StatusUnauthorized, "AUTH_ERROR", err.Error())
			return
		}
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "password updated successfully"})
}
