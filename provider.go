package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/phogolabs/cli"
)

//go:generate counterfeiter -fake-name Fetcher -o ./fake/fetcher.go . Fetcher

// Fetcher fetches secrets
type Fetcher interface {
	// Secret returns the underlying secret
	Secret(path string) (interface{}, error)
	// Stop stops the fetcher
	Stop()
}

var (
	_ cli.Provider    = &Provider{}
	_ cli.Transaction = &Provider{}
)

// Provider is a parser that populates flags from Hashi Corp Vault
type Provider struct {
	Fetcher Fetcher
}

// Provide parses the args
func (m *Provider) Provide(ctx *cli.Context) (err error) {
	if err = m.init(ctx); err != nil {
		return err
	}

	if m.Fetcher == nil {
		return nil
	}

	var (
		path   string
		secret interface{}
	)

	for _, flag := range ctx.Command.Flags {
		accessor := &cli.FlagAccessor{Flag: flag}

		if path = accessor.MetaKey("vault_key"); path == "" {
			continue
		}

		if secret, err = m.Fetcher.Secret(path); err != nil {
			return err
		}

		if err = accessor.SetValue(secret); err != nil {
			return err
		}
	}

	return nil
}

func (m *Provider) init(ctx *cli.Context) error {
	if m.Fetcher != nil {
		return nil
	}

	address := ctx.String("vault-addr")
	if address == "" {
		return nil
	}

	config := api.DefaultConfig()
	config.Address = address

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	m.Fetcher = &RepositoryTree{
		Repository: client,
		Root:       make(map[string]interface{}),
	}

	if token := ctx.String("vault-token"); token != "" {
		client.SetToken(token)
		return nil
	}

	var (
		path   = ctx.String("vault-auth-mount-path")
		secret = map[string]interface{}{
			"role": ctx.String("vault-auth-role"),
			"jwt":  ctx.String("vault-auth-kube-jwt"),
		}
	)

	return client.Auth(path, secret)
}

// Rollback stops the provider
func (m *Provider) Rollback(ctx *cli.Context) error {
	if m.Fetcher != nil {
		m.Fetcher.Stop()
	}
	return nil
}
