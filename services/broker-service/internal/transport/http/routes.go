package routes

import (
	"net/http"

	handler "github.com/bertoxic/med/services/broker-service/internal/handlers"
	"github.com/bertoxic/med/services/broker-service/internal/middlewares"
	"github.com/bertoxic/med/services/broker-service/internal/models"
	"github.com/bertoxic/med/services/broker-service/internal/transport/httputil"
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
	router.GET("/login", func(ctx *gin.Context) {
		httputil.WriteJson(ctx, true, 200, &models.JsonResponse{
			Success: true,
			Message: "duely authenticated",
			Data:    "no data currently",
		})
	})
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "authenticated no wow",
		})
	})
	router.POST("/signup", handler.SignUp)

	return router
}
