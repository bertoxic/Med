package middlewares

import (
	"log"
	"net/http"
	"strings"

	handler "github.com/bertoxic/med/services/broker-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
            log.Println("no authheader found")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, message := handler.VerifyToken(tokenString)
		if message != "" {
            log.Println(" claims error", claims,message)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("claims", claims)
		log.Println(claims)
		ctx.Next()
	}

}