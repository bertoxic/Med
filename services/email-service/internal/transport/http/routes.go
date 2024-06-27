package routes

import (
	"net/http"

	handler "github.com/bertoxic/med/services/email-service/internal/handlers"
	"github.com/bertoxic/med/services/email-service/internal/models"
	"github.com/bertoxic/med/services/email-service/internal/transport/httputil"
	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/authenticate", func(ctx *gin.Context) {
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
	router.POST("/sendotp", handler.SendMailTORecipient)

	return router
}
