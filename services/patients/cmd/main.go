package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/bertoxic/med/services/patient-service/internal/handlers"
	routes "github.com/bertoxic/med/services/patient-service/internal/transport/http"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
    conn, err := grpc.NewClient(":"+port,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("grpc did not connect: %v", err)
	}
    handler.Dialgrpc(conn)
    defer conn.Close()
    err = srv.ListenAndServe()
    if err != nil {
        log.Println("Eroooxxx "+err.Error())
        return
    }
}