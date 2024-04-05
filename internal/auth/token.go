package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	TokenRequest struct {
		Email string
	}
	TokenResponse struct {
		Token string `json:"access_token"`
	}
)

var (
	DefaultTokenExpiration   = time.Now().Add(time.Minute * 10).Unix()
	ErrTokenGenerationFailed = errors.New("failed to generate JWT token")
)

const TOKEN_SECRET = "OMS_SECRET_FOR_TOKEN"

func GenerateAccessToken(req TokenRequest) (*TokenResponse, error) {
	claims := jwt.MapClaims{
		"username": req.Email,
		"exp":      DefaultTokenExpiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(TOKEN_SECRET))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return nil, ErrTokenGenerationFailed
	}

	return &TokenResponse{Token: tokenString}, nil
}
