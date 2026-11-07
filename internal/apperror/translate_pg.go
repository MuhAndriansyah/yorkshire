package apperror

import (
	"context"
	"errors"
)

// TranslateOptions maps raw sqlc/pgx errors to domain-level application errors.
type TranslateOptions struct {
	NotFound string

	UniqueViolations       map[string]string
	DefaultUniqueViolation string

	FKViolations       map[string]*AppError
	DefaultFKViolation string

	Internal string
}

func TranslatePgError(err error, opts TranslateOptions) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	if opts.NotFound != "" && IsNoRows(err) {
		return NotFound(opts.NotFound)
	}

	if constraint, ok := UniqueViolationConstraint(err); ok {
		if msg, found := opts.UniqueViolations[constraint]; found {
			return Conflict(msg)
		}
		if opts.DefaultUniqueViolation != "" {
			return Conflict(opts.DefaultUniqueViolation)
		}
	}

	if constraint, ok := ForeignKeyViolationConstraint(err); ok {
		if appErr, found := opts.FKViolations[constraint]; found {
			return appErr
		}
		if opts.DefaultFKViolation != "" {
			return Conflict(opts.DefaultFKViolation)
		}
	}

	if opts.Internal == "" {
		opts.Internal = "internal error"
	}
	return Internal(err, opts.Internal)
}
