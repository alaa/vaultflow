package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/alaa/vaultflow/cache"
	"github.com/alaa/vaultflow/consul"
	"github.com/alaa/vaultflow/vault"
)

type Sync struct {
	cache  *cache.Cache
	consul *consul.Consul
	vault  *vault.Vault
}

func New(cache *cache.Cache, consul *consul.Consul, vault *vault.Vault) *Sync {
	return &Sync{
		cache:  cache,
		consul: consul,
		vault:  vault,
	}
}

type revisions map[string]string

func (s *Sync) getRevisions() (revisions, error) {
	remoteRevision, err := s.consul.GetRevision()
	if err != nil {
		return revisions{}, err
	}

	localRevision, err := s.cache.GetRevision()
	if err != nil {
		return revisions{}, err
	}

	rev := make(revisions)
	rev["remote"] = remoteRevision
	rev["local"] = localRevision
	return rev, nil
}

func (s *Sync) isLatest() (bool, revisions, error) {
	revisions, err := s.getRevisions()
	if err != nil {
		return false, revisions, err
	}

	if revisions["local"] == revisions["remote"] && revisions["remote"] != "" {
		return true, revisions, nil
	}

	return false, revisions, nil
}

// Lock, Accquires lock on consul
// Lock should be always followed by defer Release.
func (s *Sync) lock() (string, error) {
	sessionID, err := s.consul.AcquireLock()
	return sessionID, err
}

// Release, Invalide the lock with certain sessionID.
func (s *Sync) release(sessionID string) error {
	if err := s.consul.ReleaseLock(sessionID); err != nil {
		return err
	}
	return nil
}

// Pull, Get latest secrets from vault and update local revisionID.
func (s *Sync) Pull(force bool) error {
	synced, revisions, err := s.isLatest()
	if err != nil {
		return err
	}

	if !force && synced {
		return errors.New(fmt.Sprintf("You have the latest copy of secrets: %s\n", revisions["local"]))
	}

	secrets, err := s.vault.GetSecrets()
	if err != nil {
		return err
	}

	for key, secret := range secrets {
		json, err := json.Marshal(secret)
		if err != nil {
			return err
		}

		if err = s.cache.Put([]byte(key), json); err != nil {
			return errors.New(fmt.Sprintf("Could not write secret to the cache: %s", err))
		}
	}

	if err = s.cache.UpdateRevision(revisions["remote"]); err != nil {
		return errors.New(fmt.Sprintf("Could not update cache revision: %s\n", err))
	}
	return nil
}

// Push, Add/Update remote secret on Vault.
func (s *Sync) Push(key string) {
	synced, revisions, err := s.isLatest()
	if err != nil {
		log.Panic(err)
	}

	if !synced {
		log.Print("You have older version of Secrets")
		log.Printf("Latest RevisionID is: %s\n", revisions["remote"])
		return
	}

	sessionID, err := s.lock()
	if err != nil {
		log.Panic(err)
	}
	defer func(id string) {
		if err := s.release(id); err != nil {
			log.Fatal(err)
		}
	}(sessionID)

	data, err := s.cache.Get([]byte(key))
	if err != nil {
		log.Panic(err)
	}

	// Push bytes to vault
	_, err = s.vault.WriteSecret(key, data)
	if err != nil {
		log.Panic(err)
	}

	if err = s.cache.UpdateRevision(sessionID); err != nil {
		log.Panicf("Could not update cache revision: %s\n", err)
	}

	err = s.consul.UpdateRevision(sessionID)
	if err != nil {
		log.Panicf("Could not update revision on consul: %s\n", err)
	}
}
