package name

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

/*
func IsValidPhotoID(name string) bool {
	match, err := regexp.MatchString("^[0-9]+-[0-9]+[^.]+$", name)
	if err != nil {
		log.Fatalf("Regexp failed: %v", err)
	}
	return match
}
		if !dryRun {
			os.Rename(path, result)
		}
*/

func NamePhoto(path string) (string, error) {
	extension := filepath.Ext(path)
	basename := strings.TrimSuffix(filepath.Base(path), extension)

	match, err := regexp.MatchString("^[.](jpg|JPG)$", extension)
	if err != nil || !match {
		return "", fmt.Errorf("Invalid file extension: %s", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Failed to open file: %v", err)
	}
	defer f.Close()

	e, err := exif.Decode(f)
	if err != nil {
		return "", fmt.Errorf("Error decoding photo: %v", err)
	}

	get := func(f exif.FieldName) string {
		v, _ := e.Get(f)
		s, _ := v.StringVal()
		return s
	}

	cameraModels := map[string]string{
		"DMC-GF1": "GF1",
		"DMC-GX7": "GX7",
		"DMC-LX3": "LX3",
	}

	datetime := strings.Replace(get(exif.DateTimeOriginal), ":", "", -1)
	date := strings.Split(datetime, " ")[0]
	time := strings.Split(datetime, " ")[1]
	model, ok := cameraModels[get(exif.Model)]
	if !ok {
		return "", fmt.Errorf("Unknown camera: %s", get(exif.Model))
	}

	result := fmt.Sprintf("%s-%s-%s-%s.jpg", date, time, model, basename)

	return result, nil

}

func Handler(dest string, liveRun bool) {
	err := filepath.Walk(".",
		func(path string, fs os.FileInfo, err error) error {
			if path == "." {
				return nil
			}

			if fs.IsDir() {
				return filepath.SkipDir
			}

			name, err := NamePhoto(path)
			if err != nil {
				log.Printf("Error: %v", err)
				return nil
			}

			result := filepath.Join(dest, name[:4], name)
			if _, err := os.Stat(result); !os.IsNotExist(err) {
				log.Printf("File already exists: %v", result)
				return nil
			}

			if !liveRun {
				// log.Printf("(dry-run) cp %s %s", path, result)
				return nil
			}

			log.Printf("cp %s %s", path, result)
			os.Chmod(path, 0644)
			os.Link(path, result)
			return nil
		})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
