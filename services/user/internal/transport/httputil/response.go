package httputil

import (
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func ReadJson(ctx *gin.Context) {

}

func WriteJson(ctx *gin.Context, statuscode int, message string, data interface{}) {
	var response ApiResponse
	response.Data = data
	response.Message = message
	status := statuscode

	ctx.JSON(status, response)
}

func ErrorJson() {

}
