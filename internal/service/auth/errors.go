package auth

import "github.com/MuhAndriansyah/yorkshire/internal/apperror"

var authErrors = apperror.TranslateOptions{
	NotFound: "invalid email or password",
	UniqueViolations: map[string]string{
		"users_email_unique": "email already exists",
	},
	DefaultUniqueViolation: "duplicate user data",
	Internal:               "failed to process authentication",
}

func translateAuthPgError(err error) error {
	return apperror.TranslatePgError(err, authErrors)
}
