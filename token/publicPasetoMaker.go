package token

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	"github.com/o1egl/paseto"
	"github.com/pkg/errors"
)

type PublicPasetoMaker struct {
	paseto     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewPublicPasetoMaker(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (PublicMaker, error) {
	maker := &PublicPasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return maker, nil
}

func (maker *PublicPasetoMaker) CreatePublicToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal payload")
	}

	token, err := maker.paseto.Sign(maker.privateKey, jsonPayload, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign payload")
	}

	return token, nil
}

func (maker *PublicPasetoMaker) VerifyPublicToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
