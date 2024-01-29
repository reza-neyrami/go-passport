package bridge

import (
    "fmt"
    "time"
)

type AccessToken struct {
    ID          string   `json:"id"`
    UserID      string   `json:"user_id"`
    ClientID     string   `json:"client_id"`
    Scopes      []string `json:"scopes"`
    Revoked     bool     `json:"revoked"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    ExpiresAt   *time.Time `json:"expires_at"`
    RefreshToken *string   `json:"refresh_token"`
}

func NewAccessToken(userID string, scopes []string, clientID string) *AccessToken {
    return &AccessToken{
        ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
        UserID:      userID,
        ClientID:     clientID,
        Scopes:      scopes,
        Revoked:     false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        ExpiresAt:   nil, // ExpiresAt will be set later
        RefreshToken: nil,
    }
}

func (accessToken *AccessToken) SetExpiresAt(expiresAt time.Time) {
    accessToken.ExpiresAt = &expiresAt
}

func (accessToken *AccessToken) HasExpired() bool {
    if accessToken.ExpiresAt == nil {
        return false
    }

    return time.Now().After(*accessToken.ExpiresAt)
}
