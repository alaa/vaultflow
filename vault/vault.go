package vault

import (
	"encoding/json"
	"errors"

	vaultapi "github.com/hashicorp/vault/api"
)

type Vault struct {
	client  *vaultapi.Client
	logical *vaultapi.Logical
}

func New(vaultToken string) (*Vault, error) {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return &Vault{}, err
	}

	if vaultToken != "" {
		client.SetToken(vaultToken)
	}

	return &Vault{
		client:  client,
		logical: client.Logical(),
	}, nil
}

func (v *Vault) ReadSecret(path string) (*vaultapi.Secret, error) {
	secret, err := v.logical.Read(path)
	if err != nil {
		return &vaultapi.Secret{}, err
	}
	return secret, nil
}

func (v *Vault) WriteSecret(path string, data []byte) (*vaultapi.Secret, error) {
	var s map[string]interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return &vaultapi.Secret{}, err
	}

	vaultpath := "secret/" + path
	secret, err := v.logical.Write(vaultpath, s)
	if err != nil {
		return &vaultapi.Secret{}, err
	}
	return secret, nil
}

func (v *Vault) ListSecrets(path string) ([]string, error) {
	paths, err := v.logical.List(path)
	if err != nil {
		return nil, err
	}

	var s []string
	for _, keys := range paths.Data {
		keysSlice := keys.([]interface{})
		for _, key := range keysSlice {
			s = append(s, key.(string))
		}
		return s, nil
	}
	return nil, errors.New("Error fetching vault secrets list")
}

type Secrets map[string]map[string]interface{}

func (v *Vault) GetSecrets() (Secrets, error) {
	secrets := make(Secrets)
	keys, err := v.ListSecrets("/secret")
	if err != nil {
		return Secrets{}, nil
	}

	for _, key := range keys {
		secret, err := v.ReadSecret("secret/" + key)
		if err != nil {
			return Secrets{}, nil
		}
		secrets[key] = secret.Data
	}

	return secrets, nil
}
