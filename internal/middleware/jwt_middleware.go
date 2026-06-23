package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"domesv2/config"
	"domesv2/pkg/response"
	"go.uber.org/zap"
)

func JWTMiddleware() fiber.Handler {
	cfg := config.AppConfig
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header", "TOKEN_MISSING")
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Unauthorized(c, "Invalid authorization header format", "INVALID_TOKEN_FORMAT")
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			zap.L().Warn("JWT validation failed", zap.Error(err))
			return response.Unauthorized(c, "Invalid or expired token", "INVALID_TOKEN")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user ID in context for use in handlers
			if userID, ok := claims["user_id"].(float64); ok {
				c.Locals("user_id", uint(userID))
			}
			if email, ok := claims["email"].(string); ok {
				c.Locals("user_email", email)
			}
			return c.Next()
		}

		return response.Unauthorized(c, "Invalid token claims", "INVALID_TOKEN_CLAIMS")
	}
}
