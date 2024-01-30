package bridge

import "encoding/json"

type Scope struct {
	Name string `json:"name"`
	
	Identifier string `json:"identifier"`
}

func NewScope(name string) *Scope {
	return &Scope{
		Name: name,
	}
}

func (scope *Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(scope.Name)
}
