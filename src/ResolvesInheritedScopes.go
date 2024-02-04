package src

import "strings"

func ResolveInheritedScopes(scope string) []string {
	parts := strings.Split(scope, ":")
	partsCount := len(parts)
	scopes := []string{}

	for i := 1; i <= partsCount; i++ {
		scopes = append(scopes, strings.Join(parts[:i], ":"))
	}

	return scopes
}
