package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtMaker struct {
	secreteKey string
}

func NewJWTMaker(secreteKey string) (Maker, error) {
	return &JwtMaker{secreteKey: secreteKey}, nil
}

func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secreteKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secreteKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
