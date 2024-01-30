package bridge

import "encoding/json"

func createScope(identifier string) *Scope {
	return &Scope{
		Identifier: identifier,
	}
}


func (scope *Scope) GetIdentifier() string {
	return scope.Identifier
}

type FormatsScopesForStorage interface {
	formatScopesForStorage(scopes []Scope) string
}

func formatScopesForStorage(scopes []Scope) string {
	b, err := json.Marshal(scopesToArray(scopes))
	if err != nil {
		return ""
	}
	return string(b)
}

func scopesToArray(scopes []Scope) []string {
	mappedScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		mappedScopes[i] = scope.Identifier
	}
	return mappedScopes
}
