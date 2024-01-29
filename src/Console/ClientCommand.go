package console

import (
	"fmt"
	"strings"
)

type ClientCommand struct {
	client *passport.ClientRepository
}

func NewClientCommand(client *passport.ClientRepository) *ClientCommand {
	return &ClientCommand{client: client}
}

func (c *ClientCommand) Execute() error {
	options := map[string]string{
		"personal": "Create a personal access token client",
		"password": "Create a password grant client",
		"client": "Create a client credentials grant client",
		"name": "The name of the client",
		"provider": "The name of the user provider",
		"redirect_uri": "The URI to redirect to after authorization",
		"user_id": "The user ID the client should be assigned to",
		"public": "Create a public client (Auth code grant type only)",
	}

	for _, option := range []string{"personal", "password", "client"} {
		if c.option(option) {
			if err := c.validateOptions(); err != nil {
				return err
			}

			if option == "personal" {
				c.createPersonalClient()
			} else if option == "password" {
				c.createPasswordClient()
			} else {
				c.createClientCredentialsClient()
			}
			return nil
		}
	}

	c.createAuthCodeClient()
	return nil
}

func (c *ClientCommand) option(name string) bool {
	if name, exists := c.flags().GetString(""+name); exists {
		return true
	}
	return false
}

func (c *ClientCommand) validateOptions() error {
	if !c.option("name") {
		return fmt.Errorf("The `--name` option is required")
	}

	if !c.option("user_id") && c.option("redirect_uri") {
		return fmt.Errorf("The `--user_id` option is required when specifying the `--redirect_uri` option")
	}

	return nil
}

func (c *ClientCommand) createPersonalClient() {
	name := c.option("name")
	redirect := c.option("redirect_uri")

	client, err := c.client.CreatePersonalAccessClient(name, redirect)
	if err != nil {
		fmt.Println("Error creating personal access client:", err)
		return
	}

	fmt.Println("Personal access client created successfully.")
	c.outputClientDetails(client)
}

func (c *ClientCommand) createPasswordClient() {
	name := c.option("name")
	redirect := c.option("redirect_uri")
	provider := c.option("provider")

	client, err := c.client.CreatePasswordGrantClient(name, redirect, provider)
	if err != nil {
		fmt.Println("Error creating password grant client:", err)
		return
	}

	fmt.Println("Password grant client created successfully.")
	c.outputClientDetails(client)
}

func (c *ClientCommand) createClientCredentialsClient() {
	name := c.option("name")

	client, err := c.client.Create(nil, name, "")
	if err != nil {
		fmt.Println("Error creating client credentials grant client:", err)
		return
	}

	fmt.Println("Client credentials grant client created successfully.")
	c.outputClientDetails(client)
}

func (c *ClientCommand) createAuthCodeClient() {
	userId := c.option("user_id")
	name := c.option("name")
	redirect := c.option("redirect_uri")

	client, err := c.client.Create(userId, name, redirect, nil, false, false, false)
	if err != nil {
		fmt.Println("Error creating auth code grant client:", err)
		return
	}

	fmt.Println("Auth code grant client created successfully.")
}