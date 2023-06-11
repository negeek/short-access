package utils

import (
	"fmt"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	jwt.RegisteredClaims
	ID    string
	Email string
}

func CreateJwtToken(id string, email string, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		ID: id,
		Email: email,
	})

	// Create the actual JWT token
	signedString, err := token.SignedString([]byte(os.Getenv("AUTH_KEY")))

	if err != nil {
		return "", err.Error()
	}

	return signedString, nil
}

func VerifyJwt(tokenStr string) (interface{}, error) {
	// Parse and validate the JWT token
	token, err := jwt.ParseWithClaims(jwtToken, &UserClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_KEY")), nil
	})
	if err != nil {
		return  "", err.Error()
	}
	return UserClaim, nil
}


