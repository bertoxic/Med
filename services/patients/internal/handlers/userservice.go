package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bertoxic/med/services/patient-service/internal/models"
	"github.com/bertoxic/med/services/patient-service/internal/transport/httputil"
	"github.com/gin-gonic/gin"
)

func SignUp(ctx *gin.Context) {
	user := &models.Patient{}
	ctx.ShouldBind(user)
	var userD UserDetails
	userD.Email = user.Email
	userD.FirstName = user.FirstName
	userD.LastName = user.LastName
	userD.UserType = "patient"

	//opt, err := GenerateOTP(9)
	// if err != nil {
	//     httputil.WriteJson(ctx,false,200, &models.JsonResponse{})
	//     return
	// }

	var tokens Tokens
	tokens, err := GenerateTokenViaRestAPI(userD)
	if err != nil {
		log.Println(err)
		return
	}
	data := gin.H{
		"token":         tokens.Token,
		"refresh_token": tokens.RefreshToken,
	}
	httputil.WriteJson(ctx, true, 200, &models.JsonResponse{
		Data:    data,
		Success: true,
		Message: "suceess obtaining json response",
	})

}



func RegisterServiceGrpc(ctx *gin.Context) {
	user := &models.UserDetails{}
	jsondata := models.JsonResponse{}
	ctx.ShouldBind(user)
	id, err := RegisterViaGRPC(user)
	if err != nil {
		log.Println(err)
	}
	newdata, err := json.Marshal(id)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(newdata, &jsondata)
	if err != nil || !jsondata.Success {
		ctx.JSON(jsondata.Error.Code, jsondata)
		return
	}

	ctx.JSON(200, id)
}

func LoginUser(ctx *gin.Context) {
	var UserDetails = models.UserDetails{}
	ctx.ShouldBind(&UserDetails)
	ms, err := LoginUserViaGrpc(UserDetails)
	if err != nil {
		ctx.JSON(ms.Error.Code, ms)
	}

	ctx.JSON(200, ms)
}
func GenerateOTPResponse(ctx *gin.Context) {
	user := &models.UserDetails{}
	ctx.ShouldBind(user)

	otp, err := GenerateOTP(*user, 8)
	if err != nil {
		return
	}
	type OTP struct {
		Otp string `json:"otp"`
	}
	var otpx OTP
	otpx.Otp = otp
	json.Marshal(otpx)
	httputil.WriteJson(ctx, true, 200, &models.JsonResponse{
		Success: true,
		Message: "otp generated successfully",
		Data:    otpx,
	})
}

func GenerateTokenViaRestAPI(user UserDetails) (Tokens, error) {
	var authresp = models.JsonResponse{}
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("error: unable to marshall userdetails")
		return Tokens{}, err
	}
	authServicsURL := "http:// 0.0.0.0:9000/signup"
	request, err := http.NewRequest(http.MethodPost, authServicsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("error: unable to send request to auth")
		return Tokens{}, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("errrrrrroxx sending requuest")
		return Tokens{}, err
	}
	log.Println(response)
	var token = Tokens{}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(body, &authresp)
	if err != nil {
		return Tokens{}, err
	}
	jdata := authresp.Data.(map[string]interface{})
	jx, _ := json.Marshal(jdata)
	//token = authresp.Data.(Tokens)
	_ = json.Unmarshal(jx, &token)

	//log.Println(body)
	log.Println(authresp, "vvvvvvvvvvvvvvvvvvvvvvv")
	return token, nil
}

func VerifyOTP(ctx *gin.Context) {
	type oTp struct {
		otpstring string
	}
	var otp oTp
	ctx.ShouldBind(otp)

}
