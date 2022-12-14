package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/samandar2605/medium_user_service/config"
)

type TokenParams struct {
	UserID   int64
	Username string
	Email    string
	UserType string
	Password string
	Duration time.Duration
}

// CreateToken creates a new token
func CreateToken(cfg *config.Config, params *TokenParams) (string, *Payload, error) {
	payload, err := NewPayload(params)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(cfg.AuthSecretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func VerifyToken(cfg *config.Config, token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println(38)
			return nil, ErrInvalidToken
		}
		return []byte(cfg.AuthSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		fmt.Println(50)
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		fmt.Println(56)
		return nil, ErrInvalidToken
	}

	return payload, nil
}
