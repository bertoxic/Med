package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	handler "github.com/bertoxic/med/services/patient-service/internal/handlers"
	routes "github.com/bertoxic/med/services/patient-service/internal/transport/http"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/credentials/insecure"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
		log.Println("Error loading .env file",err)
	}
    port := os.Getenv("GRPC_PORT")	
    if port == "" {
        port = "5001" // Default to 5001 if PORT is not set
    }

    r := routes.Router()
	srv := &http.Server{
        Addr: "0.0.0.0:8086",
        Handler: r,
    }
    conn, err := grpc.Dial(
        "med-o9j9.onrender.com:443",
        grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
    )
    //  conn, err := grpc.NewClient("med-o9j9.onrender.com:"+"5001",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("grpc did not connect: %v", err)
	}
    log.Printf("listening on port : %s",port)
    handler.Dialgrpc(conn)
    defer conn.Close()
    err = srv.ListenAndServe()
    if err != nil {
        log.Println("Eroooxxx "+err.Error())
        return
    }
}