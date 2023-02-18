package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("invalid token")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Id:        id,
		Username:  username,
		ExpiresAt: time.Now().Add(duration),
		IssuedAt:  time.Now(),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if p.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
