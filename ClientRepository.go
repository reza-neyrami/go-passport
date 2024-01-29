package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
)



type ClientRepository interface {
    Create(ctx context.Context, userId int, name string, redirect string, provider *string, personalAccess bool, password bool, confidential bool) (Client, error)
    Update(ctx context.Context, client *Client, name string, redirect string) error
    Delete(ctx context.Context, client *Client) error
    Find(ctx context.Context, id int) (*Client, error)
    FindActive(ctx context.Context, id int) (*Client, error)
    FindForUser(ctx context.Context, clientId int, userId int) (*Client, error)
    ForUser(ctx context.Context, userId int) ([]Client, error)
    ActiveForUser(ctx context.Context, userId int) ([]Client, error)
    CreatePersonalAccessClient(ctx context.Context, userId int, name string, redirect string) (*Client, error)
    CreatePasswordGrantClient(ctx context.Context, userId int, name string, redirect string, provider *string) (*Client, error)
    Revoke(ctx context.Context, id int) error
    GetPersonalAccessClientId() int
    GetPersonalAccessClientSecret() string
}

var (
    ErrInvalidClientID = errors.New("invalid client ID")
    ErrInvalidClientSecret = errors.New("invalid client secret")
)

func NewClientRepository() ClientRepository {
    return &clientRepository{}
}

type clientRepository struct {
}



func (r *clientRepository) Create(ctx context.Context, userId int, name string, redirect string, provider *string, personalAccess bool, password bool, confidential bool) (Client, error) {
    client := Client{
        UserId :    userId,
        Name:      name,
        Secret:    fmt.Sprintf("%s-%d", name, time.Now().UnixNano()),
        Provider:  *provider,
        Redirect:  redirect,
        PersonalAccess: personalAccess,
        Password:      password,
        Revoked:    false,
    }

    if confidential {
        client.Secret = RandomString(40)
    }

    if err := client.Save(ctx); err != nil {
        return Client{}, err
    }

    if personalAccess {
        err := r.createPersonalAccessClient(ctx, client.ID)
        if err != nil {
            return Client{}, err
        }
    }

    return client, nil
}

func (r *clientRepository) Update(ctx context.Context, client *Client, name string, redirect string) error {
    if err := client.Validate(); err != nil {
        return err
    }

    client.Name = name
    client.Redirect = redirect

    return client.Update(ctx)
}

func (r *clientRepository) Delete(ctx context.Context, client *Client) error {
    if err := client.Validate(); err != nil {
        return err
    }

    return client.Delete(ctx)
}

func (r *clientRepository) Find(ctx context.Context, id int) (*Client, error) {
    client, err := FindClient(ctx, id)
    if err != nil {
        return nil, err
    }

    if client == nil || client.Revoked {
        return nil, ErrInvalidClientID
    }

    return client, nil
}

func (r *clientRepository) FindActive(ctx context.Context, id int) (*Client, error) {
    client, err := FindActiveClient(ctx, id)
    if err != nil {
        return nil, err
    }

    if client == nil || client.Revoked {
        return nil, ErrInvalidClientID
    }

    return client, nil
}

func (r *clientRepository) FindForUser(ctx context.Context, clientId int, userId int) (*Client, error) {
    client, err := FindClientForUser(ctx, clientId, userId)
    if err != nil {
        return nil, err
    }

    if client == nil || client.Revoked {
        return nil, ErrInvalidClientID
    }

    return client, nil
}

func (r *clientRepository) ForUser(ctx context.Context, userId int) ([]Client, error) {
    clients, err := FindClientsForUser(ctx, userId)
    if err != nil {
        return nil, err
    }

    for _, client := range clients {
        if client == nil || client.Revoked {
            clients = append(clients[:i], clients[i+1:]...)
            i--
        }
    }

    return clients, nil
}

func (r *clientRepository) ActiveForUser(ctx context.Context, userId int) ([]Client, error) {
    clients, err := FindActiveClientsForUser(ctx, userId)
    if err != nil {
        return nil, err
    }

    for _, client := range clients {
        if client == nil || client.Revoked {
            clients = append(clients[:i], clients[i+1:]...)
            i--
        }
    }

    return clients, nil
}

func (r *clientRepository) CreatePersonalAccessClient(ctx context.Context, userId int, name string, redirect string) (*Client, error) {
    client, err := r.Create(ctx, userId, name, redirect, nil, true, false, false)
    if err != nil {
        return nil, err
    }

    r.updatePersonalAccessClient(ctx, client.ID)

    return client, nil
}

func (r *clientRepository) CreatePasswordGrantClient(ctx context.Context, userId int, name string, redirect string, provider *string) (*Client, error) {
    client, err := r.Create(ctx, userId, name, redirect, provider, false, true, false)
    if err != nil {
        return nil, err
    }

    return client, nil
}

func (r *clientRepository) Revoke(ctx context.Context, id int) error {
    client, err := FindClient(ctx, id)
    if err != nil {
        return err
    }

    if client == nil {
        return ErrInvalidClientID
    }

    client.Revoked = true

    return client.Save(ctx)
}

func (r *clientRepository) GetPersonalAccessClientId() int {
    return r.personalAccessClientId
}

func (r *clientRepository) GetPersonalAccessClientSecret() string {
    return r.personalAccessClientSecret
}

func (r *clientRepository) createPersonalAccessClient(ctx context.Context, clientId int) error {
    client := Client{
        ID:        clientId,
        UserId:    0,
        Name:      "Personal Access Client",
        Secret:    "",
        Provider:  "",
        Redirect:  "",
        PersonalAccess: true,
        Password:      false,
        Revoked:    false,
    }

    err := client.Save(ctx)
    if err != nil {
        return err
    }

    r.personalAccessClientId = client.ID
    r.personalAccessClientSecret = client.Secret

    return nil
}

func (r *clientRepository) updatePersonalAccessClient(ctx context.Context, clientId int) error {
    client, err := FindClient(ctx, clientId)
    if err != nil {
        return err
    }

    if client == nil {
        return ErrInvalidClientID
    }

    client.Secret = r.personalAccessClientSecret

    return client.Save(ctx)
}


func RandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}