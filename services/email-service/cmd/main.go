package main

import (
	"log"
	"net/http"

	"github.com/bertoxic/med/services/email-service/internal/handlers"
	routes "github.com/bertoxic/med/services/email-service/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../../.env")
    if err != nil {
        log.Println("Error loading .env file")
    }
	r := routes.Router()
	srv := &http.Server{
		Addr:    "0.0.0.0:8083",
		Handler: r,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println("Eroooxxx " + err.Error())
		return
	}
	handler.Initx()
}
