package token

import "time"

// Maker is an interface for managing tokens
type LocalMaker interface {
	// CreateLocalToken creates a new local token for a specific username and duration
	CreateLocalToken(username string, duration time.Duration) (string, error)

	// VerifyLocalToken checks if the local token is valid or not
	VerifyLocalToken(token string) (*Payload, error)
}

type PublicMaker interface {
	// CreatePublicToken creates a new public token for a specific username and duration
	CreatePublicToken(username string, duration time.Duration) (string, error)

	// VerifyPublicToken checks if the public token is valid or not
	VerifyPublicToken(token string) (*Payload, error)
}
