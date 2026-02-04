// Package jwt provides function to generate and validate token
package jwt

import (
	"time"

	"github.com/SyafaHadyan/freepass-2026/internal/infra/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTItf interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type JWT struct {
	secretKey   string
	expiredTime uint
}

type Claims struct {
	ID   uuid.UUID
	Role string
	jwt.RegisteredClaims
}

func New(env *env.Env) *JWT {
	return &JWT{
		secretKey:   env.JWTSecretKey,
		expiredTime: env.JWTExpiredDays,
	}
}

func (j *JWT) GenerateToken(userID uuid.UUID, role string) (string, error) {
	claim := Claims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * 24 * time.Duration(j.expiredTime)),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)

	tokenStirng, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStirng, nil
}

func (j *JWT) ValidateToken(tokenString string) (uuid.UUID, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, err
	}

	userID := claims.ID

	return userID, nil
}
