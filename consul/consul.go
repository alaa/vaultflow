package consul

import (
	"errors"
	"fmt"
	"log"
	"os"

	consulapi "github.com/hashicorp/consul/api"
	multierror "github.com/hashicorp/go-multierror"
)

const lockPath = "vaultflow/lock"
const revisionPath = "vaultflow/revision"

type Consul struct {
	Client *consulapi.Client
}

func New() (*Consul, error) {
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		log.Print("CONSUL_ADDR is not set, using localhost:8500")
		consulAddr = "localhost:8500"
	}

	config := &consulapi.Config{
		Address: consulAddr,
		Scheme:  "http",
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return &Consul{}, err
	}

	return &Consul{
		Client: client,
	}, nil
}

func (c *Consul) createSession() (string, error) {
	session := c.Client.Session()
	sessionID, _, err := session.Create(&consulapi.SessionEntry{
		Behavior: consulapi.SessionBehaviorDelete,
		TTL:      "15s",
	}, nil)
	return sessionID, err
}

func (c *Consul) destroySession(sessionID string) error {
	session := c.Client.Session()
	_, err := session.Destroy(sessionID, nil)
	return err
}

func (c *Consul) AcquireLock() (string, error) {
	err := c.isLocked()
	if err != nil {
		return "", err
	}

	sessionID, err := c.createSession()
	if err != nil {
		return "", err
	}

	kv := c.Client.KV()
	kvpair := &consulapi.KVPair{
		Key:     lockPath,
		Session: sessionID,
		Value:   []byte(sessionID),
	}
	kv.Acquire(kvpair, nil)

	return sessionID, nil
}

func (c *Consul) isLocked() error {
	kv := c.Client.KV()
	pair, _, err := kv.Get(lockPath, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not fetch lock %s", pair))
	}

	if pair != nil && pair.Session != "" {
		return errors.New(fmt.Sprintf("vaultflow session is locked with sessionID: %s", pair.Session))
	}

	return nil
}

func (c *Consul) GetRevision() (string, error) {
	kv := c.Client.KV()
	pair, _, err := kv.Get(revisionPath, nil)
	if err != nil || pair == nil {
		return "", err
	}
	return string(pair.Value), nil
}

func (c *Consul) UpdateRevision(id string) error {
	kv := c.Client.KV()
	_, err := kv.Put(&consulapi.KVPair{
		Key:   revisionPath,
		Value: []byte(id),
	}, nil)
	return err
}

func (c *Consul) ReleaseLock(sessionID string) error {
	var result error

	kv := c.Client.KV()
	_, _, err := kv.Release(&consulapi.KVPair{
		Key:     lockPath,
		Session: sessionID,
	}, nil)
	if err != nil {
		result = multierror.Append(result, err)
	}

	err = c.destroySession(sessionID)
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}
