package console

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Token struct {
	*gorm.Model
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

type AuthCode struct {
	*gorm.Model
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

type RefreshToken struct {
	*gorm.Model
	ID        uint
	Revoked   bool
	ExpiresAt time.Time
}

type PurgeCommand struct {
	db        *gorm.DB
	Option    func(string) bool
	OptionInt func(string) int
}

func (s *PurgeCommand) Handle() error {

	// Parse the options
	revoked := s.Option("revoked")
	expired := s.Option("expired")
	hours := s.OptionInt("hours")

	// Calculate the expiration date
	expiredDate := time.Now().Add(-time.Duration(hours) * time.Hour)

	// Purge revoked and expired tokens
	if revoked && expired {
		s.db.Where("revoked = ?", true).Or("expires_at <= ?", expiredDate).Delete(&Token{})
		s.db.Where("revoked = ?", true).Or("expires_at <= ?", expiredDate).Delete(&AuthCode{})
		s.db.Where("revoked = ?", true).Or("expires_at <= ?", expiredDate).Delete(&RefreshToken{})

		if hours > 0 {
			return errors.New("cannot purge revoked and expired tokens with hours option")
		}

		return nil
	} else if revoked {
		s.db.Where("revoked = ?", true).Delete(&Token{})
		s.db.Where("revoked = ?", true).Delete(&AuthCode{})
		s.db.Where("revoked = ?", true).Delete(&RefreshToken{})

		return nil
	} else if expired {
		s.db.Where("expires_at <= ?", expiredDate).Delete(&Token{})
		s.db.Where("expires_at <= ?", expiredDate).Delete(&AuthCode{})
		s.db.Where("expires_at <= ?", expiredDate).Delete(&RefreshToken{})

		return nil
	} else {
		return errors.New("must specify either --revoked or --expired")
	}
}
