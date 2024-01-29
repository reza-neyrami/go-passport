package bridge

type User struct {
    Identifier string `json:"identifier"`
}

func NewUser(identifier string) *User {
    return &User{
        Identifier: identifier,
    }
}

func (user *User) GetIdentifier() string {
    return user.Identifier
}

func (user *User) SetIdentifier(identifier string) {
    user.Identifier = identifier
}
