package bridge

import (
    "database/sql"
    "fmt"
)

type RefreshTokenRepository struct {
    db *sql.DB
    eventDispatcher *events.Dispatcher
}

func NewRefreshTokenRepository(db *sql.DB, eventDispatcher *events.Dispatcher) *RefreshTokenRepository {
    return &RefreshTokenRepository{
        db: db,
        eventDispatcher: eventDispatcher,
    }
}

func (r *RefreshTokenRepository) GetNewRefreshToken() *RefreshToken {
    return &RefreshToken{}
}

func (r *RefreshTokenRepository) PersistNewRefreshToken(token *RefreshToken) error {
    stmt, err := r.db.Prepare("INSERT INTO refresh_tokens (id, access_token_id, revoked, expires_at) VALUES (?, ?, ?, ?)")
    if err != nil {
        return err
    }

    defer stmt.Close()

    if _, err := stmt.Exec(token.ID, token.AccessTokenID, false, token.ExpiresAt); err != nil {
        return err
    }

    r.eventDispatcher.Dispatch(NewRefreshTokenCreated(token.ID, token.AccessTokenID))
    return nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(tokenID string) error {
    _, err := r.db.Exec("UPDATE refresh_tokens SET revoked = true WHERE id = ?", tokenID)
    return err
}

func (r *RefreshTokenRepository) IsRefreshTokenRevoked(tokenID string) bool {
    var revoked bool
    err := r.db.QueryRow("SELECT revoked FROM refresh_tokens WHERE id = ?", tokenID).Scan(&revoked)
    if err != nil {
        return false
    }
    return revoked
}
