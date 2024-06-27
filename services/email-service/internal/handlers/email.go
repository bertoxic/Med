package handler

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/bertoxic/med/services/email-service/internal/models"
	"github.com/bertoxic/med/services/email-service/internal/transport/httputil"
	"github.com/gin-gonic/gin"
)

type EmailService struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(host string, port int, username, password string) *EmailService {
	return &EmailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (es *EmailService) SendMail(data models.EmailData) error {
	Initx()
	auth := smtp.PlainAuth("", es.Username, es.Password, es.Host)

	msg := fmt.Sprintf("from: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s ", es.Username, data.To, data.Subject, data.Body)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", es.Host, es.Port), auth, es.Username, []string{data.To}, []byte(msg))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SendMailTORecipient(ctx *gin.Context) {
    num, err := strconv.Atoi(smtpPort)
    if err != nil {
        fmt.Println("Error converting string to int:", err)
        return
    }
	emailService := NewEmailService(smtpHost, num, smtpUser, smtpPass)
	var user = &models.UserDetails{}
	err = ctx.ShouldBind(user)
	log.Println(user)
	if err != nil {

		log.Println("error ocuxxxxz", err)
		return
	}
	emailData := models.EmailData{
		To:      emailRecipient,
		Subject: "Verification mail",
		Body:    "verfici code is 23312",
	}

	err = emailService.SendMail(emailData)
	if err != nil {
		log.Println(err.Error())
		httputil.WriteJson(ctx, false, 400, &models.JsonResponse{
			Success: false,
			Message: fmt.Sprintf("err occured :%s", err.Error()),
		})
		return
	}
	httputil.WriteJson(ctx, true, 200, &models.JsonResponse{
		Success: true,
		Message: "emaal sent successfully",
	})
}
