package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bertoxic/med/services/patient-service/grpc"
	"github.com/bertoxic/med/services/patient-service/internal/models"
	gp "google.golang.org/grpc"
)
var conn = &gp.ClientConn{}
func Dialgrpc(mainconn *gp.ClientConn ){
	conn = mainconn
}


func RegisterViaGRPC(user *models.UserDetails) (models.JsonResponse, error) {

	// conn, err := gp.NewClient(":5001",gp.WithTransportCredentials(insecure.NewCredentials()))

	// if err != nil {
	// 	log.Printf("grpc did not connect: %v", err)
	// 	return  models.JsonResponse{Success: false,Message:"",Error:&models.ErrorJson{Code: 500,Message: err.Error()}}, err
	// }
	//defer conn.Close()
	client := grpc.NewUserAuthServiceClient(conn)
	dat, err := json.Marshal(user)
	if err != nil {
		log.Println("Error:", err)
		return models.JsonResponse{Success: false, Message: "", Error: &models.ErrorJson{Code: 500, Message: err.Error()}}, err
	}
	jsonresp, err := client.RegisterUser(context.Background(), &grpc.JsonRequest{
		Success: true,
		Message: "sending user data",
		Data:    string(dat),
	})
	if err != nil {
		log.Printf("jsonresp not obtained: %v", err)
		log.Println(jsonresp)
		return models.JsonResponse{Success: false, Message: "", Error: &models.ErrorJson{Code: 500, Message: err.Error()}}, err
	}
	if !jsonresp.Success {
		error := &models.ErrorJson{}
		err := json.Unmarshal([]byte(jsonresp.Error), error)
		if err != nil {
			log.Println(err)
		}
		return models.JsonResponse{Success: jsonresp.Success, Message: jsonresp.Message, Error: error}, fmt.Errorf(jsonresp.Error)
	}

	var data map[string]interface{}

	err = json.Unmarshal([]byte(jsonresp.Data), &data)
	if err != nil {
		log.Printf("json unmarshal fail: %v", err)
		return models.JsonResponse{Success: false, Message: "", Error: &models.ErrorJson{Code: 500, Message: err.Error()}}, err
	}
	return models.JsonResponse{Success: true, Message: "", Data: data}, nil
}

func LoginUserViaGrpc(user models.UserDetails)(models.JsonResponse, error){
	var jsonresponse = &models.JsonResponse{}
	client := grpc.NewUserAuthServiceClient(conn)
	var jsonreq = &models.JsonResponse{
		Success: true,
		Message: "sending user login details",
		Data: user ,
	}
	dat, err := json.Marshal(jsonreq)
	if err != nil {
		var error = models.ErrorJson{
			Code: 500,
			Message: err.Error(),
		}
		return models.JsonResponse{Success: false,Error:&error }, err
	}
	jsonresp, err := client.LoginUser(context.Background(),&grpc.JsonRequest{
		Success: true,
		Message: "sending user login details",
		Data:string(dat) ,
	})
	if err != nil{
		var error = models.ErrorJson{
			Code: 500,
			Message: err.Error(),
		}
		return models.JsonResponse{Success: false,Error:&error }, err
	}
	jsonresponse.Success = jsonreq.Success
	jsonresponse.Message= jsonreq.Message


	err =  json.Unmarshal([]byte(jsonresp.Data),jsonresponse)
	if err != nil {
		log.Println("bb cannot unmarshall")
		var error = models.ErrorJson{
			Code: 500,
			Message: err.Error(),
		}
		return models.JsonResponse{Success: false,Error:&error }, err
	}
	return *jsonresponse, nil
}