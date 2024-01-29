package bridge

import "time"

type Client struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	IsConfidential bool      `json:"is_confidential"`
	RedirectURIs   []string  `json:"redirect_uris"`
	Provider       string    `json:"provider"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewClient(identifier string, name string, redirectURI string, isConfidential bool, provider string) *Client {
	return &Client{
		ID:             identifier,
		Name:           name,
		IsConfidential: isConfidential,
		RedirectURIs:   []string{redirectURI},
		Provider:  time provider,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (client *Client) GetIdentifier() string {
	return client.ID
}
