package bridge

import (
    "database/sql"
    "fmt"
    "time"
)

type AccessTokenRepository struct {
    db *sql.DB
}

func NewAccessTokenRepository(db *sql.DB) *AccessTokenRepository {
    return &AccessTokenRepository{
        db: db,
    }
}

func (r *AccessTokenRepository) getNewAccessToken(clientID string, scopes []string, userID string) (*AccessToken, error) {
    token := &AccessToken{
        ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
        ClientID:     clientID,
        Scopes:      scopes,
        UserID:      userID,
        Revoked:     false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        ExpiresAt:   nil, // ExpiresAt will be set later
    }

    return token, nil
}

func (r *AccessTokenRepository) persistNewAccessToken(accessToken *AccessToken) error {
    stmt, err := r.db.Prepare("INSERT INTO access_tokens (id, user_id, client_id, scopes, revoked, created_at, updated_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(accessToken.ID, accessToken.UserID, accessToken.ClientID, accessToken.ScopesString(), accessToken.Revoked, accessToken.CreatedAt, accessToken.UpdatedAt, accessToken.ExpiresAt)
    if err != nil {
        return err
    }

    return nil
}

func (r *AccessTokenRepository) revokeAccessToken(tokenID string) error {
    _, err := r.db.Exec("UPDATE access_tokens SET revoked = true WHERE id = ?", tokenID)
    if err != nil {
        return err
    }

    return nil
}

func (r *AccessTokenRepository) isAccessTokenRevoked(tokenID string) (bool, error) {
    var revoked bool
    err := r.db.QueryRow("SELECT revoked FROM access_tokens WHERE id = ?", tokenID).Scan(&revoked)
    if err != nil {
        return false, err
    }

    return revoked, nil
}
