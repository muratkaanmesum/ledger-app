package jwt

import (
	"errors"
	"fmt"
	"os"
	"ptm/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Id       uint   `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set in environment variables")
	}

	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"id":       user.ID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (*CustomClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set in environment variables")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
