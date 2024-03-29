package manager

import (
	"os"
	"path/filepath"
	"strings"
)

func ExpandPath(name string) string {
	if strings.HasPrefix(name, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return name
		}

		return strings.Replace(name, "~", homeDir, 1)
	}

	return name
}

func resolvePath(root, name string) string {
	baseRoot := ExpandPath(name)
	if baseRoot == "" || !filepath.IsAbs(baseRoot) {
		baseRoot = filepath.Join(root, name)
	}
	return baseRoot
}
