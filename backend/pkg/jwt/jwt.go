package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) *TokenManager {
	return &TokenManager{
		signingKey: signingKey,
	}
}

func (m *TokenManager) CreateToken(userID int64, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, m.validateSigningMethod)
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["sub"].(float64))
		return userID, nil
	}
	return 0, fmt.Errorf("invalid token")
}

func (m *TokenManager) validateSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(m.signingKey), nil
}
