package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bertoxic/med/services/authentication/internal/transport/httputil"
)

func Router() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/authenticate", func(ctx *gin.Context) { httputil.WriteJson(ctx, 200, "you are about to be authenticated", nil) })
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "authenticated no wow",
		})
	})

	return router
}
