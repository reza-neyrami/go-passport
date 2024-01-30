package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TokenRepository interface {
	Create(attributes map[string]interface{}) (*Token, error)
	Find(id string) (*Token, error)
	FindForUser(userId uint) (*Token, error)
	ForUser(userId uint) ([]*Token, error)
	GetValidToken(user interface{}, client *Client) (*Token, error)
	Save(token *Token) error
	RevokeAccessToken(id string) error
	IsAccessTokenRevoked(id string) bool
}
type tokenRepository struct {
	Tokens map[string]*Token
}

func NewTokenRepository() TokenRepository {
	return &tokenRepository{
		Tokens: make(map[string]*Token),
	}
}

func (r *tokenRepository) Create(attributes map[string]interface{}) (*Token, error) {
	token := &Token{
		ID:          uuid.New().String(),
		TOKEN:       strings.TrimSpace(attributes["token"].(string)),
		ACCESSTOKEN: strings.TrimSpace(attributes["access_token"].(string)),
		REVOKED:     attributes["revoked"].(bool),
		EXPIRES_AT:  time.Now().Add(time.Duration(attributes["expires_in"].(float64)) * time.Second),
		USERID:      attributes["user_id"].(uint),
		CLIENTID:    attributes["client_id"].(uint),
	}

	id := uuid.New().String()
	token.ID = id

	r.Tokens[id] = token

	return token, nil
}

func (r *tokenRepository) Find(id string) (*Token, error) {
	token, ok := r.Tokens[id]
	if !ok {
		return nil, nil
	}

	return token, nil
}

func (r *tokenRepository) FindForUser(userId uint) (*Token, error) {
	for _, token := range r.Tokens {
		if token.USERID == userId {
			return token, nil
		}
	}

	return nil, nil
}

func (r *tokenRepository) ForUser(userId uint) ([]*Token, error) {
	tokens := make([]*Token, 0)
	for _, token := range r.Tokens {
		if token.USERID == userId {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (r *tokenRepository) GetValidToken(user interface{}, client *Client) (*Token, error) {
	for _, token := range r.Tokens {
		if strconv.FormatUint(uint64(token.USERID), 10) == user && strconv.Itoa(int(token.CLIENTID)) == strconv.Itoa(client.ID) && !token.REVOKED && token.EXPIRES_AT.After(time.Now()) {
			return token, nil
		}
	}

	return nil, nil
}

func (r *tokenRepository) Save(token *Token) error {
	r.Tokens[token.ID] = token

	return nil
}

func (r *tokenRepository) RevokeAccessToken(id string) error {
	token, ok := r.Tokens[id]
	if !ok {
		return fmt.Errorf("token not found")
	}

	token.REVOKED = true

	return nil
}

func (r *tokenRepository) IsAccessTokenRevoked(id string) bool {
	token, ok := r.Tokens[id]
	if !ok {
		return true
	}

	return token.REVOKED
}
