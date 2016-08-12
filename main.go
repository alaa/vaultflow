package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alaa/vaultflow/consul"
	"github.com/alaa/vaultflow/vault"
)

func main() {
	vaultToken := os.Getenv("VAULT_TOKEN")
	vault, err := vault.New(vaultToken)
	if err != nil {
		log.Println("Cloud not initialize vault client", err)
	}
	client := consul.New()

	// Lock
	sessionID, err := client.AcquireLock()
	if err != nil {
		log.Fatal(err)
	}

	// Update revision
	err = client.UpdateRevision(sessionID)
	if err != nil {
		log.Fatal(err)
	}

	// Critical section
	keys, err := vault.ListSecrets("/secret")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, key := range keys {
		secret, err := vault.ReadSecret("secret/" + key)
		if err != nil {
			log.Println(err)
		}
		formatter(key, secret.Data)
	}

	// Unlock
	if err = client.ReleaseLock(sessionID); err != nil {
		log.Fatalf("Could not release global lock with err: %s \n", err)
	}
}

func formatter(key string, secret map[string]interface{}) {
	json, err := json.Marshal(secret)
	if err != nil {
		log.Fatalf("Could not marshal secret: %s", err)
	}
	fmt.Printf("==> %s:\n%s\n\n", key, json)
}
