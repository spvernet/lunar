package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Placeholder: añade tu validación (JWT, API key, etc.)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Si quieres proteger lecturas, valida aquí.
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Next()
	}
}

// Healthcheck rápido (opcional)
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
