package jwtauth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTSession struct {
	SecretKey []byte
}

func NewJWTSession(secretKey []byte) *JWTSession {

	return &JWTSession{SecretKey: secretKey}

}

func CreateJWTToken(email string, secretKey []byte) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,                            // Subject (user identifier)
		"iss": "icealpha",                       // Issuer
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
