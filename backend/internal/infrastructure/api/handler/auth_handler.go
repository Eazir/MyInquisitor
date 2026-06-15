package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appAuth "github.com/myinquisitor/backend/internal/application/auth"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type AuthHandler struct {
	registerUC *appAuth.RegisterUseCase
	loginUC    *appAuth.LoginUseCase
	refreshUC  *appAuth.RefreshUseCase
}

func NewAuthHandler(registerUC *appAuth.RegisterUseCase, loginUC *appAuth.LoginUseCase, refreshUC *appAuth.RefreshUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.registerUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, appAuth.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		if errors.Is(err, appAuth.ErrInvalidInviteToken) || errors.Is(err, appAuth.ErrInviteTokenUsed) {
			response.Error(c, http.StatusBadRequest, "INVITE_ERROR", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, appAuth.ErrEmailNotFound) {
			response.Error(c, http.StatusNotFound, "EMAIL_NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, appAuth.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "WRONG_PASSWORD", err.Error())
			return
		}
		if errors.Is(err, appAuth.ErrUserInactive) {
			response.Error(c, http.StatusUnauthorized, "INACTIVE_ACCOUNT", err.Error())
			return
		}
		if errors.Is(err, appAuth.ErrAccountIssue) {
			response.Error(c, http.StatusInternalServerError, "ACCOUNT_ERROR", "account configuration error. please contact support")
			return
		}
		c.Error(err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred. please try again later")
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input dto.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "invalid request body", err.Error())
		return
	}

	result, err := h.refreshUC.Execute(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, appAuth.ErrInvalidRefreshToken) || errors.Is(err, appAuth.ErrUserInactive) {
			response.Error(c, http.StatusUnauthorized, "AUTH_ERROR", err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
		return
	}

	response.Success(c, http.StatusOK, result)
}
