package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{Data: data})
}

func SuccessWithMeta(c *gin.Context, status int, data interface{}, meta *Meta) {
	c.JSON(status, APIResponse{Data: data, Meta: meta})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{Error: &APIError{Code: code, Message: message}})
}

func ValidationError(c *gin.Context, message string, details interface{}) {
	c.JSON(422, APIResponse{Error: &APIError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Details: details,
	}})
}
