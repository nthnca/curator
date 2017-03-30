package disk

import (
	"os"
	"path/filepath"
	"strings"
)

func List(root string) map[string]string {
	m := make(map[string]string)
	visit := func(path string, fs os.FileInfo, err error) error {
		base := filepath.Base(path)
		key := strings.TrimSuffix(base, ".jpg")
		if key == base {
			return nil
		}

		m[key] = path
		return nil
	}

	filepath.Walk(root, visit)
	return m
}
