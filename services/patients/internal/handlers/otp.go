package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bertoxic/med/services/patient-service/internal/models"
)

const (
	maxOTPLength   = 8
	otpExpiryDur   = 5 * time.Minute
	maxOTPAttempts = 3
	encodingBase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

func SendOTP() {}

func GenerateOTP(user models.UserDetails, length int) (string, error) {
	authreq := &models.AuthPayload{
		UserDetails: user,
	}
	authresp := &models.JsonResponse{}
	jsonData, err := json.Marshal(authreq)
	if err != nil {
		return "", err
	}
	//buffer := make([]byte, length)
	// n, err := rand.Read(buffer)
	// if n != len(buffer) || err != nil {
	// 	return "", err
	// }

	// otpNum := new(big.Int).SetBytes(buffer).String()
	// log.Println(otpNum)
	// if len(otpNum) < length {
	// 	otpNum = fmt.Sprintf("%0*s", length, otpNum)
	// }
	// log.Println(otpNum)
	// otp := otpNum[:length]
	authURL := "http://localhost:8085/otp"
	request, err := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(jsonData))	
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	response, err :=client.Do(request)
	if err != nil {
		log.Println("errrrrrroxx sending requuest")
		return "", err
	}
	defer response.Body.Close()
	log.Println(response.Body)
	body, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(body,&authresp)
	if err != nil {
		log.Println(err,"cannot unmashallli")
		return "", err
	}
	log.Println(authresp)
	jdata := authresp.Data.(map[string]interface{})
	var OTP  struct{
		Otp string
	}
	jx, _ := json.Marshal(jdata)
	err = json.Unmarshal(jx,&OTP)
	if err != nil {
		log.Println(err,"ccccccccccc")
		return "", err
	}
	log.Println(OTP)

	return OTP.Otp, nil
}

func verifyOTP(generatedOTP, userOTP string) bool {
	attempts := 0
	validUntil := time.Now().Add(otpExpiryDur)

	for attempts < maxOTPAttempts && time.Now().Before(validUntil) {
		if userOTP == generatedOTP {
			return true
		}
		attempts++
	}

	return false
}
