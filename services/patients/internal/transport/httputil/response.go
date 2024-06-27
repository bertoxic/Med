package httputil

import (
	"github.com/bertoxic/med/services/patient-service/internal/models"
	"github.com/gin-gonic/gin"
)



func ReadJson(ctx *gin.Context) {

}

func WriteJson(ctx *gin.Context, success bool, statuscode int, response *models.JsonResponse) {
	var resp models.JsonResponse
	ctx.Request.Header.Set("Content-Type", "application/json")
	header := ctx.Writer.Header()
	header.Write(ctx.Writer)
	resp.Success = success
	var status int
	if !success {
		resp.Data = response.Data
		resp.Message = response.Message
		status = statuscode
		resp.Error = &models.ErrorJson{Code: 400,Message: "erroro???"}

	} else {		
		resp.Data = response.Data
		resp.Message = response.Message
		status = 200
	}

	ctx.JSON(status, response)
}

func Success(data interface{}, message string) *models.JsonResponse {
	return &models.JsonResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func Fail(statuscode int, err error) *models.JsonResponse {
	return &models.JsonResponse{
		Success: false,
		Error: &models.ErrorJson{
			Code:    statuscode,
			Message: err.Error(),
		},
	}
}
