package newphotos

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"google.golang.org/api/iterator"
)

const (
	dryRun = false
)

var (
	photoData util.PhotoInfo
)

func Handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	t := time.Now().UnixNano()
	photoData.Load(ctx, client)
	log.Printf("PhotoInfo read took %v seconds",
		float64(time.Now().UnixNano()-t)/1000000000.0)

	/*
		for _, iter := range photoData.Data {
			fmt.Printf("%s\n", hex.EncodeToString(iter.GetSha256Sum()))
		}
	*/

	for it := client.Buckets(ctx, config.PhotoQueueProject); ; {
		bktiter, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through buckets: %v", err)
		}

		bkt := client.Bucket(bktiter.Name)
		for it := bkt.Objects(ctx, nil); ; {
			obj, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate through objects: %v", err)
			}

			photo, haveInfo := getPhotoInfo(obj)
			havePhoto := false
			if photo != nil {
				havePhoto = addPhotoToLongTerm(obj, photo)
			}

			if haveInfo && havePhoto {
				removePhotoFromQueue(obj)
			}
		}
	}
}

func getPhotoInfo(attr *storage.ObjectAttrs) (*message.Photo, bool) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	rc, err := client.Bucket(attr.Bucket).Object(attr.Name).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		log.Fatalf("Failed to read photo: %v", err)
	}

	md := md5.Sum([]byte(slurp))
	if !bytes.Equal(attr.MD5, md[:]) {
		log.Fatalf("File corrupted?")
	}

	var pd message.Photo
	sha := sha256.Sum256([]byte(slurp))
	if found := photoData.Lookup(sha[:], &pd); found {
		return &pd, true
	}

	err = ioutil.WriteFile(attr.Name, slurp, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	photo, err := util.IdentifyPhoto(attr.Name, attr.MD5, sha[:])
	if err != nil {
		log.Printf("Attempting to identify file: %v", err)
		return nil, false
	}

	if err := os.Remove(attr.Name); err != nil {
		log.Fatalf("Attempting to remove file: %v", err)
	}

	log.Printf("Saving PhotoInfo\n")
	if !dryRun {
		photoData.Save(ctx, client, photo)
	}
	return photo, false
}

func addPhotoToLongTerm(attr *storage.ObjectAttrs, photo *message.Photo) bool {
	longTerm := config.PhotoStorageBucket
	pname := photo.GetPath()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Bucket(longTerm).Object(pname).Attrs(ctx)
	if err != nil {
		log.Printf("Copying: %v/%v -> %v/%v\n", attr.Bucket, attr.Name, longTerm, pname)
		if !dryRun {
			src := client.Bucket(attr.Bucket).Object(attr.Name)
			dest := client.Bucket(longTerm).Object(pname)
			_, err = dest.CopierFrom(src).Run(ctx)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
		return false
	}
	log.Printf("Already have: %v/%v -> %v/%v\n", attr.Bucket, attr.Name, longTerm, pname)
	return true
}

func removePhotoFromQueue(attr *storage.ObjectAttrs) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Printf("Deleting")
	if !dryRun {
		if err := client.Bucket(attr.Bucket).Object(attr.Name).Delete(ctx); err != nil {
			log.Fatalf("Failed to delete: %v", err)
		}
	}
}
