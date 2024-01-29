package bridge


import (
    "database/sql"
    "fmt"
    "time"
)

type AuthCode struct {
    ID           string   `json:"id"`
    ClientID     string   `json:"client_id"`
    UserID       string   `json:"user_id"`
    RedirectURI  string   `json:"redirect_uri"`
    Scopes       []string `json:"scopes"`
    Code         string   `json:"code"`
    ExpiresAt    *time.Time `json:"expires_at"`
    Revoked      bool     `json:"revoked"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func NewAuthCode(clientID string, userID string, redirectURI string, scopes []string, code string) *AuthCode {
    return &AuthCode{
        ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
        ClientID:     clientID,
        UserID:      userID,
        RedirectURI:  redirectURI,
        Scopes:      scopes,
        Code:         code,
        ExpiresAt:   nil, // ExpiresAt will be set later
        Revoked:     false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}

func (authCode *AuthCode) SetExpiresAt(expiresAt time.Time) {
    authCode.ExpiresAt = &expiresAt
}

func (authCode *AuthCode) HasExpired() bool {
    if authCode.ExpiresAt == nil {
        return false
    }

    return time.Now().After(*authCode.ExpiresAt)
}
