package auth

import (
	"time"

	db "github.com/MuhAndriansyah/yorkshire/internal/database/sqlc"
	"github.com/MuhAndriansyah/yorkshire/internal/jwt"
)

type Auth struct {
	Repo          db.Querier
	Token         jwt.TokenService
	TokenDuration time.Duration
}

func NewAuth(repo db.Querier, token jwt.TokenService, tokenDuration time.Duration) *Auth {
	return &Auth{
		Repo:          repo,
		Token:         token,
		TokenDuration: tokenDuration,
	}
}
