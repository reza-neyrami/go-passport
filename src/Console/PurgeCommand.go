package console

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Token struct {
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

type AuthCode struct {
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

type RefreshToken struct {
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

func Purge(db *sqlx.DB, hours int) error {
	// Revoke all revoked tokens
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	query := "UPDATE tokens SET revoked = true"
	_, err = tx.ExecContext(context.Background(), query)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Delete all expired tokens
	expired := time.Now().Add(-time.Duration(hours) * time.Hour)
	query = "DELETE FROM tokens WHERE expires_at < $1"
	_, err = db.ExecContext(context.Background(), query, expired)
	if err != nil {
		return err
	}

	// Revoke all revoked auth codes
	tx, err = db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	query = "UPDATE auth_codes SET revoked = true"
	_, err = tx.ExecContext(context.Background(), query)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Delete all expired auth codes
	expired = time.Now().Add(-time.Duration(hours) * time.Hour)
	query = "DELETE FROM auth_codes WHERE expires_at < $1"
	_, err = db.ExecContext(context.Background(), query, expired)
	if err != nil {
		return err
	}

	// Revoke all revoked refresh tokens
	tx, err = db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	query = "UPDATE refresh_tokens SET revoked = true"
	_, err = tx.ExecContext(context.Background(), query)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Delete all expired refresh tokens
	expired = time.Now().Add(-time.Duration(hours) * time.Hour)
	query = "DELETE FROM refresh_tokens WHERE expires_at < $1"
	_, err = db.ExecContext(context.Background(), query, expired)
	if err != nil {
		return err
	}

	return nil
}
