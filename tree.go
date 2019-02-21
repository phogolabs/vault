package vault

import (
	"path/filepath"
)

//go:generate counterfeiter -fake-name Repository -o ./fake/repository.go . Repository

// Repository is a repository of the secrets
type Repository interface {
	// Secret provides the secret from the backend
	Secret(path string) (map[string]interface{}, error)
	// Stop stops the repository
	Stop()
}

var _ Repository = &RepositoryTree{}

// RepositoryTree caches the secrets
type RepositoryTree struct {
	Repository Repository
	Root       map[string]map[string]interface{}
}

// Secret returns value from a tree
func (r *RepositoryTree) Secret(path string) (map[string]interface{}, error) {
	path = filepath.Join(split(path)...)
	node, found := r.Root[path]

	if !found {
		var err error

		if node, err = r.Repository.Secret(path); err != nil {
			return nil, err
		}
	}

	r.Root[path] = node
	return node, nil
}

// Stop stops the tree
func (r *RepositoryTree) Stop() {
	r.Repository.Stop()
}
