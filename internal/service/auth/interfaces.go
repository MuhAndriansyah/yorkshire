package auth

import "context"

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (AuthOutput, error)
}
