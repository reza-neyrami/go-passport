package main

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Client struct {
	ID                   int      `db:"id"`
	UserId               uint      `db:"user_id"`
	Name                 string   `db:"name"`
	Secret               string   `db:"secret"`
	Provider             string   `db:"provider"`
	Redirect             string   `db:"redirect"`
	PersonalAccessClient bool     `db:"personal_access_client"`
	PasswordClient       bool     `db:"password_client"`
	Revoked              bool     `db:"revoked"`
	GrantTypes           []string `db:"grant_types"`
	Scopes               []string `db:"scopes"`
	db                   *sqlx.DB
}

func (c *Client) User() User {
	// Implement the logic to get the user associated with the client here.
	return User{}
}

func (c *Client) AuthCodes() []AuthCode {
	// Implement the logic to get the auth codes associated with the client here.
	return []AuthCode{}
}

func (c *Client) Tokens() []Token {
	// Implement the logic to get the tokens associated with the client here.
	return []Token{}
}

func (c *Client) GetGrantTypesAttribute() []string {
	return c.GrantTypes
}

func (c *Client) GetScopesAttribute() []string {
	return c.Scopes
}

func (c *Client) SetScopesAttribute(scopes []string) {
	c.Scopes = scopes
}

func (c *Client) GetPlainSecretAttribute() string {
	return c.Secret
}

func (c *Client) SetSecretAttribute(value string) {
	c.Secret = value
}

func (c *Client) FirstParty() bool {
	return c.PersonalAccessClient || c.PasswordClient
}

func (c *Client) SkipsAuthorization() bool {
	return false
}

func (c *Client) HasScope(scope string) bool {
	if c.Scopes == nil {
		return true
	}

	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}

func (c *Client) Confidential() bool {
	return c.Secret != ""
}

func (c *Client) GetKeyType() string {
	return "string"
}

func (c *Client) GetIncrementing() bool {
	return false
}



func (c *Client) FindByID(ctx context.Context) (*Client, error) {
	var client Client
	err := c.db.GetContext(ctx, &client, `
		SELECT
			id,
			name,
			secret,
			provider,
			redirect,
			personal_access_client,
			password_client,
			revoked,
			grant_types,
			scopes
		FROM clients
		WHERE id = $1
	`,
		c.ID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &client, nil
}

func (c *Client) FindByClientID(ctx context.Context) (*Client, error) {
	var client Client
	err := c.db.GetContext(ctx, &client, `
		SELECT
			id,
			name,
			secret,
			provider,
			redirect,
			personal_access_client,
			password_client,
			revoked,
			grant_types,
			scopes
		FROM clients
		WHERE client_id = $1
	`,
		c.ClientID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &client, nil
}

func (c *Client) FindByName(ctx context.Context) (*Client, error) {
	var client Client
	err := c.db.GetContext(ctx, &client, `
		SELECT
			id,
			name,
			secret,
			provider,
			redirect,
			personal_access_client,
			password_client,
			revoked,
			grant_types,
			scopes
		FROM clients
		WHERE name = $1
	`,
		c.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &client, nil
}

func (c *Client) Create(ctx context.Context) error {
	// Define a transaction object
	tx, err := c.db.Beginx(ctx)
	if err != nil {
		return err
	}

	// Insert the client record
	result, err := tx.ExecContext(ctx, `
		INSERT INTO clients (
			name,
			secret,
			provider,
			redirect,
			personal_access_client,
			password_client,
			revoked,
			grant_types,
			scopes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`,
		c.Name,
		c.Secret,
		c.Provider,
		c.Redirect,
		c.PersonalAccessClient,
		c.PasswordClient,
		c.Revoked,
		c.GrantTypes,
		c.Scopes,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Retrieve the client ID from the result
	clientID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the client ID if it's not empty
	if c.ID != 0 {
		_, err := tx.ExecContext(ctx, `
			UPDATE clients
			SET id = $1
			WHERE id = $2
		`,
			clientID,
			c.ID,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	c.ID = clientID
	return nil
}

func (c *Client) Update(ctx context.Context) error {
	// Define a transaction object
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	// Update the client record
	_, err = tx.ExecContext(ctx, `
		UPDATE clients
		SET
			name = $1,
			secret = $2,
			provider = $3,
			redirect = $4,
			personal_access_client = $5,
			password_client = $6,
			revoked = $7,
			grant_types = $8,
			scopes = $9
		WHERE id = $10
	`,
		c.Name,
		c.Secret,
		c.Provider,
		c.Redirect,
		c.PersonalAccessClient,
		c.PasswordClient,
		c.Revoked,
		c.GrantTypes,
		c.Scopes,
		c.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}


func (c *Client) Delete(ctx context.Context) error {
	// Define a transaction object
	tx, err := c.db.Beginx()
	if err != nil {
		return err
	}

	// Delete the client record
	_, err = tx.ExecContext(ctx, `
		DELETE FROM clients
		WHERE id = $1
	`,
		c.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

