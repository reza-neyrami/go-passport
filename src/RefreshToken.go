package src

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
    *gorm.Model
    id          string
    token       string
    access_token string
    revoked     bool
    expires_at  time.Time
}

func NewRefreshToken(token string, accessToken string, revoked bool, expiresAt time.Time) *RefreshToken {
    return &RefreshToken{
        id:          uuid.New().String(),
        token:       token,
        access_token: accessToken,
        revoked:     revoked,
        expires_at:  expiresAt,
    }
}

func (r *RefreshToken) Revoke() error {
    r.revoked = true

    err := r.Save()
    if err != nil {
        return fmt.Errorf("Error revoking refresh token: %v", err)
    }

    return nil
}

func (r *RefreshToken) Transient() bool {
    return false
}
