package main

import (
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
	pullForceCommand  = pullCommand.Flag("force", "force pull secrets").Bool()
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

	cache, err := cache.New(".cache2")
	if err != nil {
		log.Fatal(err)
	}

	sync := sync.New(cache, consul, vault)

	switch kingpin.Parse() {
	case "pull":
		_ = *pullCommand
		if err := sync.Pull(*pullForceCommand); err != nil {
			log.Println(err)
		}

	case "push":
		sync.Push(*pushSecretCommand)
	}
}
