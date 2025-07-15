package postgres

import (
	"strings"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/lib/pq"
)

// DatabaseOperation represents different types of database operations
type DatabaseOperation string

const (
	OpCreate DatabaseOperation = "create"
	OpRead   DatabaseOperation = "read"
	OpUpdate DatabaseOperation = "update"
	OpDelete DatabaseOperation = "delete"
)

// ErrorContext provides additional context for error handling
type ErrorContext struct {
	Operation DatabaseOperation
	Table     string
}

// HandleDatabaseError provides centralized PostgreSQL error handling
func HandleDatabaseError(err error, ctx ErrorContext) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return errs.NewInternalError("query execution failed", err, nil)
	}

	// Handle common errors that apply to all operations
	switch pqErr.Code {
	case "08003", "08006":
		// Connection errors
		return errs.NewDatabaseError("database connection lost", err, nil)
	case "53300":
		// Too many connections
		return errs.NewDatabaseError("database connection limit reached", err, nil)
	case "53400":
		// Configuration limit exceeded
		return errs.NewDatabaseError("database configuration limit exceeded", err, nil)
	case "57014":
		// Query canceled
		return errs.NewDatabaseError("query was canceled", err, nil)
	case "40001":
		// Serialization failure
		return errs.NewDatabaseError("transaction conflict, please retry", err, nil)
	}

	// Handle operation-specific errors
	switch ctx.Operation {
	case OpCreate:
		return handleCreateErrors(pqErr, ctx)
	case OpRead:
		return handleReadErrors(pqErr, ctx)
	case OpUpdate:
		return handleUpdateErrors(pqErr, ctx)
	case OpDelete:
		return handleDeleteErrors(pqErr, ctx)
	}

	// Default case for unhandled errors
	return errs.NewInternalError("postgres error", err, map[string]any{
		"postgresErrCode":    pqErr.Code,
		"postgresErrMessage": pqErr.Message,
		"postgresErrDetail":  pqErr.Detail,
		"operation":          ctx.Operation,
		"table":              ctx.Table,
	})
}

func handleCreateErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "23505":
		// Unique constraint violation
		if ctx.Table == "users" {
			if strings.Contains(pqErr.Detail, "email") {
				return errs.NewAlreadyExistsError("email already exists", map[string]any{"email": extractEmailFromDetail(pqErr.Detail)})
			} else if strings.Contains(pqErr.Detail, "uuid") {
				return errs.NewAlreadyExistsError("user ID already exists", nil)
			}
		}
		return errs.NewAlreadyExistsError("duplicate value violates unique constraint", map[string]any{"detail": pqErr.Detail})
	case "23503":
		// Foreign key violation
		return errs.NewDatabaseError("invalid reference to related resource", pqErr, nil)
	case "23502":
		// Not-null constraint violation
		return errs.NewDatabaseError("missing required field", pqErr, map[string]any{"column": pqErr.Column})
	case "23514":
		// Check constraint violation
		return errs.NewValidationError("data violates check constraint", map[string]any{"constraint": pqErr.Constraint})
	case "22001":
		// String data right truncation
		return errs.NewValidationError("field value too long", map[string]any{"column": pqErr.Column})
	case "22P02":
		// Invalid text representation
		return errs.NewValidationError("invalid data format", map[string]any{"detail": pqErr.Detail})
	}
	return nil // Will fall through to default handling
}

func handleReadErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "22P02":
		// Invalid text representation
		return errs.NewValidationError("invalid query parameter format", map[string]any{"detail": pqErr.Detail})
	case "42703":
		// Undefined column
		return errs.NewInternalError("query references undefined column", pqErr, map[string]any{"column": pqErr.Column})
	case "42P01":
		// Undefined table
		return errs.NewInternalError("query references undefined table", pqErr, nil)
	}
	return nil
}

func handleUpdateErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "23505":
		// Unique constraint violation
		if ctx.Table == "users" && strings.Contains(pqErr.Detail, "email") {
			return errs.NewAlreadyExistsError("email already exists", map[string]any{"email": extractEmailFromDetail(pqErr.Detail)})
		}
		return errs.NewAlreadyExistsError("duplicate value violates unique constraint", map[string]any{"detail": pqErr.Detail})
	case "23503":
		// Foreign key violation
		return errs.NewDatabaseError("invalid reference to related resource", pqErr, nil)
	case "23502":
		// Not-null constraint violation
		return errs.NewDatabaseError("missing required field", pqErr, map[string]any{"column": pqErr.Column})
	case "23514":
		// Check constraint violation
		return errs.NewValidationError("data violates check constraint", map[string]any{"constraint": pqErr.Constraint})
	case "22001":
		// String data right truncation
		return errs.NewValidationError("field value too long", map[string]any{"column": pqErr.Column})
	case "22P02":
		// Invalid text representation
		return errs.NewValidationError("invalid data format", map[string]any{"detail": pqErr.Detail})
	}
	return nil
}

func handleDeleteErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "23503":
		// Foreign key violation (trying to delete referenced row)
		return errs.NewDatabaseError("cannot delete record that is referenced by other records", pqErr, nil)
	case "22P02":
		// Invalid text representation
		return errs.NewValidationError("invalid query parameter format", map[string]any{"detail": pqErr.Detail})
	}
	return nil
}

// Helper function to extract email from error detail
func extractEmailFromDetail(detail string) string {
	if strings.Contains(detail, "email") {
		start := strings.Index(detail, "=(")
		if start != -1 {
			start += 2
			end := strings.Index(detail[start:], ")")
			if end != -1 {
				return detail[start : start+end]
			}
		}
	}
	return ""
}
