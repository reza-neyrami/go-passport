package bridge

import (
    "database/sql"
    "time"
)

type PersonalAccessGrant struct {
    db *sql.DB
}

func NewPersonalAccessGrant(db *sql.DB) *PersonalAccessGrant {
    return &PersonalAccessGrant{
        db: db,
    }
}

func (grant *PersonalAccessGrant) RespondToAccessTokenRequest(
    request *http.Request,
    responseType *ResponseTypeInterface,
    accessTokenTTL time.Duration) (*ResponseTypeInterface, error) {

    // Validate request
    client, err := validateClient(request)
    if err != nil {
        return responseType, err
    }

    scopes, err := validateScopes(request)
    if err != nil {
        return responseType, err
    }

    userIdentifier, err := request.FormValue("user_id")
    if err != nil {
        return responseType, err
    }

    // Finalize the requested scopes
    scopes = finalizeScopes(
        scopes,
        grant.Identifier(),
        client,
        userIdentifier,
    )

    // Issue and persist access token
    accessToken, err := grant.issueAccessToken(
        accessTokenTTL,
        client,
        userIdentifier,
        scopes,
    )
    if err != nil {
        return responseType, err
    }

    // Inject access token into response type
    responseType.SetAccessToken(accessToken)

    return responseType, nil
}

func (grant *PersonalAccessGrant) GetIdentifier() string {
    return "personal_access"
}
