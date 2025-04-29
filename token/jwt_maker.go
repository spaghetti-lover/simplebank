package token

import (
	"errors"
	"fmt"
	"time"

<<<<<<< HEAD
	"github.com/dgrijalva/jwt-go"
=======
	"github.com/golang-jwt/jwt/v5"
>>>>>>> d4d0e58 (refactor)
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
<<<<<<< HEAD
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
=======
func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration, tokenType)
>>>>>>> d4d0e58 (refactor)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
<<<<<<< HEAD
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
=======
func (maker *JWTMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
>>>>>>> d4d0e58 (refactor)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
<<<<<<< HEAD
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
=======

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
>>>>>>> d4d0e58 (refactor)
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

<<<<<<< HEAD
=======
	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}

>>>>>>> d4d0e58 (refactor)
	return payload, nil
}
