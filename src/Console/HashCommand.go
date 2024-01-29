package console

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type HashCommand struct {
	db *sql.DB
	chunkSize int
	force bool
}

func NewHashCommand(db *sql.DB, chunkSize int) *HashCommand {
	return &HashCommand{db: db, chunkSize: chunkSize, force: false}
}

func (c *HashCommand) SetFlags() {
	c.force = c.flag("force").Bool("force", false, "Force the operation to run without confirmation")
}

func (c *HashCommand) Execute() error {
	if c.force || c.confirm("Are you sure you want to hash all client secrets? This cannot be undone.") {
		rows, err := c.db.Query("SELECT secret FROM clients WHERE secret IS NOT NULL")
		if err != nil {
			return fmt.Errorf("error querying for client secrets: %w", err)
		}
		defer rows.Close()

		var hashedSecrets []string
		for rows.Next() {
			var secret string
			err := rows.Scan(&secret)
			if err != nil {
				return fmt.Errorf("error scanning client secret: %w", err)
			}

			hashedSecret, err := bcrypt.HashPassword([]byte(secret), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("error hashing client secret: %w", err)
			}

			hashedSecrets = append(hashedSecrets, string(hashedSecret))

			if len(hashedSecrets) >= c.chunkSize {
				err = c.saveChunk(hashedSecrets)
				if err != nil {
					return err
				}
				hashedSecrets = []string{}
			}
		}

		// ذخیره باقی‌مانده کلیدها
		if len(hashedSecrets) > 0 {
			err = c.saveChunk(hashedSecrets)
			if err != nil {
				return err
			}
		}

		return nil
	}
	return nil
}

func (c *HashCommand) saveChunk(hashedSecrets []string) error {
	ids := []string{}
	for i, _ := range hashedSecrets {
	  id := fmt.Sprintf("%d", i)
	  ids = append(ids, id)
	}
  
	query := fmt.Sprintf("UPDATE clients SET secret = '%s' WHERE id IN (%s)",
	  strings.Join(hashedSecrets, ","),
	  strings.Join(ids, ","))
	_, err := c.db.Exec(query)
	return err
  }
  
