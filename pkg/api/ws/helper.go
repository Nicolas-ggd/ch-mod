package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
)

var (
	secret = os.Getenv("PRIVATE_KEY")
)

// ExtractBearerToken takes JWTToken as a string parameter and returns token
func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("invalid credentials, missing header, or bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	return jwtToken[0], nil
}

// ParseJWTClaims function takes JWTToken as a string parameter and returns extract jwt payload
func ParseJWTClaims(header string) (*models.TokenClaim, error) {
	token, err := ExtractBearerToken(header)
	if err != nil {
		return nil, fmt.Errorf("failed to extract bearer token: %v", err)
	}

	parsedToken, err := ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate JWT token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected type for JWT claims")
	}

	var tokenClaim models.TokenClaim

	marshalled, err := json.Marshal(claims["user"])
	if err != nil {
		return nil, fmt.Errorf("marshalling error")
	}

	err = json.Unmarshal(marshalled, &tokenClaim)
	if err != nil {
		return nil, fmt.Errorf("error converting 'user' claim to float64")
	}

	return &tokenClaim, nil
}

// ValidateJWT function takes JWTToken as a string parameter and checks if it's valid generated token
func ValidateJWT(jwtToken string) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}
