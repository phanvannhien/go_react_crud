package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// ContextKeyUserID is the key used to store user ID in echo.Context.
const ContextKeyUserID = "user_id"

// JWTAuth creates an Echo middleware that validates JWT Bearer tokens.
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"data": nil,
					"error": map[string]string{
						"code":    "unauthorized",
						"message": "Missing authorization header",
					},
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"data": nil,
					"error": map[string]string{
						"code":    "unauthorized",
						"message": "Invalid authorization header format",
					},
				})
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"data": nil,
					"error": map[string]string{
						"code":    "unauthorized",
						"message": "Invalid or expired token",
					},
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"data": nil,
					"error": map[string]string{
						"code":    "unauthorized",
						"message": "Invalid token claims",
					},
				})
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"data": nil,
					"error": map[string]string{
						"code":    "unauthorized",
						"message": "Invalid user ID in token",
					},
				})
			}

			c.Set(ContextKeyUserID, userID)

			return next(c)
		}
	}
}

// GetUserIDFromContext extracts the user ID from the echo.Context.
func GetUserIDFromContext(c echo.Context) string {
	userID, _ := c.Get(ContextKeyUserID).(string)
	return userID
}
