package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID            uuid.UUID `json:"id"`
	UserName      string    `json:"username"`
	WalletAddress string    `json:"wallet"`
	IssuedAt      time.Time `json:"iat"`
	ExpiredAt     time.Time `json:"exp"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, wallet string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:            tokenID,
		UserName:      username,
		WalletAddress: wallet,
		IssuedAt:      time.Unix(int64(time.Now().Second()), int64(time.Now().Nanosecond())),
		ExpiredAt:     time.Unix(int64(time.Now().Add(duration).Second()), int64(time.Now().Add(duration).Nanosecond())),
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
