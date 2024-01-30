package main

import (
    "errors"
    "fmt"
    "reflect"
)

type RefreshTokenRepository struct {
    model *RefreshToken
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
    return &RefreshTokenRepository{
        model: &RefreshToken{},
    }
}

// Create creates a new refresh token.
func (r *RefreshTokenRepository) Create(attributes map[string]interface{}) (*RefreshToken, error) {
    token, err := r.model.Create(attributes)
    if err != nil {
        return nil, fmt.Errorf("Error creating refresh token: %v", err)
    }

    return token, nil
}

// Find gets a refresh token by the given ID.
func (r *RefreshTokenRepository) Find(id string) (*RefreshToken, error) {
    token := r.model.FindByID(id)
    if token == nil {
        return nil, errors.New("Refresh token not found")
    }

    return token, nil
}

// Save stores the given token instance.
func (r *RefreshTokenRepository) Save(token *RefreshToken) error {
    if err := token.Save(); err != nil {
        return fmt.Errorf("Error saving refresh token: %v", err)
    }

    return nil
}

// RevokeRefreshToken revokes the refresh token.
func (r *RefreshTokenRepository) RevokeRefreshToken(id string) error {
    err := r.model.Where("id", id).Update("revoked", true)
    if err != nil {
        return fmt.Errorf("Error revoking refresh token: %v", err)
    }

    return nil
}

// RevokeRefreshTokensByAccessTokenId revokes refresh tokens by access token ID.
func (r *RefreshTokenRepository) RevokeRefreshTokensByAccessTokenId(tokenId string) error {
    err := r.model.Where("access_token_id", tokenId).Update("revoked", true)
    if err != nil {
        return fmt.Errorf("Error revoking refresh tokens: %v", err)
    }

    return nil
}

// IsRefreshTokenRevoked checks if the refresh token has been revoked.
func (r *RefreshTokenRepository) IsRefreshTokenRevoked(id string) bool {
    token, err := r.Find(id)
    if err != nil {
        return true
    }

    return token.Revoked
}
