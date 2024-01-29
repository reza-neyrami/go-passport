package console

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type InstallCommand struct {
	uuids  bool
	force  bool
	length int
}

func NewInstallCommand() *InstallCommand {
	uuids := flag.Bool("uuids", false, "Use UUIDs for all client IDs")
	force := flag.Bool("force", false, "Overwrite keys they already exist")
	length := flag.Int("length", 4096, "The length of the private key")

	flag.Parse()

	return &InstallCommand{
		uuids:  *uuids,
		force:  *force,
		length: *length,
	}
}

func (c *InstallCommand) Handle() {
	fmt.Println("Running the commands necessary to prepare Passport for use...")

	if c.uuids {
		c.configureUuids()
	}

	fmt.Println("Creating Personal Access Client...")
	fmt.Println("Creating Password Grant Client...")
}

func (c *InstallCommand) configureUuids() {
	fmt.Println("Configuring Passport for client UUIDs...")

	c.replaceInFile("passport.go", "'client_uuids' => false", "'client_uuids' => true")
	c.replaceInFile("oauth_auth_codes_table.go", "client_id bigint", "client_id uuid")
	c.replaceInFile("oauth_access_tokens_table.go", "client_id bigint", "client_id uuid")
	c.replaceInFile("oauth_clients_table.go", "id bigint", "id uuid")
	c.replaceInFile("oauth_personal_access_clients_table.go", "client_id bigint", "client_id uuid")

	fmt.Println("Finished configuring client UUIDs.")
}

func (c *InstallCommand) replaceInFile(filename, search, replace string) {
	read, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	newContents := strings.Replace(string(read), search, replace, -1)

	err = os.WriteFile(filename, []byte(newContents), 0)
	if err != nil {
		log.Fatal(err)
	}
}


