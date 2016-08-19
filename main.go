package main

import (
	"fmt"
	"log"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/alaa/vaultflow/cache"
	"github.com/alaa/vaultflow/consul"
	"github.com/alaa/vaultflow/sync"
	"github.com/alaa/vaultflow/vault"
)

var (
	pullCommand       = kingpin.Command("pull", "pull latest secrets from vault.")
	pushCommand       = kingpin.Command("push", "Push new secret to vault.")
	pushSecretCommand = pushCommand.Flag("secret", "secret file name stored under .cache/ directory").Required().String()
)

func main() {
	vaultToken := os.Getenv("VAULT_TOKEN")
	vault, err := vault.New(vaultToken)
	if err != nil {
		log.Fatal(err)
	}

	consul, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := cache.New()
	if err != nil {
		log.Fatal(err)
	}

	sync := sync.New(cache, consul, vault)

	// CLI options:
	// Pull: vaultflow pull
	// Push: vaultflow push --secret=sercret_file
	switch kingpin.Parse() {
	case "pull":
		_ = *pullCommand
		if err := sync.Pull(); err != nil {
			log.Println(err)
		}

	case "push":
		fmt.Println("just triggered push secret command")
		fmt.Println(*pushSecretCommand)
		sync.Push(*pushSecretCommand)
	}
}
