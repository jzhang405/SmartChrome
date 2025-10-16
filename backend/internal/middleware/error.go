package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle if there are errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Skip if response has already been written
			if c.Writer.Status() != http.StatusOK {
				return
			}

			switch e := err.Err.(type) {
			case *AppError:
				c.JSON(e.StatusCode, ErrorResponse{
					Error:   e.Message,
					Code:    e.Code,
					Details: e.Details,
				})
			default:
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error: "Internal server error",
					Code:  "INTERNAL_ERROR",
					Details: err.Error(),
				})
			}
		}
	}
}

type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Details    interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func NewAppErrorWithDetails(statusCode int, code, message string, details interface{}) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// Common error types
var (
	ErrBadRequest       = NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Bad request")
	ErrUnauthorized     = NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	ErrForbidden        = NewAppError(http.StatusForbidden, "FORBIDDEN", "Forbidden")
	ErrNotFound         = NewAppError(http.StatusNotFound, "NOT_FOUND", "Resource not found")
	ErrInternalServer   = NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	ErrServiceUnavailable = NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service unavailable")
)