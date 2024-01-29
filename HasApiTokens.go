package main

import (
	"context"
	"fmt"
	"time"
)

type User struct {
    ID        int
    Name      string
    email     string
    clients   []Client
    tokens   []Token
    accessToken Token
}

func (u *User) Clients(ctx context.Context) ([]Client, error) {
    clients, err := FindClientsForUser(ctx, u.ID)
    if err != nil {
        return nil, err
    }

    for i, client := range clients {
        if client.Revoked {
            clients = append(clients[:i], clients[i+1:]...)
            i--
        }
    }

    return clients, nil
}

func (u *User) Tokens(ctx context.Context) ([]Token, error) {
    tokens, err := FindTokensForUser(ctx, u.ID)
    if err != nil {
        return nil, err
    }

    for i, token := range tokens {
        if token.Revoked {
            tokens = append(tokens[:i], tokens[i+1:]...)
            i--
        }
    }

    return tokens, nil
}

func (u *User) AccessToken(ctx context.Context) (Token, error) {
    tokens, err := u.Tokens(ctx)
    if err != nil {
        return Token{}, err
    }

    for _, token := range tokens {
        if token.IsPersonalAccess {
            return token, nil
        }
    }

    return Token{}, ErrNoAccessToken
}

func (u *User) CreateToken(ctx context.Context, name string, scopes []string) (PersonalAccessTokenResult, error) {
    token, err := CreatePersonalAccessClient(ctx, u.ID, name, fmt.Sprintf("%s:%d", u.ID, time.Now().UnixNano()))
    if err != nil {
        return PersonalAccessTokenResult{}, err
    }

    u.accessToken = token

    return PersonalAccessTokenResult{AccessToken: token}, nil
}
