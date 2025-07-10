package helper

import (
	"fmt"
	"store/config"
	"store/constant"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string, now time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Duration(config.AccessTokenExpiration) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(config.JwtSecret))
}

func GenerateRefreshToken(username string, now time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"type":     "refresh",
		"exp":      time.Now().Add(time.Duration(config.RefreshTokenExpiration) * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(config.RefreshSecret))
}

func GenerateRefreshTokenFromOld(token string, now time.Time) (string, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(config.RefreshSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return "", "", fmt.Errorf("token is not a refresh token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid token: username not found")
	}

	newClaims := jwt.MapClaims{
		"username": username,
		"type":     "refresh",
		"exp":      time.Now().Add(time.Duration(config.RefreshTokenExpiration) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	refreshToken, err := newToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", "", constant.ErrJwtSigned
	}

	return username, refreshToken, nil
}
