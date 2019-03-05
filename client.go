package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/phogolabs/log"
)

var _ Repository = &Client{}

// Client of the vault
type Client struct {
	client   *api.Client
	renewers []*api.Renewer
}

// NewClient creates a new client
func NewClient(config *api.Config) (*Client, error) {
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:   client,
		renewers: []*api.Renewer{},
	}, nil
}

// SetToken sets the token
func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}

// Auth authenticates the client
func (c *Client) Auth(method string, config map[string]interface{}) error {
	client, err := c.client.Clone()
	if err != nil {
		return err
	}

	request := client.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/login", method))
	request.SetJSONBody(&config)

	response, err := client.RawRequest(request)
	if err != nil {
		return err
	}

	secret := &api.Secret{}
	if err = response.DecodeJSON(secret); err != nil {
		return err
	}

	if err = c.renew(client, "token", secret); err != nil {
		return err
	}

	c.client.SetToken(secret.Auth.ClientToken)
	return nil
}

// Stop stops renewing
func (c *Client) Stop() {
	for _, renewer := range c.renewers {
		renewer.Stop()
	}
}

// Mount return the mount for given path
func (c *Client) Mount(path string) (*api.MountOutput, error) {
	request := c.client.NewRequest("GET", fmt.Sprintf("/v1/sys/internal/ui/mounts/%s", path))

	response, err := c.client.RawRequest(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	props := map[string]interface{}{}

	if err = response.DecodeJSON(&props); err != nil {
		return nil, err
	}

	type mount struct {
		Renewable bool            `json:"renewable"`
		Data      api.MountOutput `json:"data" mapstructure:"data"`
	}

	output := &mount{}
	if err = mapstructure.Decode(props, output); err != nil {
		return nil, err
	}

	return &output.Data, nil
}

// Secret returns a secrete
func (c *Client) Secret(path string) (map[string]interface{}, error) {
	mnt, err := c.Mount(path)
	if err != nil {
		return nil, err
	}

	path = abs(path, mnt)

	secret, err := c.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	if err = c.renew(c.client, path, secret); err != nil {
		return nil, err
	}

	switch mnt.Type {
	case "kv":
		if version, ok := mnt.Options["version"]; ok && version == "2" {
			if data, found := secret.Data["data"].(map[string]interface{}); found {
				secret.Data = data
			}
		}
	}

	return secret.Data, nil
}

func (c *Client) renew(client *api.Client, key string, secret *api.Secret) error {
	renewer, err := client.NewRenewer(&api.RenewerInput{
		Secret: secret,
	})

	if err != nil {
		return err
	}

	c.renewers = append(c.renewers, renewer)
	logger := log.WithField("secret", key)

	go func() {
		for {
			logger.Info("renewing")

			select {
			case _ = <-renewer.RenewCh():
				logger.Info("renewed")
			case err = <-renewer.DoneCh():
				logger.WithError(err).Info("renewing stopped")
				break
			}
		}
	}()

	go renewer.Renew()
	return nil
}
