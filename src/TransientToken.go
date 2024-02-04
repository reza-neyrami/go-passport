package src

type TransientToken interface {
    Can(scope string) bool
    Cant(scope string) bool
    Transient() bool
}

type transientToken struct {
    scopes []string
}

func NewTransientToken(scopes []string) TransientToken {
    return &transientToken{
        scopes: scopes,
    }
}

func (t *transientToken) Can(scope string) bool {
    for _, tokenScope := range t.scopes {
        if tokenScope == scope {
            return true
        }
    }

    return false
}

func (t *transientToken) Cant(scope string) bool {
    return !t.Can(scope)
}

func (t *transientToken) Transient() bool {
    return true
}
