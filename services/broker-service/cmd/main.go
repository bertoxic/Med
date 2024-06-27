package main

import (
	"log"
	"net/http"

	routes "github.com/bertoxic/med/services/broker-service/internal/transport/http"
)

func main() {
    r := routes.Router()
	srv := &http.Server{
        Addr: "0.0.0.0:8087",
        Handler: r,
    }

    err := srv.ListenAndServe()
    if err != nil {
        log.Println("Eroooxxx "+err.Error())
        return
    }
}