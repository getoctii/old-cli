package utils

import (
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/havoc-io/go-keytar"
)


func Authorization() (string, error) {
	keychain, err := keytar.GetKeychain()

	if err != nil {
		return "", err
	}

	authorization, err := keychain.GetPassword("octii", "authorization")

	if err != nil {
		return "", err
	}

	return authorization, nil
}

func ParseAuthorization(token string) (jwt.StandardClaims, error) {
	var claims = jwt.StandardClaims{}

	_, _, err := jwt.NewParser().ParseUnverified(token, &claims)

	if err != nil {
		return jwt.StandardClaims{}, err
	}

	return claims, nil
}