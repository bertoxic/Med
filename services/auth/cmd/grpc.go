package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	config "github.com/bertoxic/med/services/authentication/configs"
	"github.com/bertoxic/med/services/authentication/grpc"
	handler "github.com/bertoxic/med/services/authentication/internal/handlers"
	"github.com/bertoxic/med/services/authentication/internal/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	gp "google.golang.org/grpc"
)

type UserServer struct {
	grpc.UnimplementedUserAuthServiceServer
}

var db *mongo.Client
var conf *config.AppConfig
var dataMed *mongo.Database

type DataB struct {
	app *config.AppConfig
}



func NewDataB(appconfig *config.AppConfig) {
	db = appconfig.Client
	conf = appconfig
	dataMed = db.Database(appconfig.Config.DBNAME)
}

func (s *UserServer) RegisterUser(ctx context.Context, req *grpc.JsonRequest) (*grpc.JsonResponse, error) {
	log.Println("in reg user now")
	user := &models.UserDetails{}
	jsonreq := req.Data
	err := json.Unmarshal([]byte(jsonreq), user)	
	if err != nil {
		return nil, err
	}
	log.Println(user)
	patientCollection := dataMed.Collection("patients")
	ctxs, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	emailExist, err :=  emailExists(user.Email, dataMed)
	if err != nil {
		error := models.ErrorJson{
			Code: 500,
			Message: "something went wrong during signup",
		}
			errstr , err := json.Marshal(error)
			if err != nil {
				log.Println(err)
			}
				return &grpc.JsonResponse{
			Success: false,
			Message: "registation unsuccessful"+err.Error(),
			Error: string(errstr),
		}, err
	}
	if emailExist{
			var errorJson = models.ErrorJson{
				Code: http.StatusBadRequest,
				Message:errors.New("this email has already been registered please signin").Error() ,
			}
			jsondat, err := json.MarshalIndent(errorJson,""," ")
			if err != nil {
				log.Println(err)
			}
			error := string(jsondat)
		
			return &grpc.JsonResponse{
				Success: false,
				Message: "registation unsuccessful",
				Error: error ,
			}, nil
	}
	log.Println("trying too register")
	bsoncmd := bson.M{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":	 user.Email,
		"password":   user.PassWord,
		"role":       user.UserType,
	}
	result, err := patientCollection.InsertOne(ctxs, bsoncmd)
	if err != nil {
		return  &grpc.JsonResponse{
			Success: false,
			Message: "registation successful",
			Error: err.Error(),
		}, nil
	}
	insertedId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to assert type of insertedIdzzzzzzzz")
	}
	insertedIdStr := insertedId.Hex()
	data := map[string]interface{}{
		"id": insertedIdStr,
	}
	jsondat, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return  &grpc.JsonResponse{
			Success: false,
			Message: "registation successful",
			Data:    string(jsondat),
			Error: err.Error(),
		}, nil
	}
	return &grpc.JsonResponse{
		Success: true,
		Message: "registation successful",
		Data:    string(jsondat),
	}, nil
}


