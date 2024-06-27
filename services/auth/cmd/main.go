package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	config "github.com/bertoxic/med/services/authentication/configs"
	routes "github.com/bertoxic/med/services/authentication/internal/transport/http"
	"github.com/bertoxic/med/services/authentication/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	app, err := config.NewAppConfig()
	if err != nil {
		log.Println("error occured while setting env", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	client := database.NewMongodClient(ctx, app.Config.MONGODB_URI)
	app.InProduction = false
	app.Client = client
	NewDataB(app)
	grpcListen()
	routes := routes.Router()
	port := os.Getenv("PORT")
	portnum, err := strconv.Atoi(port)
	if err != nil {
		log.Println("error getting or converting port")
	}
	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", portnum),
		Handler: routes,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println("Eroooxxx " + err.Error())
		return
	}
}
