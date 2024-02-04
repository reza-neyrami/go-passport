package src

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthCode struct {
	ID          string     `db:"id"`
	Code        string     `db:"code"`
	ClientID    string     `db:"client_id"`
	UserID      int64      `db:"user_id"`
	RedirectURI string     `db:"redirect_uri"`
	Scopes      []string   `db:"scopes"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	RevokedAt   *time.Time `db:"revoked_at"`
}

func NewAuthCode(code string, clientID string, userID int64, redirectURI string, scopes []string) *AuthCode {
	return &AuthCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scopes:      scopes,
	}
}

func (c *AuthCode) Create(db *sqlx.DB) error {
	ctx := context.WithValue(context.Background(), "scopes", c.Scopes)
	ctx = context.WithValue(ctx, "clientID", c.ClientID)
	stmt, err := db.PreparexContext(ctx, "INSERT INTO oauth_auth_codes (code, client_id, user_id, redirect_uri, scopes) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, c.Code, c.ClientID, c.UserID, c.RedirectURI, c.Scopes)
	if err != nil {
		return err
	}

	return nil
}

func FindAuthCode(db *sqlx.DB, code string) (*AuthCode, error) {
	ctx := context.WithValue(context.Background(), "code", code)
	row := db.QueryRowxContext(ctx, "SELECT id, code, client_id, user_id, redirect_uri, scopes, created_at, updated_at, revoked_at FROM oauth_auth_codes WHERE code = $1", code)

	var authCode AuthCode
	err := row.Scan(&authCode.ID, &authCode.Code, &authCode.ClientID, &authCode.UserID, &authCode.RedirectURI, &authCode.Scopes, &authCode.CreatedAt, &authCode.UpdatedAt, &authCode.RevokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &authCode, nil
}

func RevokeAuthCode(db *sqlx.DB, code string) error {
	ctx := context.WithValue(context.Background(), "code", code)
	stmt, err := db.PreparexContext(ctx, "UPDATE oauth_auth_codes SET revoked_at = NOW() WHERE code = $1")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, code)
	if err != nil {
		return err
	}

	return nil
}
