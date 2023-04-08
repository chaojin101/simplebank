package token

import (
	"testing"
	"time"

	"github.com/chaojin101/simplebank/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(minSecretKeySize))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.Equal(t, payload.IssuedAt.Add(duration), payload.ExpiredAt)
}

func TestExpiredJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(minSecretKeySize))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(minSecretKeySize))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
