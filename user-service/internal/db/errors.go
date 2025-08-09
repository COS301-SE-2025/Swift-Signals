package db

import (
	"strings"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/lib/pq"
)

type DatabaseOperation string

const (
	OpCreate DatabaseOperation = "create"
	OpRead   DatabaseOperation = "read"
	OpUpdate DatabaseOperation = "update"
	OpDelete DatabaseOperation = "delete"
)

type ErrorContext struct {
	Operation DatabaseOperation
	Table     string
}

func HandleDatabaseError(err error, ctx ErrorContext) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return errs.NewInternalError("query execution failed", err, nil)
	}

	switch pqErr.Code {
	case "08003", "08006":
		return errs.NewDatabaseError("database connection lost", err, nil)
	case "53300":
		return errs.NewDatabaseError("database connection limit reached", err, nil)
	case "53400":
		return errs.NewDatabaseError("database configuration limit exceeded", err, nil)
	case "57014":
		return errs.NewDatabaseError("query was canceled", err, nil)
	case "40001":
		return errs.NewDatabaseError("transaction conflict, please retry", err, nil)
	}

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
		if ctx.Table == "users" {
			if strings.Contains(pqErr.Detail, "email") {
				return errs.NewAlreadyExistsError(
					"email already exists",
					map[string]any{"email": extractEmailFromDetail(pqErr.Detail)},
				)
			} else if strings.Contains(pqErr.Detail, "uuid") {
				return errs.NewAlreadyExistsError("user ID already exists", nil)
			}
		}
		return errs.NewAlreadyExistsError(
			"duplicate value violates unique constraint",
			map[string]any{"detail": pqErr.Detail},
		)
	case "23503":
		return errs.NewInternalError("invalid reference to related resource", pqErr, nil)
	case "23502":
		return errs.NewInternalError(
			"missing required field",
			pqErr,
			map[string]any{"column": pqErr.Column},
		)
	case "22001":
		return errs.NewValidationError(
			"field value too long",
			map[string]any{"column": pqErr.Column},
		)
	}
	return nil
}

func handleReadErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "22P02":
		return errs.NewValidationError(
			"invalid query parameter format",
			map[string]any{"detail": pqErr.Detail},
		)
	case "42703":
		return errs.NewInternalError(
			"query references undefined column",
			pqErr,
			map[string]any{"column": pqErr.Column},
		)
	case "42P01":
		return errs.NewInternalError("query references undefined table", pqErr, nil)
	}
	return nil
}

func handleUpdateErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "23505":
		if ctx.Table == "users" && strings.Contains(pqErr.Detail, "email") {
			return errs.NewInternalError(
				"email already in use by another user", pqErr,
				map[string]any{"email": extractEmailFromDetail(pqErr.Detail)},
			)
		}
		return errs.NewDatabaseError(
			"value conflicts with existing record", pqErr,
			map[string]any{"detail": pqErr.Detail},
		)
	case "23503":
		return errs.NewValidationError(
			"cannot update: referenced record does not exist",
			map[string]any{
				"detail": pqErr.Detail,
			},
		)
	case "23502":
		return errs.NewValidationError(
			"cannot set required field to empty",
			map[string]any{"column": pqErr.Column},
		)
	case "42703":
		return errs.NewValidationError(
			"cannot update: column does not exist",
			map[string]any{"column": pqErr.Column},
		)
	}
	return nil
}

func handleDeleteErrors(pqErr *pq.Error, ctx ErrorContext) error {
	switch pqErr.Code {
	case "23503":
		return errs.NewDatabaseError(
			"cannot delete record that is referenced by other records",
			pqErr,
			nil,
		)
	case "22P02":
		return errs.NewValidationError(
			"invalid query parameter format",
			map[string]any{"detail": pqErr.Detail},
		)
	}
	return nil
}

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
