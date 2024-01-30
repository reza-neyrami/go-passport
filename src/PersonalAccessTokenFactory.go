package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type PersonalAccessTokenFactory struct {
    server *AuthorizationServer
    clients *ClientRepository
    tokens *TokenRepository
    jwt JwtParser
}

func NewPersonalAccessTokenFactory(server *AuthorizationServer, clients *ClientRepository, tokens *TokenRepository, jwt JwtParser) *PersonalAccessTokenFactory {
    return &PersonalAccessTokenFactory{
        server: server,
        clients: clients,
        tokens: tokens,
        jwt: jwt,
    }
}

func (p *PersonalAccessTokenFactory) Make(userId interface{}, name string, scopes []string) (PersonalAccessTokenResult, error) {
    response, err := p.dispatchRequestToAuthorizationServer(p.createRequest(p.clients.PersonalAccessClient(), userId, scopes))
    if err != nil {
        return PersonalAccessTokenResult{}, fmt.Errorf("Error creating personal access token: %v", err)
    }

    token, err := p.findAccessToken(response)
    if err != nil {
        return PersonalAccessTokenResult{}, fmt.Errorf("Error finding access token: %v", err)
    }

    err = p.tokens.Save(token.ForceFill(map[string]interface{}{
        "user_id": userId,
        "name": name,
    }))
    if err != nil {
        return PersonalAccessTokenResult{}, fmt.Errorf("Error updating token: %v", err)
    }

    return PersonalAccessTokenResult{
        AccessToken: response["access_token"],
        Token: token,
    }, nil
}

func (p *PersonalAccessTokenFactory) createRequest(client Client, userId interface{}, scopes []string) *ServerRequest {
    var secret string
    if Passport.HashesClientSecrets {
        secret = p.clients.GetPersonalAccessClientSecret()
    } else {
        secret = client.Secret
    }
    
    request := &ServerRequest{Method: "POST", Uri: "not-important"}
    request.SetBody([]byte(`grant_type=personal_access&client_id=` + client.Key() + `&client_secret=` + secret + `&user_id=` + fmt.Sprint(userId) + `&scope=` + strings.Join(scopes, " ")), "application/x-www-form-urlencoded")

    return request
}

func (p *PersonalAccessTokenFactory) dispatchRequestToAuthorizationServer(request *ServerRequest) (map[string]interface{}) {
    response, err := p.server.RespondToAccessTokenRequest(request, new(Response))
    if err != nil {
        log.Printf("Error dispatching request: %v", err)
        return nil
    }

    var parsedResponse map[string]interface{}
    err = json.Unmarshal([]byte(response.Body.String()), &parsedResponse)
    if err != nil {
        log.Printf("Error parsing response: %v", err)
        return nil
    }

    return parsedResponse
}

func (p *PersonalAccessTokenFactory) findAccessToken(response map[string]interface{}) (Token) {
    return p.tokens.Find(p.jwt.Parse(response["access_token"]).Claims().Get("jti"))
}
