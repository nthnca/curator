package sync

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/disk"
	"github.com/nthnca/curator/data/gcs"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Should implement this to speed this up
// https://gobyexample.com/worker-pools
func StoreFile(path, key string) {
	tmpfile, err := ioutil.TempFile("", "image-")
	if err != nil {
		log.Fatalf("Creating tmp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	log.Printf("exec: convert %v -quality 80 -resize 1920x1080 %v",
		path, tmpfile.Name())
	cmd := exec.Command("convert", path, "-quality", "80", "-resize",
		"1920x1080", tmpfile.Name())
	if err := cmd.Run(); err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}

	dest := fmt.Sprintf("gs://%v/%v.jpg", config.StorageBucket, key)
	log.Printf("exec: gsutil -h Content-Type:image/jpg cp -a public-read %v %v",
		tmpfile.Name(), dest)
	cmd2 := exec.Command("gsutil", "-h", "Content-Type:image/jpeg",
		"cp", "-a", "public-read", tmpfile.Name(), dest)
	if err := cmd2.Run(); err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
}

func Handler(_ *kingpin.ParseContext) error {
	var wg sync.WaitGroup

	var mf map[string]string
	wg.Add(1)
	go func() {
		mf = disk.List(config.PhotoPath)
		log.Printf("Photos found in '%s': %v\n",
			config.PhotoPath, len(mf))
		wg.Done()
	}()

	var mb map[string]string
	wg.Add(1)
	go func() {
		mb, _ = gcs.List(config.StorageBucket)
		log.Printf("Photos in bucket: %v\n", len(mb))
		wg.Done()
	}()

	wg.Wait()

	unknown := false
	for key := range mb {
		if _, ok := mf[key]; ok {
			continue
		}
		log.Printf("Unknown file in storage bucket: %v\n", key)
		unknown = true
	}

	if unknown {
		log.Fatalf("Can't continue because of unknown files.")
	}

	for key, path := range mf {
		if _, ok := mb[key]; ok {
			continue
		}
		log.Printf("Key %v %v\n", key, path)
		StoreFile(path, key)
	}
	return nil
}
