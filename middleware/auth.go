package middleware

import (
	"net/http"
	"store/config"
	"store/repository"
	"store/utils"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	redis "github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(authRepository repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "Authorization header is missing", nil, true)
			return
		}

		// Expect "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "Invalid Authorization header format", nil, true)
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "Invalid token : "+err.Error(), nil, true)

			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
				utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "expired token", nil, true)
				return
			}

			username := claims["username"]
			accessKey := username.(string)
			tokenRedis, err := authRepository.GetByKey(accessKey)
			if err != nil {
				if err == redis.Nil {
					utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "AccessToken missing on redis", nil, true)
					return
				}
			}

			if tokenRedis != tokenString {
				utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "AccessToken not match", nil, true)
				return
			}
		} else {
			utils.Response(c, http.StatusUnauthorized, nil, nil, nil, "Invalid token claims", nil, true)
			return
		}

		c.Next()
	}
}
