package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model
	ID          string
	TOKEN       string
	ACCESSTOKEN string
	REVOKED     bool
	EXPIRES_AT  time.Time
	USERID      uint
	CLIENTID    uint
    Scopes      []string
}

func NewToken(token, access_token string, revoked bool, expires_at time.Time, user_id uint, client_id uint) *Token {
	return &Token{
		ID:          uuid.New().String(),
		TOKEN:       token,
		ACCESSTOKEN: access_token,
		REVOKED:     revoked,
		EXPIRES_AT:  expires_at,
		USERID:      user_id,
		CLIENTID:    client_id,
	}
}

func (t *Token) Revoke() error {
	t.REVOKED = true

	err := t.Save(t).Error
	if err != nil {
		return fmt.Errorf("Error revoking token: %v", err)
	}

	return nil
}

func (t *Token) HasScope(scope string) bool {
	for _, s := range t.Scopes {
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
