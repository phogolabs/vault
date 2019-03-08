package vault

import (
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
)

func splitBy(path, separator string) []string {
	parts := []string{}

	for _, part := range strings.Split(path, separator) {
		if part == "" {
			continue
		}

		parts = append(parts, part)
	}

	return parts
}

func abs(key string, mnt *api.MountOutput) string {
	switch mnt.Type {
	case "kv":
		if version, ok := mnt.Options["version"]; ok && version == "2" {
			dir, name := path.Split(key)
			key = path.Join(dir, "data", name)
		}
	}

	return key
}
