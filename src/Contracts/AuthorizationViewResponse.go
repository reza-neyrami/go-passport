package main

type Parameters map[string]interface{}

type AuthorizationViewResponse interface {
	WithParameters(parameters Parameters) AuthorizationViewResponse
}

type MyAuthorizationViewResponse struct {
	parameters Parameters
}

func (r *MyAuthorizationViewResponse) WithParameters(parameters Parameters) AuthorizationViewResponse {
	r.parameters = parameters
	return r
}

func main() {
	response := &MyAuthorizationViewResponse{}
	response.WithParameters(Parameters{"key": "value"})
}
