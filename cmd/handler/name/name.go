package name

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
)

var cameraModels = make(map[string]string)

func initCameraModels() {
	cameraModels = map[string]string{
		"DMC-GF1": "GF1",
		"DMC-GX7": "GX7",
		"DMC-LX3": "LX3",
	}
}

func isValidPhotoName(path string) bool {
	match, err := regexp.MatchString(".*[.](jpg|JPG)$", path)
	return err == nil && match
}

// getPhotoName creates a canonical name for a photo based on its original
// name plus EXIF data.
// TODO: Switch this to using ImageMagick
func getPhotoName(path string) (string, error) {
	extension := filepath.Ext(path)
	basename := strings.TrimSuffix(filepath.Base(path), extension)

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

func generateFileNames(foo chan<- string) error {
	return filepath.Walk(".",
		func(path string, fs os.FileInfo, err error) error {
			if path == "." {
				return nil
			}

			if fs.IsDir() {
				return filepath.SkipDir
			}

			foo <- path

			return nil
		})
}

func processFiles(f func(path string) error) bool {
	files := make(chan string)
	errors := false

	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for path := range files {
				if err := f(path); err != nil {
					log.Printf("%v", err)
					errors = true
				}
			}
			wg.Done()
		}()
	}

	generateFileNames(files)
	close(files)
	wg.Wait()

	return errors
}

type copy struct {
	dest   string
	doCopy bool
}

func (c *copy) checkAndCopy(path string) error {
	if !isValidPhotoName(path) {
		return fmt.Errorf("Skipping file: %v", path)
	}

	name, err := getPhotoName(path)
	if err != nil {
		return fmt.Errorf("Failed to determine name: %v", err)
	}

	if !c.doCopy {
		return nil
	}

	result := filepath.Join(c.dest, name[:4], name)
	if _, err := os.Stat(result); !os.IsNotExist(err) {
		return fmt.Errorf("File already exists: %v", result)
	}

	// log.Printf("cp %s %s", path, result)
	if err := os.Chmod(path, 0644); err != nil {
		return err
	}

	if err := os.Link(path, result); err != nil {
		return err
	}

	return nil
}

func Handler(dest string, liveRun bool) {
	log.SetFlags(0)
	initCameraModels()
	c := &copy{dest: dest}
	if processFiles(c.checkAndCopy) {
		log.Printf("Didn't copy files, fix errors and re-run")
		return
	}

	c.doCopy = true
	processFiles(c.checkAndCopy)
}
