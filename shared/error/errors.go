// Package error provides a structured error type and helper functions
package error

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServiceError represents a structured error with context
type ServiceError struct {
	Code    ErrorCode      `json:"code"`
	Message string         `json:"message"`
	Cause   error          `json:"-"` // Don't serialize the original error
	Context map[string]any `json:"context,omitempty"`
}

// Error returns the error message
func (e *ServiceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %s)", e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the original error
func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrValidation    ErrorCode = "VALIDATION_ERROR"
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrForbidden     ErrorCode = "FORBIDDEN"
	ErrDatabase      ErrorCode = "DB_ERROR"
	ErrInternal      ErrorCode = "INTERNAL_ERROR"
	ErrExternal      ErrorCode = "EXTERNAL_ERROR"
)

// NewValidationError creates a new validation error
func NewValidationError(message string, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrValidation,
		Message: message,
		Context: context,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrNotFound,
		Message: message,
		Context: context,
	}
}

// NewAlreadyExistsError creates a new already exists error
func NewAlreadyExistsError(message string, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrAlreadyExists,
		Message: message,
		Context: context,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrUnauthorized,
		Message: message,
		Context: context,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrForbidden,
		Message: message,
		Context: context,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, cause error, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrDatabase,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, cause error, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrInternal,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// NewExternalError creates a new external error
func NewExternalError(message string, cause error, context map[string]any) *ServiceError {
	return &ServiceError{
		Code:    ErrExternal,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// HandleServiceError converts service errors to appropriate gRPC status errors
func HandleServiceError(err error) error {
	if err == nil {
		return nil
	}

	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		switch svcErr.Code {
		case ErrValidation:
			return status.Error(codes.InvalidArgument, svcErr.Message)
		case ErrNotFound:
			return status.Error(codes.NotFound, svcErr.Message)
		case ErrAlreadyExists:
			return status.Error(codes.AlreadyExists, svcErr.Message)
		case ErrUnauthorized:
			return status.Error(codes.Unauthenticated, svcErr.Message)
		case ErrForbidden:
			return status.Error(codes.PermissionDenied, svcErr.Message)
		case ErrDatabase, ErrInternal, ErrExternal:
			return status.Error(codes.Internal, "internal server error")
		}
	}

	return status.Error(codes.Internal, "internal server error")
}
