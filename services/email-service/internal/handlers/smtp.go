package handler

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

var smtpHost string = os.Getenv("SMTPHOST")
var smtpPort string = os.Getenv("SMTPPORT")
var smtpUser string = os.Getenv("SMTPUSER")
var smtpPass string = os.Getenv("SMTPPASS")
var emailRecipient string = os.Getenv("EMAIL_RECEPIENT")

func Initx() {
	// SMTP server configuration

	// Message details
	from := smtpUser
	to := []string{emailRecipient}
	subject := "Test Email"
	body := "This is a test email sent using Go."

	// Connect to the SMTP server
	addr := smtpHost + ":" + smtpPort
	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("Failed to connect to SMTP server: %v", err)
		return
	}
	defer conn.Close()

	// Send the initial greeting
	if err = conn.Hello("localhost"); err != nil {
		log.Printf("Failed to send HELO: %v", err)
		return
	}

	// Authenticate with the SMTP server
	if err = conn.Auth(smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)); err != nil {
		log.Printf("Failed to authenticate: %v", err)
		return
	}

	// Set the sender's email address
	if err = conn.Mail(from); err != nil {
		log.Printf("Failed to set sender: %v", err)
		return
	}

	// Set the recipient's email address(es)
	for _, addr := range to {
		if err = conn.Rcpt(addr); err != nil {
			log.Printf("Failed to set recipient: %v", err)
			return
		}
	}

	// Send the DATA command and the email body
	wc, err := conn.Data()
	if err != nil {
		log.Printf("Failed to get data stream: %v", err)
		return
	}
	defer wc.Close()

	// Compose the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to[0], subject, body)

	// Write the email body data
	_, err = fmt.Fprint(wc, message)
	if err != nil {
		log.Printf("Failed to send message data: %v", err)
		return
	}

	// Terminate the connection
	err = conn.Quit()
	if err != nil {
		log.Printf("Failed to quit connection: %v", err)
		return
	}

	log.Println("Email sent successfully!")
}
