package bridge

import (
    "database/sql"
    "fmt"
)

type ScopeRepository struct {
    db *sql.DB
}

func NewScopeRepository(db *sql.DB) *ScopeRepository {
    return &ScopeRepository{
        db: db,
    }
}

func (r *ScopeRepository) GetScopeEntityByIdentifier(identifier string) (*Scope, error) {
    if _, err := r.db.Exec("INSERT INTO scopes (name) VALUES (?)", identifier); err != nil {
        return nil, err
    }

    return &Scope{
	Name: identifier,
}, nil
}

func (r *ScopeRepository) FinalizeScopes(
    scopes []Scope,
    grantType string,
    client *Client,
    userIdentifier string) []Scope {

    if grantType != "password" && grantType != "personal_access" && grantType != "client_credentials" {
        scopes = scopes.Filter(func(scope Scope) bool {
            return scope.Name != "*"
        })
    }

    // TODO: Implement client validation
    return scopes
}
