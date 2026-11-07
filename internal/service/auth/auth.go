package auth

import (
	"context"
	"strings"

	"github.com/MuhAndriansyah/yorkshire/internal/apperror"
	db "github.com/MuhAndriansyah/yorkshire/internal/database/sqlc"
	"github.com/MuhAndriansyah/yorkshire/internal/jwt"
	"github.com/MuhAndriansyah/yorkshire/internal/security"
	"github.com/google/uuid"
)

func (s *Auth) Register(ctx context.Context, input RegisterInput) (AuthOutput, error) {
	passwordHashed, err := security.HashAndSalt(input.Password)
	if err != nil {
		return AuthOutput{}, apperror.Internal(err, "failed to process registration")
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return AuthOutput{}, apperror.Internal(err, "failed to generate user id")
	}

	user, err := s.Repo.CreateUser(ctx, db.CreateUserParams{
		ID:           newUUID,
		Username:     strings.TrimSpace(input.Username),
		Email:        strings.ToLower(strings.TrimSpace(input.Email)),
		PasswordHash: passwordHashed,
	})
	if err != nil {
		return AuthOutput{}, translateAuthPgError(err)
	}

	token, err := s.Token.Generate(jwt.UserIdentity{
		UserID: user.ID.String(),
		Email:  user.Email,
	})
	if err != nil {
		return AuthOutput{}, apperror.Internal(err, "failed to generate access token")
	}

	return AuthOutput{
		User:        toUserOutput(user),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.TokenDuration.Seconds()),
	}, nil
}

func (s *Auth) Login(ctx context.Context, input LoginInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if apperror.IsNoRows(err) {
			return AuthOutput{}, apperror.ErrUnauthorized("invalid email or password")
		}
		return AuthOutput{}, translateAuthPgError(err)
	}

	if !security.ComparePasswords(user.PasswordHash, input.Password) {
		return AuthOutput{}, apperror.ErrUnauthorized("invalid email or password")
	}

	token, err := s.Token.Generate(jwt.UserIdentity{
		UserID: user.ID.String(),
		Email:  user.Email,
	})
	if err != nil {
		return AuthOutput{}, apperror.Internal(err, "failed to generate access token")
	}

	return AuthOutput{
		User:        toUserOutput(user),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.TokenDuration.Seconds()),
	}, nil
}
