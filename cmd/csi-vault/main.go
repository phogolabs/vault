package main

import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/vault/driver"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "endpoint",
		Usage:    "CSI Unix Socket Endpoint",
		EnvVar:   "VAULT_CSI_ENDPOINT",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "node-id",
		Usage:    "Identifier of the node running the driver",
		EnvVar:   "VAULT_CSI_NODE_ID",
		Required: true,
	},
	&cli.StringFlag{
		Name:   "vault-token",
		Usage:  "Hashi Corp Vault Token",
		EnvVar: "VAULT_CSI_TOKEN",
	},
	&cli.StringFlag{
		Name:   "vault-addr",
		Usage:  "Hashi Corp Vault Address",
		EnvVar: "VAULT_CSI_ADDR",
	},
	&cli.StringFlag{
		Name:   "vault-auth-mount-path",
		Usage:  "Hashi Corp Vault Auth Mount",
		EnvVar: "VAULT_CSI_AUTH_MOUNT_PATH",
		Value:  "kubernetes",
	},
	&cli.StringFlag{
		Name:   "vault-auth-role",
		Usage:  "Hashi Corp Vault Auth Role",
		EnvVar: "VAULT_CSI_AUTH_ROLE",
	},
	&cli.StringFlag{
		Name:   "vault-auth-kube-jwt",
		Usage:  "Hashi Corp Vault Kube Jwt",
		EnvVar: "VAULT_CSI_AUTH_KUBE_TOKEN",
	},
}

func main() {
	app := &cli.App{
		Name:      "csi-vault",
		HelpName:  "csi-vault",
		Usage:     "Vault Container Storage Interface",
		UsageText: "csi-vault [global options]",
		Version:   "1.0-beta-01",
		Flags:     flags,
		Action:    run,
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("exited")
	}
}

func run(ctx *cli.Context) error {
	config := &driver.Config{
		Endpoint: ctx.String("endpoint"),
		Node:     ctx.String("node-id"),
	}

	service, err := driver.New(config)
	if err != nil {
		return err
	}

	return service.Run()
}
