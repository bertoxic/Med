package routes

import (
	"net/http"

	handler "github.com/bertoxic/med/services/patient-service/internal/handlers"
	"github.com/bertoxic/med/services/patient-service/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api := router.Group("",middlewares.AuthMiddleware())
	api.GET("/profile", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "authenticated no wow",
		})
	})
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "userservices oohh wow",
		})
	})
	router.POST("/sendotp", handler.GenerateOTPResponse)
	router.POST("/signup", handler.RegisterServiceGrpc)
	router.POST("/login", handler.LoginUser)

	return router
}
