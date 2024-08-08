package token

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/ed25519"
)

type PasetoPayload struct {
	Email              string    `json:"email"`
	IssueDateTime      time.Time `json:"issue_datetime"`
	ExpirationDateTime time.Time `json:"expiration_datetime"`
}

func (pp *PasetoPayload) isExpired() bool {
	return time.Now().After(pp.ExpirationDateTime)
}

type TokenManger interface {
	CreateToken(email string, duration time.Duration) (string, error)
	ValidateToken(token string) (PasetoPayload, error)
}

type PasetoTokenManager struct {
	paseto     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func (ptm *PasetoTokenManager) CreateToken(email string, duration time.Duration) (string, error) {
	return ptm.paseto.Sign(ptm.privateKey, PasetoPayload{
		Email:              email,
		IssueDateTime:      time.Now(),
		ExpirationDateTime: time.Now().Add(duration),
	}, nil)
}

func (ptm *PasetoTokenManager) ValidateToken(token string) (PasetoPayload, error) {
	pp := PasetoPayload{}
	if err := ptm.paseto.Verify(token, ptm.publicKey, &pp, nil); err != nil {
		return pp, err
	}

	if pp.isExpired() {
		return pp, fmt.Errorf("paseto token is expired")
	}

	return pp, nil
}

func newPasetoTokenManager(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) *PasetoTokenManager {
	return &PasetoTokenManager{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func BuildTokenManager() (*PasetoTokenManager, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return newPasetoTokenManager(publicKey, privateKey), nil
}
