package console

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type KeysCommand struct {
	force  bool
	length int
}

func NewKeysCommand() *KeysCommand {
	force := flag.Bool("force", false, "Overwrite keys if they already exist")
	length := flag.Int("length", 4096, "The length of the private key")

	flag.Parse()

	return &KeysCommand{
		force:  *force,
		length: *length,
	}
}

func (c *KeysCommand) Handle() {
	publicKeyPath := "oauth-public.key"
	privateKeyPath := "oauth-private.key"

	if fileExists(publicKeyPath) || fileExists(privateKeyPath) && !c.force {
		fmt.Println("Encryption keys already exist. Use the --force option to overwrite them.")
		return
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, c.length)
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile(privateKeyPath, privateKeyPEM, 0600)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(publicKeyPath, publicKeyPEM, 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encryption keys generated successfully.")
}

func (c *KeysCommand) WriteToFile(publicKeyPath, privateKeyPath string) error {
	publicKey, err := rsa.GenerateKey(rand.Reader, c.length)
	if err != nil {
		return err
	}

	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(publicKeyPath, publicKeyPEM, 0600)
	if err != nil {
		return err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(publicKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile(privateKeyPath, publicKeyPEM, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *KeysCommand) confirm(message string) bool {
	var response string
	fmt.Println(message)
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return false
	}

	response = strings.ToLower(response)
	if response == "y" || response == "yes" {
		return true
	} else if response == "n" || response == "no" {
		return false
	}
	fmt.Println("Invalid response. Please enter 'y' or 'n'.")
	return c.confirm(message)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}
