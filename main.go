package main

import (
	"log"
	"os"
	"strings"

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
	pushSecretCommand = pushCommand.Flag("secret", "secret file name stored under vault-data/ directory").Required().String()
	cachedir = kingpin.Flag("cache-dir", "cache directory").Default("vault-data").String()
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

	opts := kingpin.Parse()

	cache, err := cache.New(*cachedir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using cache-dir=%s\n", *cachedir)

	sync := sync.New(cache, consul, vault)

	switch opts {
	case "pull":
		_ = *pullCommand
		if err := sync.Pull(*pullForceCommand); err != nil {
			log.Println(err)
		}

	case "push":
		key := strings.Replace(*pushSecretCommand, *cachedir + "/", "", 1)
		log.Printf("Pushing %s to vault\n", key)
		sync.Push(key)
	}
}
