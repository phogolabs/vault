package vault

import (
	"fmt"
	"path/filepath"
)

//go:generate counterfeiter -fake-name Repository -o ./fake/repository.go . Repository

// Repository is a repository of the secrets
type Repository interface {
	// Secret provides the secret from the backend
	Secret(path string) (map[string]interface{}, error)
}

var _ Fetcher = &RepositoryTree{}

// RepositoryTree caches the secrets
type RepositoryTree struct {
	Repository Repository
	Root       map[string]interface{}
}

// Secret returns value from a tree
func (r *RepositoryTree) Secret(path string) (interface{}, error) {
	secret := base(path)
	parts := split(path)
	path = filepath.Join(parts...)

	var (
		err     error
		current string
		parent  = r.Root
	)

	for _, part := range parts {
		current = filepath.Join(current, part)
		node, found := parent[part]

		if !found {
			if current == secret {
				if node, err = r.Repository.Secret(secret); err != nil {
					return nil, err
				}
			} else {
				node = make(map[string]interface{})
			}
		}

		parent[part] = node

		if current == path {
			return node, nil
		}

		next, ok := node.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("vault: invalid type '%T' for path '%v'", node, current)
		}

		parent = next
	}

	return nil, fmt.Errorf("vault: path '%s' not found", path)
}
