package newphotos

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
	"google.golang.org/api/iterator"
)

var (
	photoData []message.Photo
)

func Handler() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	meta := client.Bucket(config.MetadataStorageBucket)
	for it := meta.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects: %v", err)
		}

		rc, err := client.Bucket(config.MetadataStorageBucket).Object(obj.Name).NewReader(ctx)
		if err != nil {
			log.Fatalf("Failed to create reader: %v", err)
		}

		slurp, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		var ph message.Photo
		if err := proto.Unmarshal(slurp, &ph); err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		photoData = append(photoData, ph)
	}

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

			process(obj)
		}
	}
}

func getPhotoInfo(attr *storage.ObjectAttrs) *message.Photo {
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
		log.Fatalf("Failed to read file: %v", err)
	}

	md := md5.Sum([]byte(slurp))
	if !bytes.Equal(attr.MD5, md[:]) {
		log.Fatalf("File corrupted?")
	}

	sha := sha256.Sum256([]byte(slurp))
	// Is this photo already known?
	for _, pd := range photoData {
		if bytes.Equal(sha[:], pd.GetSha256Sum()) {
			return &pd
		}
	}

	err = ioutil.WriteFile(attr.Name, slurp, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	photo, err := util.IdentifyPhoto(attr.Name, attr.MD5, sha[:])
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	data, err := proto.Marshal(photo)
	if err != nil {
		log.Fatalf("Failed to marshal photo proto: %v", err)
	}

	log.Printf("Saving metadata\n")
	wc := client.Bucket(config.MetadataStorageBucket).Object(attr.Name + ".pb").NewWriter(ctx)
	foo := md5.Sum(data)
	wc.MD5 = foo[:]
	if _, err := wc.Write(data); err != nil {
		log.Fatalf("Failed writing meta data: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("Failed closing meta data write: %v", err)
	}
	return photo
}

func addPhotoToLongTerm(attr *storage.ObjectAttrs, photo *message.Photo) {
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
		src := client.Bucket(attr.Bucket).Object(attr.Name)
		dest := client.Bucket(longTerm).Object(pname)
		_, err = dest.CopierFrom(src).Run(ctx)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		return
	}
	log.Printf("Already have: %v/%v -> %v/%v\n", attr.Bucket, attr.Name, longTerm, pname)
}

// Download file
func process(attr *storage.ObjectAttrs) {
	photo := getPhotoInfo(attr)
	addPhotoToLongTerm(attr, photo)
}
