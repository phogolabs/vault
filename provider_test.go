package vault_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/vault"
	"github.com/phogolabs/vault/fake"
)

var _ = Describe("Provider", func() {
	var (
		provider *vault.Provider
		server   *Server
		ctx      *cli.Context
		handlers []http.HandlerFunc
	)

	BeforeEach(func() {
		handlers = []http.HandlerFunc{
			newAuthHandler(),
			newGetMntHandler(),
			newGetKVHandler(),
		}

		provider = &vault.Provider{}
	})

	JustBeforeEach(func() {
		server = NewServer()
		server.AppendHandlers(handlers...)

		ctx = &cli.Context{
			Command: &cli.Command{
				Name: "app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "vault-addr",
						Value: server.URL(),
					},
					&cli.StringFlag{
						Name:  "vault-auth-mount",
						Value: "kubernetes",
					},
					&cli.StringFlag{
						Name:  "vault-auth-role",
						Value: "admin",
					},
					&cli.StringFlag{
						Name:  "vault-auth-kube-jwt",
						Value: "kubo",
					},
					&cli.StringFlag{
						Name: "password",
						Metadata: map[string]string{
							"vault_key": "/app/kv/config::password",
						},
					},
				},
			},
		}
	})

	AfterEach(func() {
		Expect(provider.Rollback(ctx)).To(Succeed())

		server.Close()
	})

	It("parses the flags successfully", func() {
		Expect(provider.Provide(ctx)).To(Succeed())
		Expect(ctx.String("password")).To(Equal("swordfish"))
	})

	Context("when the token is provided", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newGetMntHandler(),
				newGetKVHandler(),
			}
		})

		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "vault-addr",
							Value: server.URL(),
						},
						&cli.StringFlag{
							Name:  "vault-token",
							Value: "my-token",
						},
						&cli.StringFlag{
							Name: "password",
							Metadata: map[string]string{
								"vault_key": "/app/kv/config::password",
							},
						},
					},
				},
			}
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(ctx.String("password")).To(Equal("swordfish"))
		})
	})

	Context("when the authentication fails", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandlerFailed(),
				newAuthHandlerFailed(),
				newAuthHandlerFailed(),
			}
		})

		JustBeforeEach(func() {
			server.SetAllowUnhandledRequests(true)
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx).Error()).To(ContainSubstring("Code: 500"))
		})
	})

	Context("when the fetcher fails", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandler(),
				newGetMntHandlerFailed(),
				newGetMntHandlerFailed(),
				newGetMntHandlerFailed(),
			}
		})

		JustBeforeEach(func() {
			server.SetAllowUnhandledRequests(true)
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx).Error()).To(ContainSubstring("Code: 500"))
		})
	})

	Context("when setting the flag fails", func() {
		JustBeforeEach(func() {
			flags := ctx.Command.Flags
			flags[len(flags)-1] = &cli.IntFlag{
				Name: "password",
				Metadata: map[string]string{
					"vault_key": "/app/kv/config::password",
				},
			}
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx)).To(MatchError("strconv.ParseInt: parsing \"swordfish\": invalid syntax"))
		})
	})

	Context("when the fetcher is already initialized", func() {
		var fetcher *fake.Fetcher

		BeforeEach(func() {
			fetcher = &fake.Fetcher{}
			fetcher.SecretReturns("swordfish", nil)

			provider.Fetcher = fetcher
		})

		AfterEach(func() {
			provider.Fetcher = nil
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(ctx.String("password")).To(Equal("swordfish"))
		})
	})

	Context("when the fetcher cannot be initialized", func() {
		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "password",
							Metadata: map[string]string{
								"vault_key": "/app/kv/config::password",
							},
						},
					},
				},
			}
		})

		It("parses the flags successfully", func() {
			Expect(provider.Provide(ctx)).To(Succeed())
			Expect(provider.Fetcher).To(BeNil())
			Expect(ctx.String("password")).To(BeEmpty())
		})
	})

	Context("when the valut-addr is not valid", func() {
		JustBeforeEach(func() {
			ctx = &cli.Context{
				Command: &cli.Command{
					Name: "app",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "vault-addr",
							Value: "://address",
						},
						&cli.StringFlag{
							Name: "password",
							Metadata: map[string]string{
								"vault_key": "/app/kv/config::password",
							},
						},
					},
				},
			}
		})

		It("returns an error", func() {
			Expect(provider.Provide(ctx)).To(MatchError("parse ://address: missing protocol scheme"))
		})
	})
})