func (us *UserServer) LoginUser(ctx context.Context, req *grpc.JsonRequest) (*grpc.JsonResponse, error) {
	log.Println("in login user now")

	 user := models.UserDetails{}
	 Jsondata := &models.JsonResponse{}
	if !req.Success {
		log.Println("aaaa",req.Data)
		return &grpc.JsonResponse{
			Success: false,
			Error:   "error occured:" + req.Error,
		}, errors.New("invalid values obained from jsonrequest")
	}
	err := json.Unmarshal([]byte(req.Data),Jsondata)
	if err != nil {
		return nil, err
	}
	
	b, err := json.Marshal(Jsondata.Data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b,&user)
	if err != nil {
		return nil, err
	}
	emailExist, err :=  emailExists(user.Email, dataMed)
	if err != nil {
		error := models.ErrorJson{
			Code: 500,
			Message: "something went wrong during signup",
		}
			errstr , err := json.Marshal(error)
			if err != nil {
				log.Println(err)
			}
				return &grpc.JsonResponse{
			Success: false,
			Message: "login unsuccessful"+err.Error(),
			Error: string(errstr),
		}, err
	}
	if !emailExist{
		var errorJson = models.ErrorJson{
			Code: http.StatusBadRequest,
			Message:errors.New("this email is not registered, please register").Error(),

		}
		jsondat, err := json.MarshalIndent(errorJson,""," ")
		if err != nil {
			log.Println(err)
		}
		error := string(jsondat)
	
		return &grpc.JsonResponse{
			Success: false,
			Message: "login unsuccessful",
			Error: error ,
		}, nil
}

	
	patientCollection := dataMed.Collection("patients")
	ctxs, cancel := context.WithTimeout(ctx, time.Second * 20)
	defer cancel()
	bsonCmd := bson.M{
		"email":user.Email,
		"password": user.PassWord,
	
}

	userdetail := models.UserDetails{}
	cursor := patientCollection.FindOne(ctxs,bsonCmd)
	cursor.Decode(&userdetail)
	if user.Email !=userdetail.Email && user.PassWord != userdetail.PassWord{
		var error = models.ErrorJson{
			Code: http.StatusBadRequest,
			Message: "email or password is incorrect",
		}
		var jsondata = models.JsonResponse{
			Success: false,
			Message: "",
			Error: &error,

		}
		jsondat, err := json.Marshal(jsondata)
		if err != nil {
			log.Println(err)
		}
		data := string(jsondat)
		
		return &grpc.JsonResponse{
			Success: false,
			Message: "login unsuccessful",
			Data: data,
		}, nil
	}
	user_patient := models.Patient{}
	if userdetail.UserType == "patient"{
		
		cursor.Decode(&user_patient)
	}
	token, err := handler.GenerateToken(userdetail)
	
	if err != nil {
		 	
		return &grpc.JsonResponse{
			Success: false,
			Message: "unable to generate token",
		}, nil
	}
	// datax := struct{token handler.Tokens; user models.UserDetails;}{
	// 	token: token,
	// 	user: user ,
	// }
	data := map[string]interface{}{
		"token" : token,
		"user" : user_patient,
	}
	var jsondata = models.JsonResponse{
		Success: true,
		Message: "",
		Data: data,
	}
	
	jsondat, err := json.Marshal(jsondata)
		if err != nil {
			log.Println(err)
		}
		dat := string(jsondat)
	return &grpc.JsonResponse{
		Success: true,
		Message: "login successful",
		Data: dat,
	}, nil

}


func emailExists (email string, dataB *mongo.Database) (bool, error){
	user := models.UserDetails{}
	patientCollection:=  dataB.Collection("patients")
	ctxs, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"email":email}
	cursor := patientCollection.FindOne(ctxs, filter)
	
	err := cursor.Decode(&user) 
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GenerateTokenViaGRPC(userDetails models.UserDetails) (models.Tokens, error) {
	var tokens models.Tokens
	claims := &models.SignedDetails{
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Email:     userDetails.Email,
		UserType:  userDetails.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 25).Unix(),
		},
	}
	refreshClaims := &models.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 196).Unix(),
		},
	}
	secret_key := []byte("bert")
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.StandardClaims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshClaims.StandardClaims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}
	tokens = models.Tokens{
		Token:        token,
		RefreshToken: refreshtoken,
	}

	return tokens, nil
}

func ValidateTokenViaGRPC(signedToken string) (*models.SignedDetails, string) {
	var userclaims models.SignedDetails

	token, err := jwt.ParseWithClaims(signedToken, &userclaims.StandardClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte("bert"), nil
	})
	if err != nil {

		return nil, "cannot parse token"
	}
	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok {
		return nil, "invalid token"
	}
	if claims.StandardClaims.ExpiresAt < time.Now().Local().Unix() {
		return nil, "expired token"
	}
	return claims, ""
}



func grpcListen() {
    port := os.Getenv("GRPC_PORT")
    if port == "" {
        port = "10000" // Default to 10000 if GRPC_PORT is not set
    }

    lis, err := net.Listen("tcp4", "0.0.0.0:"+port)
    if err != nil {
        log.Fatalf("failed to listen on port %v: %v", port, err)
    }

    s := gp.NewServer(
        gp.MaxRecvMsgSize(1024*1024*10), // 10MB
        gp.MaxSendMsgSize(1024*1024*10), // 10MB
    )
    log.Printf("running on port..>> %s", port)
    grpc.RegisterUserAuthServiceServer(s, &UserServer{})
    log.Printf("server listening at %v on port %v", lis.Addr(), port)

    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// func grpcListen(){
// 	port := os.Getenv("GRPC_PORT")	
//     if port == "" {
//         port = "10000" // Default to 10000 if PORT is not set
//     }

//     lis, err := net.Listen("tcp", "0.0.0.0:"+port)
//     if err != nil {
//         log.Fatalf("failed to listen to: %v, error is : %v", port,err)
//     }
// 	//lis, err := net.Listen("tcp", ":5001")
// 	// if err != nil {
// 	// 	log.Printf("did not listen failed dto isten: %v", err)
// 	// }
// 	s := gp.NewServer()
// 	log.Printf("running on port..>> %s",port)
// 	grpc.RegisterUserAuthServiceServer(s, &UserServer{})
// 	log.Printf("server listening at %v on port %v", lis.Addr(),port)
// 	if err := s.Serve(lis); err != nil {
// 		log.Printf("failed to serve: %v", err)
// 	}

// }