package src

type PersonalAccessTokenResult struct {
    accessToken string
    token       Token
}

func NewPersonalAccessTokenResult(accessToken string, token Token) *PersonalAccessTokenResult {
    return &PersonalAccessTokenResult{
        accessToken: accessToken,
        token: token,
    }
}

func (p *PersonalAccessTokenResult) ToJSON() (string, error) {
    result, err := json.Marshal(p)
    if err != nil {
        return "", fmt.Errorf("Error marshalling personal access token result: %v", err)
    }

    return string(result), nil
}

func (p *PersonalAccessTokenResult) ToArray() (map[string]interface{}, error) {
    result := map[string]interface{}{
        "accessToken": p.accessToken,
        "token":       p.token,
    }

    return result, nil
}
