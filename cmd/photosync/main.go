package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nthnca/curator/config"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func ReadAllPhotos() (map[string]string, error) {
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

	filepath.Walk(".", visit)
	return m, nil
}

func ListBucket() (map[string]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	m := make(map[string]string)
	bkt := client.Bucket(config.StorageBucket)
	for it := bkt.Objects(ctx, nil); ; {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}

		m[strings.SplitN(objAttrs.Name, ".", 2)[0]] = objAttrs.Name
	}
	return m, nil
}

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

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	/*
		var md map[string]string
		wg.Add(1)
		go func() {
			photoList, _ := client.LoadAllPhotos()
			for _, e := range photoList {
				md[e.GetName()] = e.GetName()
			}
			log.Printf("Photos in datastore: %v\n", len(md))
			wg.Done()
		}()
	*/

	var mf map[string]string
	wg.Add(1)
	go func() {
		mf, _ = ReadAllPhotos()
		log.Printf("Photos on disk: %v\n", len(mf))
		wg.Done()
	}()

	var mb map[string]string
	wg.Add(1)
	go func() {
		mb, _ = ListBucket()
		log.Printf("Photos in bucket: %v\n", len(mb))
		wg.Done()
	}()

	wg.Wait()

	unknown := false
	for key, _ := range mb {
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
}
