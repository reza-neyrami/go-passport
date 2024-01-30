package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	id           string
	token        string
	access_token string
	revoked      bool
	expires_at   time.Time
	user_id      string
	client_id    string
}

func NewToken(token, access_token string, revoked bool, expires_at time.Time, user_id, client_id string) *Token {
	return &Token{
		id:           uuid.New().String(),
		token:        token,
		access_token: access_token,
		revoked:      revoked,
		expires_at:   expires_at,
		user_id:      user_id,
		client_id:    client_id,
	}
}

func (t *Token) Revoke() error {
	t.revoked = true

	err := t.Save()
	if err != nil {
		return fmt.Errorf("Error revoking token: %v", err)
	}

	return nil
}

func (t *Token) HasScope(scope string) bool {
	for _, s := range t.scopes {
		if s == scope {
			return true
		}
	}

	return false
}

func (t *Token) Can(scope string) bool {
	return t.HasScope(scope)
}

func (t *Token) Cant(scope string) bool {
	return !t.Can(scope)
}

func (t *Token) transient() bool {
	return false
}
