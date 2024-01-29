package bridge

import (
	"fmt"
	"time"
)

type RefreshToken struct {
	ID        string     `json:"id"`
	ClientID  string     `json:"client_id"`
	UserID    string     `json:"user_id"`
	Token     *string    `json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	Scopes    []string   `json:"scopes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewRefreshToken(
	clientID,
	userID,
	token string,
	expiresAt time.Time,
	scopes []string) *RefreshToken {

	return &RefreshToken{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		ClientID:  clientID,
		UserID:    userID,
		Token:     &token,
		ExpiresAt: expiresAt,
		RevokedAt: nil,
		Scopes:    scopes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (token *RefreshToken) GetIdentifier() string {
	return token.ID
}

func (token *RefreshToken) SetIdentifier(identifier string) {
	token.ID = identifier
}
