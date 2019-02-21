package vault_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/vault"
	"github.com/phogolabs/vault/fake"
)

var _ = Describe("RepositoryTree", func() {
	var (
		repository *fake.Repository
		tree       *vault.RepositoryTree
	)

	BeforeEach(func() {
		repository = &fake.Repository{}

		tree = &vault.RepositoryTree{
			Repository: repository,
			Root:       make(map[string]map[string]interface{}),
		}
	})

	AfterEach(func() {
		tree.Stop()
		Expect(repository.StopCallCount()).To(Equal(1))
	})

	fetch := func(config map[string]interface{}) map[string]map[string]interface{} {
		return map[string]map[string]interface{}{
			"app/service-api/config": config,
		}
	}

	ItReturnsTheSecretSuccessfully := func(count int) {
		It("returns the secret successfully", func() {
			secret, err := tree.Secret("/app/service-api/config")
			Expect(err).To(BeNil())
			Expect(secret).To(HaveKeyWithValue("password", "swordfish"))
			Expect(repository.SecretCallCount()).To(Equal(count))
		})
	}

	Context("when the secrets is not cached", func() {
		BeforeEach(func() {
			secrets := map[string]interface{}{
				"password": "swordfish",
			}

			repository.SecretReturns(secrets, nil)
		})

		ItReturnsTheSecretSuccessfully(1)
	})

	Context("when the secret is already fetched", func() {
		BeforeEach(func() {
			secrets := map[string]interface{}{
				"password": "swordfish",
			}

			tree.Root = fetch(secrets)
		})

		ItReturnsTheSecretSuccessfully(0)
	})

	Context("when the provider fails", func() {
		BeforeEach(func() {
			repository.SecretReturns(nil, fmt.Errorf("oh no!"))
		})

		It("returns an error", func() {
			secret, err := tree.Secret("/app/service-api/config")
			Expect(err).To(MatchError("oh no!"))
			Expect(secret).To(BeNil())
			Expect(tree.Root).NotTo(HaveKey("/app/service-api/config"))
		})
	})
})
