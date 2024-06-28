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
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file",err)
	}
	app, err := config.NewAppConfig()
	if err != nil {
		log.Println("error occured while setting env", err)
	}
    port := os.Getenv("PORT")
	portnum, err := strconv.Atoi(port)


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	client := database.NewMongodClient(ctx, app.Config.MONGODB_URI)
	app.InProduction = false
	app.Client = client
	NewDataB(app)
	
	routes := routes.Router()

	if err != nil {
		log.Println("error getting or converting port")
	}
	srv := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", portnum),
		Handler: routes,
	}
   go grpcListen()
	err = srv.ListenAndServe()
	if err != nil {
        log.Println("vvvvvvvv",portnum)
		log.Println("Eroooxxx " + err.Error())
		return
	}
   

}
