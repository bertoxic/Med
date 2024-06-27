package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
}

type SignedDetails struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	UserType       string `json:"user_type"`
	StandardClaims jwt.StandardClaims
}

type Tokens struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func generateToken(userDetails UserDetails) (Tokens, error) {
	claims := &SignedDetails{
		FirstName: userDetails.FirstName,
		LastName:  userDetails.LastName,
		Email:     userDetails.Email,
		UserType:  userDetails.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}
	var tokens Tokens
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * 192).Unix()},
	}
	secret_key := []byte("bert")

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret_key)
	if err != nil {
		return tokens, err
	}
	tokens.Token = token
	tokens.RefreshToken = refreshToken
	return tokens, nil

}

func VerifyToken(signedToken string) (*SignedDetails, string) {
	var userclaims SignedDetails

	token, err := jwt.ParseWithClaims(signedToken, &userclaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("bert"), nil // Assuming this is the correct key
	})
	if err != nil {

		return nil, "cannot parse token" + err.Error()
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, "invalid token"
	}
	if claims.StandardClaims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Sprintf("expiredToken, %d>> %d", claims.StandardClaims.ExpiresAt, time.Now().Local().Unix())
	}

	return claims, ""
}

func (c *SignedDetails) Valid() error {
	return c.StandardClaims.Valid()
}

func RefreshTokenEndpoint(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the refresh token
	claims, msg := VerifyToken(requestBody.RefreshToken)
	if msg != "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}
	var user UserDetails
	user.Email = claims.Email
	user.FirstName = claims.FirstName
	user.LastName = claims.LastName
	user.UserType = "patient"
	// Generate new tokens
	tokens, err := GenerateTokenViaRestAPI(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.Token,
		"refresh_token": tokens.RefreshToken,
	})
}
