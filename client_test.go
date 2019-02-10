package vault_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	"github.com/hashicorp/vault/api"
	"github.com/phogolabs/vault"
)

var _ = Describe("Client", func() {
	var (
		client   *vault.Client
		server   *Server
		handlers []http.HandlerFunc
	)

	JustBeforeEach(func() {
		var err error

		server = NewServer()
		server.AppendHandlers(handlers...)

		config := api.DefaultConfig()
		config.Address = server.URL()

		client, err = vault.NewClient(config)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Auth", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandler(),
			}
		})

		AfterEach(func() {
			client.Stop()
		})

		It("authenticates the client", func() {
			Expect(client.Auth("kubernetes", nil)).To(Succeed())
		})

		Context("when the auth fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					newAuthHandlerFailed(),
					newAuthHandlerFailed(),
					newAuthHandlerFailed(),
				}
			})

			It("returns an error", func() {
				Expect(client.Auth("kubernetes", nil).Error()).To(ContainSubstring("Code: 500"))
			})
		})
	})

	Describe("Mount", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandler(),
				newGetMntHandler(),
			}
		})

		JustBeforeEach(func() {
			Expect(client.Auth("kubernetes", nil)).To(Succeed())
		})

		It("return the mount options", func() {
			mnt, err := client.Mount("/app/kv/config")
			Expect(err).NotTo(HaveOccurred())
			Expect(mnt).NotTo(BeNil())
		})

		Context("when getting the options fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					newAuthHandler(),
					newGetMntHandlerFailed(),
					newGetMntHandlerFailed(),
					newGetMntHandlerFailed(),
				}
			})

			It("returns an error", func() {
				mnt, err := client.Mount("/app/kv/config")
				Expect(mnt).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Code: 500"))
			})
		})

		Context("when decoding the response fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					newAuthHandler(),
					newGetMntHandlerBadResponse(),
				}
			})

			It("returns an error", func() {
				mnt, err := client.Mount("/app/kv/config")
				Expect(mnt).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("json: cannot unmarshal string into Go value of type map[string]interface {}"))
			})
		})
	})

	Describe("Secret", func() {
		BeforeEach(func() {
			handlers = []http.HandlerFunc{
				newAuthHandler(),
				newGetMntHandler(),
				newGetKVHandler(),
			}
		})

		JustBeforeEach(func() {
			Expect(client.Auth("kubernetes", nil)).To(Succeed())
		})

		It("return the secret successuflly", func() {
			secret, err := client.Secret("/app/kv/config")
			Expect(err).NotTo(HaveOccurred())
			Expect(secret).NotTo(BeNil())
			Expect(secret).To(HaveKeyWithValue("password", "swordfish"))
		})

		Context("when getting the secret fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					newAuthHandler(),
					newGetMntHandler(),
					newGetKVHandlerFailed(),
					newGetKVHandlerFailed(),
					newGetKVHandlerFailed(),
				}
			})

			It("returns an error", func() {
				mnt, err := client.Secret("/app/kv/config")
				Expect(mnt).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Code: 500"))
			})
		})

		Context("when decoding the secret fails", func() {
			BeforeEach(func() {
				handlers = []http.HandlerFunc{
					newAuthHandler(),
					newGetMntHandler(),
					newGetKVHandlerBadResponse(),
				}
			})

			It("returns an error", func() {
				mnt, err := client.Secret("/app/kv/config")
				Expect(mnt).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("json: cannot unmarshal string into Go value of type api.Secret"))
			})
		})
	})
})
