package main

type Scope struct {
    ID   string
    DESC string
}

func NewScope(id string, description string) *Scope {
    return &Scope{
        id:   id,
        desc: description,
    }
}

func (s *Scope) ToJSON(options int) string {
    data, err := json.Marshal(s)
    if err != nil {
        return ""
    }

    return string(data)
}

func (s *Scope) ToMap() map[string]interface{} {
    data := make(map[string]interface{})
    data["id"] = s.id
    data["description"] = s.desc

    return data
}
