package disk

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/util/need"
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

var needDataPhotoList need.NeedData

func NeedPhotoList() func() map[string]string {
	n := needDataPhotoList.Need(func() interface{} {
		d := List(config.PhotoPath)
		return d
	})
	return func() map[string]string {
		return n().(map[string]string)
	}
}
