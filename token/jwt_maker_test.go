package token

import (
	"lesson/simple-bank/utils"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	userName := utils.RandomOwner()
	duration := time.Minute

	issuedTime := time.Now()
	expiredTime := issuedTime.Add(duration)

	token, payload, err := maker.CreateToken(userName, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	require.WithinDuration(t, payload.ExpiresAt, expiredTime, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedTime, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrExpiredToken)
	require.Nil(t, payload)
}

func TestInvalidJWT(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
