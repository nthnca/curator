package sync

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/client"
	"github.com/nthnca/curator/data/disk"
	"github.com/nthnca/curator/data/gcs"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"

	"github.com/nthnca/datastore"
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

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- *message.Photo) {
	defer wg.Done()
	for j := range jobs {
		photo, err := util.IdentifyPhoto(j)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		results <- photo
	}
}

func Handler(_ *kingpin.ParseContext) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		photoListGetter := disk.NeedPhotoList()
		storageListGetter := gcs.NeedStorageList()

		mf := photoListGetter()
		log.Printf("Photos found in '%s': %v\n",
			config.PhotoPath, len(mf))
		mb := storageListGetter()
		log.Printf("Photos in bucket: %v\n", len(mb))

		for key := range mb {
			if _, ok := mf[key]; ok {
				continue
			}

			log.Printf("Unknown file in storage bucket: %v\n", key)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		photoListGetter := disk.NeedPhotoList()
		storageListGetter := gcs.NeedStorageList()

		mf := photoListGetter()
		mb := storageListGetter()

		for key, path := range mf {
			if _, ok := mb[key]; ok {
				continue
			}
			log.Printf("Key %v %v\n", key, path)
			StoreFile(path, key)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clt, _ := datastore.NewCloudClient(config.ProjectID)
		photoListGetter := disk.NeedPhotoList()
		photoCacheGetter := client.NeedPhotoCache(clt)

		mf := photoListGetter()
		mp := photoCacheGetter()
		log.Printf("Photos in database: %v\n", len(mp))

		if len(mp) == len(mf) {
			log.Printf("You should sync the database")
			return
		}

		jobs := make(chan string, 100)
		results := make(chan *message.Photo, 100)

		wg2 := &sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			wg2.Add(1)
			go worker(wg2, jobs, results)
		}

		go func() {
			wg2.Wait()
			close(results)
		}()

		go func() {
			for _, path := range mf {
				jobs <- path

			}
			close(jobs)
		}()

		set := message.PhotoSet{}
		for r := range results {
			log.Printf("%v", r.GetKey())
			set.Photo = append(set.Photo, r)

			/*
				if len(set.Photo) > 1000 {
					_, err = client.CreatePhotoCache(clt, &set)
					if err != nil {
						log.Fatalf("Foo: %v", err)
					}
					set = message.PhotoSet{}
				}
			*/
		}
		/*
			_, err = client.CreatePhotoCache(clt, &set)
			if err != nil {
				log.Fatalf("Foo: %v", err)
			}
		*/
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		clt, _ := datastore.NewCloudClient(config.ProjectID)
		photoCacheGetter := client.NeedPhotoCache(clt)

		mp := photoCacheGetter()

		for _, photo := range mp {
			log.Printf("%v", photo.GetKey())
			/*
				client.UpdatePhoto(clt, photo.GetKey(), photo)
			*/
		}
	}()

	wg.Wait()
	return nil
}
