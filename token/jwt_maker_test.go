package token

import (
	"testing"
	"time"

<<<<<<< HEAD
	"github.com/dgrijalva/jwt-go"
=======
	"github.com/golang-jwt/jwt/v5"
>>>>>>> d4d0e58 (refactor)
	"github.com/spaghetti-lover/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
<<<<<<< HEAD
	duration := time.Minute

	//issuedAt := time.Now()
	//expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
=======
	role := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration, TokenTypeAccessToken)
>>>>>>> d4d0e58 (refactor)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

<<<<<<< HEAD
	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotZero(t, payload.ID)
	require.NoError(t, err)

	token, payload, err = maker.CreateToken(util.RandomOwner(), -time.Minute)
=======
	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, -time.Minute, TokenTypeAccessToken)
>>>>>>> d4d0e58 (refactor)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

<<<<<<< HEAD
	payload, err = maker.VerifyToken(token)
=======
	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
>>>>>>> d4d0e58 (refactor)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
<<<<<<< HEAD
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
=======
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute, TokenTypeAccessToken)
>>>>>>> d4d0e58 (refactor)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
<<<<<<< HEAD
	payload, err = maker.VerifyToken(token)
=======

	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestJWTWrongTokenType(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, time.Minute, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeRefreshToken)
>>>>>>> d4d0e58 (refactor)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
