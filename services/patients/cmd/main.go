package main

import (
	"log"
	"net/http"

	handler "github.com/bertoxic/med/services/patient-service/internal/handlers"
	routes "github.com/bertoxic/med/services/patient-service/internal/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    r := routes.Router()
	srv := &http.Server{
        Addr: "0.0.0.0:8086",
        Handler: r,
    }
    conn, err := grpc.NewClient(":5001",grpc.WithTransportCredentials(insecure.NewCredentials()))
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