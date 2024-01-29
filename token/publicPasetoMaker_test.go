package token

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/Streamfair/streamfair-idp-svc/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

func TestPublicPasetoMaker(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	maker, err := NewPublicPasetoMaker(privateKey, publicKey)
	require.NoError(t, err)

	username := util.RandomUsername()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreatePublicToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyPublicToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPublicToken(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	maker, err := NewPublicPasetoMaker(privateKey, publicKey)
	require.NoError(t, err)

	username := util.RandomUsername()
	duration := -time.Minute

	token, err := maker.CreatePublicToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyPublicToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPublicToken(t *testing.T) {
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	maker, err := NewPublicPasetoMaker(privateKey, publicKey)
	require.NoError(t, err)

	payload, err := maker.VerifyPublicToken("invalid_token")
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
