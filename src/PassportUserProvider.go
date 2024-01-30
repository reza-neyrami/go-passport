package main

import (
    "fmt"
    "log"
)

type PassportUserProvider struct {
    provider UserProvider
    providerName string
}

func NewPassportUserProvider(provider UserProvider, providerName string) *PassportUserProvider {
    return &PassportUserProvider{
        provider: provider,
        providerName: providerName,
    }
}

func (p *PassportUserProvider) RetrieveById(identifier interface{}) (Authenticatable, error) {
    user, err := p.provider.RetrieveById(identifier)
    if err != nil {
        log.Printf("Error retrieving user by ID from passport provider: %v", err)
        return nil, err
    }

    return user, nil
}

func (p *PassportUserProvider) RetrieveByToken(identifier interface{}, token string) (Authenticatable, error) {
    user, err := p.provider.RetrieveByToken(identifier, token)
    if err != nil {
        log.Printf("Error retrieving user by token from passport provider: %v", err)
        return nil, err
    }

    return user, nil
}

func (p *PassportUserProvider) UpdateRememberToken(user Authenticatable, token string) error {
    err := p.provider.UpdateRememberToken(user, token)
    if err != nil {
        log.Printf("Error updating remember token for user from passport provider: %v", err)
        return err
    }

    return nil
}

func (p *PassportUserProvider) RetrieveByCredentials(credentials map[string]interface{}) (Authenticatable, error) {
    user, err := p.provider.RetrieveByCredentials(credentials)
    if err != nil {
        log.Printf("Error retrieving user by credentials from passport provider: %v", err)
        return nil, err
    }

    return user, nil
}

func (p *PassportUserProvider) ValidateCredentials(user Authenticatable, credentials map[string]interface{}) bool {
    valid := p.provider.ValidateCredentials(user, credentials)
    if !valid {
        log.Printf("Error validating credentials for user from passport provider: %v", err)
    }

    return valid
}

func (p *PassportUserProvider) GetProviderName() string {
    return p.providerName
}
