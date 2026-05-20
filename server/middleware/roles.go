package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// RoleMiddleware создает middleware для проверки ролей
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*models.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Проверяем, есть ли роль пользователя в списке разрешенных
		hasRole := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "insufficient permissions",
				"required_roles": allowedRoles,
				"user_role":      claims.Role,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware проверяет, что пользователь является администратором
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// UserMiddleware проверяет, что пользователь является обычным пользователем или администратором
func UserMiddleware() gin.HandlerFunc {
	return RoleMiddleware("user", "admin")
}

// GetUserIDFromContext извлекает ID пользователя из контекста
func GetUserIDFromContext(c *gin.Context) (int, bool) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return 0, false
	}

	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		return 0, false
	}

	return claims.UserID, true
}

// GetUserRoleFromContext извлекает роль пользователя из контекста
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return "", false
	}

	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		return "", false
	}

	return claims.Role, true
}

// ContentAccessMiddleware проверяет доступ к контенту курса/мастер-класса
func ContentAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Заглушка: в реальной реализации здесь должна быть проверка активной покупки
		// Для фазы 2 пропускаем запрос
		// Пример логики:
		// 1. Получить product_id из параметров запроса
		// 2. Проверить наличие активной покупки в таблице purchases
		// 3. Если покупки нет - вернуть 403
		c.Next()
	}
}

// OptionalAuthMiddleware - middleware, который устанавливает claims если токен есть, но не требует аутентификации
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.Next()
			return
		}

		// Парсим токен, но не прерываем если он невалидный
		// Это middleware только для установки claims если токен валидный
		// Реальная проверка токена будет в AuthMiddleware
		// Пока просто пропускаем
		c.Next()
	}
}
