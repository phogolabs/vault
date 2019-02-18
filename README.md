# vault

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A package that facilitates working with http://vaultproject.io/ in context of
[CLI](https://github.com/phogolabs/cli). It increases the security of Golang
applications by populating a command line arguments from the vault.

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install vault, simply run:
```
$ go get github.com/phogolabs/vault
```

## Getting Started

In order to have the provider enabled, you need to set its token either
directly or authenticating the client with Kuberenetes. For that purpose, you
will need to set the following flags in your application:

```golang
import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/vault"
)

func main() {
	app := &cli.App{
		Name:      "prana",
		HelpName:  "prana",
		Usage:     "Golang Database Manager",
		UsageText: "prana [global options]",
		Version:   "1.0-beta-04",
		Action:    run,
		Providers: []cli.Provider{
			&vault.Provider{},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "vault-token",
				Usage:  "Hashi Corp Vault Token",
				EnvVar: "VAULT_TOKEN",
			},
			&cli.StringFlag{
				Name:   "vault-addr",
				Usage:  "Hashi Corp Vault Address",
				EnvVar: "VAULT_ADDR",
			},
			&cli.StringFlag{
				Name:   "vault-auth-mount",
				Usage:  "Hashi Corp Vault Auth Mount",
				EnvVar: "VAULT_AUTH_MOUNT",
				Value:  "kubernetes",
			},
			&cli.StringFlag{
				Name:   "vault-auth-role",
				Usage:  "Hashi Corp Vault Auth Role",
				EnvVar: "VAULT_AUTH_ROLE",
				Value:  "demo",
			},
			&cli.StringFlag{
				Name:   "vault-auth-kube-jwt",
				Usage:  "Hashi Corp Vault Kube Jwt",
				EnvVar: "VAULT_AUTH_KUBE_TOKEN",
			},
			&cli.StringFlag{
				Name:   "config",
				Usage:  "Aplication's config",
				EnvVar: "APP_CONFIG",
				Metadata: map[string]string{
					"vault_key": "/app/service-api/kv/config",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}
```

As you can see in order to match the flag with a given secret you should set
the `vault_key` in the meta data map.

Let's assume that we have the following JSON in your KV config:

```json
{
  "username": "root",
  "password": "swordfish"
}
```

If you want to populate a flag's value with the password field you should use
the following syntax for your `vault_key`:

```
/app/service-api/kv/config::password
```

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/vault).

[travis-img]: https://travis-ci.org/phogolabs/vault.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/vault
[report-img]: https://goreportcard.com/badge/github.com/phogolabs/vault
[report-url]: https://goreportcard.com/report/github.com/phogolabs/vault
[codecov-url]: https://codecov.io/gh/phogolabs/vault
[codecov-img]: https://codecov.io/gh/phogolabs/vault/branch/master/graph/badge.svg
[godoc-url]: https://godoc.org/github.com/phogolabs/vault
[godoc-img]: https://godoc.org/github.com/phogolabs/vault?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
