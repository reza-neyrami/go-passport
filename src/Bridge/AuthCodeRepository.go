package bridge

import (
	"database/sql"
	"fmt"
	"time"
)

type AuthCodeRepository struct {
	db *sql.DB
}

func NewAuthCodeRepository(db *sql.DB) *AuthCodeRepository {
	return &AuthCodeRepository{
		db: db,
	}
}

func (r *AuthCodeRepository) getNewAuthCode() (*AuthCode, error) {
	code := &AuthCode{}

	return code, nil
}

func (r *AuthCodeRepository) persistNewAuthCode(authCode *AuthCode) error {
	stmt, err := r.db.Prepare("INSERT INTO auth_codes (id, user_id, client_id, scopes, revoked, expires_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(authCode.ID, authCode.UserID, authCode.ClientID, authCode.ScopesString(), authCode.Revoked, authCode.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthCodeRepository) revokeAuthCode(codeID string) error {
	_, err := r.db.Exec("UPDATE auth_codes SET revoked = true WHERE id = ?", codeID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthCodeRepository) isAuthCodeRevoked(codeID string) (bool, error) {
	var revoked bool
	err := r.db.QueryRow("SELECT revoked FROM auth_codes WHERE id = ?", codeID).Scan(&revoked)
	if err != nil {
		return false, err
	}

	return revoked, nil
}
