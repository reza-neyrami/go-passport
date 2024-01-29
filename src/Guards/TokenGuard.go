package guards

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"

    "github.com/dgrijalva/jwt-go"
)

type Guard interface {
    User() (user *User, err error)
    Validate(credentials map[string]interface{}) bool
    Client() (client *Client, err error)
}

type TokenGuard struct {
    provider        UserProvider
    tokenRepository TokenRepository
    encrypter       Encrypter
    request         *http.Request
    client          *Client
}

func NewTokenGuard(
    provider UserProvider,
    tokenRepository TokenRepository,
    encrypter Encrypter,
    request *http.Request,
) *TokenGuard {
    return &TokenGuard{
        provider:         provider,
        tokenRepository:  tokenRepository,
        encrypter:        encrypter,
        request:         request,
        client:          nil,
    }
}

func (g *TokenGuard) User() (user *User, err error) {
    user, err = g.userFromToken(g.request)
    if err != nil {
        return
    }

    g.client = g.tokenRepository.FindClient(user.ClientID)

    return
}

func (g *TokenGuard) Validate(credentials map[string]interface{}) bool {
    return g.userFromCredentials(credentials) != nil
}

func (g *TokenGuard) Client() (client *Client, err error) {
    client, err = g.clientFromToken(g.request)
    if err != nil {
        return
    }

    return
}

func (g *TokenGuard) userFromToken(request *http.Request) (user *User, err error) {
    token, err := g.decodeJwtToken(request)
    if err != nil {
        return
    }

    user, err = g.provider.RetrieveByID(token["user_id"].(string))
    if err != nil {
        return
    }

    return
}

func (g *TokenGuard) clientFromToken(request *http.Request) (client *Client, err error) {
    token, err := g.decodeJwtToken(request)
    if err != nil {
        return
    }

    client, err = g.tokenRepository.FindClient(token["client_id"].(string))
    if err != nil {
        return
    }

    return
}

func (g *TokenGuard) decodeJwtToken(request *http.Request) (token map[string]interface{}, err error) {
    jwt := request.Header.Get("Authorization")
    if jwt == "" {
        return nil, errors.New("Authorization header is missing")
    }

    token, err = jwt.Decode(g.encrypter.PublicKey())
    if err != nil {
        return nil, errors.New("Failed to decode JWT token")
    }

    if clientID, ok := token["client_id"].(string); !ok {
        return nil, errors.New("JWT token is missing client_id claim")
    }

    if g.provider.Clients[clientID] == nil {
        return nil, errors.New("Client with ID '" + clientID + "' is not registered")
    }

    return token, nil
}

func (g *TokenGuard) userFromCredentials(credentials map[string]interface{}) (user *User) {
    tokenID, ok := credentials["token_id"].(string)
    if !ok {
        return nil
    }

    token, err := g.tokenRepository.Find(tokenID)
    if err != nil {
        return nil
    }

    user, err = g.provider.RetrieveByID(token.UserID)
    if err != nil {
        return nil
    }

    return user
}

func (g *TokenGuard) SetRequest(request *http.Request) *TokenGuard {
    g.request = request

    return g
}

