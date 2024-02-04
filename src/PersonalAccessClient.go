package src

import (
	"log"

	"gorm.io/gorm"
)

type PersonalAccessClient struct {
	*gorm.Model
	table   string
	guarded []string
}

func (p *PersonalAccessClient) Client() (Client, error) {
	// Fix compiler error by removing call to undefined method belongsTo
	client, err := passport.ClientModel()
	if err != nil {
		log.Printf("Error retrieving client for personal access token: %v", err)
		return  err
	}

	return client, nil
}
