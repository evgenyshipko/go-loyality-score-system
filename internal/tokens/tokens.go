package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var jwtSecretKey = []byte("your-secret-key")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access токен действителен 15 минут
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your-app",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh токен действителен 7 дней
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// TODO: реализовать refresh-хендлер
func RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := ParseJWT(refreshToken)
	if err != nil {
		return "", fmt.Errorf("невалидный refresh токен: %v", err)
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return "", fmt.Errorf("refresh токен истек")
	}

	return GenerateAccessToken(claims.UserID)
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи токена правильный
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("невалидный токен")
	}

	return claims, nil
}
