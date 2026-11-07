package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	Generate(identity UserIdentity) (string, error)
	Validate(token string) (*UserIdentity, error)
}

type JWTService struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	issuer        string
	audience      string
	tokenDuration time.Duration
}

type NewJWTOptions struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	Issuer        string
	Audience      string
	TokenDuration time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type UserIdentity struct {
	UserID string
	Email  string
}

func New(opts *NewJWTOptions) *JWTService {
	return &JWTService{
		privateKey:    opts.PrivateKey,
		publicKey:     opts.PublicKey,
		issuer:        opts.Issuer,
		audience:      opts.Audience,
		tokenDuration: opts.TokenDuration,
	}
}

func (j *JWTService) Generate(identity UserIdentity) (string, error) {
	now := time.Now()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewV7().String(),
			Issuer:    j.issuer,
			Subject:   identity.UserID,
			Audience:  jwt.ClaimStrings{j.audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenDuration)),
		},
		UserID: identity.UserID,
		Email:  identity.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(j.privateKey)
}

func (j *JWTService) Validate(token string) (*UserIdentity, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		token, // 1. Hindari nama variabel yang sama dengan hasil kembalian
		claims,
		func(t *jwt.Token) (any, error) {
			// 2. Gunakan pengecekan tipe interface (*jwt.SigningMethodRSA)
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return j.publicKey, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithAudience(j.audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}), // 3. Pembatasan algoritma bawaan jwt v5
	)

	if err != nil {
		return nil, fmt.Errorf("failed parse jwt token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &UserIdentity{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil

}
