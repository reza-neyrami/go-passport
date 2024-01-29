package bridge

import (
    "database/sql"
    "fmt"
)

type ClientRepository struct {
    db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
    return &ClientRepository{
        db: db,
    }
}

func (r *ClientRepository) getClientEntity(clientID string) (*Client, error) {
    client := &Client{}

    err := r.db.QueryRow("SELECT name, redirect_uri, confidential, provider FROM clients WHERE id = ?", clientID).Scan(&client.Name, &client.RedirectURIs, &client.IsConfidential, &client.Provider)
    if err != nil {
        return nil, err
    }

    return client, nil
}

func (r *ClientRepository) validateClient(clientID, clientSecret string, grantType string) (bool, error) {
    client := &Client{}

    err := r.db.QueryRow("SELECT secret, first_party, personal_access_client, password_client FROM clients WHERE id = ?", clientID).Scan(&client.Secret, &client.FirstParty, &client.PersonalAccessClient, &client.PasswordClient)
    if err != nil {
        return false, err
    }

    if clientSecret != "" && !r.verifySecret(clientSecret, client.Secret) {
        return false, fmt.Errorf("Invalid client secret")
    }

    switch grantType {
    case "authorization_code":
        return !client.FirstParty, nil
    case "personal_access":
        return client.PersonalAccessClient && client.IsConfidential, nil
    case "password":
        return client.PasswordClient, nil
    case "client_credentials":
        return client.IsConfidential, nil
    default:
        return true, nil
    }
}

func (r *ClientRepository) verifySecret(clientSecret, storedHash string) bool {
    if Passport.HashesClientSecrets {
        return password.Verify(clientSecret, storedHash) == nil
    }

    return storedHash == clientSecret
}
