package vault

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
)

func base(path string) string {
	parts := strings.Split(path, "::")
	parts = split(parts[0])
	return filepath.Join(parts...)
}

func split(path string) []string {
	parts := []string{}

	for _, part := range strings.Split(path, "/") {
		if part == "" {
			continue
		}

		elements := strings.Split(part, "::")
		parts = append(parts, elements...)
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
