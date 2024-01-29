package bridge



type Scope struct {
	Identifier string `json:"identifier"`
}

func NewScope(identifier string) *Scope {
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
	return json.Marshal(scopesToArray(scopes))
}

func scopesToArray(scopes []Scope) []string {
	return array_map(func(scope Scope) string {
		return scope.Identifier
	}, scopes)
}
