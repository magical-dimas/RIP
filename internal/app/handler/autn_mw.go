package handler

import (
	"net/http"
	"rip_project/internal/app/role"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

const jwtSecret = "secret_neclear_key_228"
const jwtPrefix = "Bearer "

type JWTClaims struct {
	jwt.StandardClaims
	EngineerID uint      `json:"engineer_id"`
	Role       role.Role `json:"role"`
}

func (h *Handler) WithAuthCheck(allowedRoles ...role.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, jwtPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Необходим Bearer токен авторизации"})
			return
		}

		tokenStr := authHeader[len(jwtPrefix):]

		err := h.Repository.GetRedis().Get(c.Request.Context(), "blacklist:"+tokenStr).Err()
		if err == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен недействителен"})
			return
		} else if err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки токена в Redis"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный токен"})
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ошибка чтения данных токена"})
			return
		}

		hasAccess := false
		for _, r := range allowedRoles {
			if claims.Role == r || claims.Role == role.Technician {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен для вашей роли"})
			return
		}

		c.Set("engineer_id", claims.EngineerID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (h *Handler) getEngineerID(c *gin.Context) uint {
	id, exists := c.Get("engineer_id")
	if exists {
		return id.(uint)
	}
	return 0
}
