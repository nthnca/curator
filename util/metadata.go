package util

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/message"
	"google.golang.org/api/iterator"
)

type Metadata struct {
	Data []*message.Photo
}

func (md *Metadata) Lookup(sha256 []byte, photo *message.Photo) bool {
	for _, iter := range md.Data {
		if bytes.Equal(sha256[:], iter.GetSha256Sum()) {
			*photo = *iter
			return true
		}
	}
	return false
}

func (md *Metadata) Save(ctx context.Context, client *storage.Client, photo *message.Photo) {
	set := message.PhotoSet{}
	set.Photo = append(set.Photo, photo)

	md.save(ctx, client, "file."+photo.GetKey()+".md", &set)
}

func (md *Metadata) SaveAll(ctx context.Context, client *storage.Client, photos []*message.Photo) {
	set := message.PhotoSet{}
	for _, photo := range photos {
		set.Photo = append(set.Photo, photo)
	}

	t := time.Now().Unix()
	md.save(ctx, client, fmt.Sprintf("set.%d.md", t), &set)
}

func (md *Metadata) save(ctx context.Context, client *storage.Client, filename string, photos *message.PhotoSet) {
	data, err := proto.Marshal(photos)
	if err != nil {
		log.Fatalf("Failed to marshal photo set proto: %v", err)
	}

	wc := client.Bucket(config.MetadataStorageBucket).Object(filename).NewWriter(ctx)
	checksum := md5.Sum(data)
	wc.MD5 = checksum[:]
	if _, err := wc.Write(data); err != nil {
		log.Fatalf("Failed writing meta data: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("Failed closing meta data write: %v", err)
	}
}

func (md *Metadata) Load(ctx context.Context, client *storage.Client) {
	meta := client.Bucket(config.MetadataStorageBucket)
	for it := meta.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects: %v", err)
		}

		reader, err := client.Bucket(config.MetadataStorageBucket).Object(obj.Name).NewReader(ctx)
		if err != nil {
			log.Fatalf("Failed to create reader: %v", err)
		}

		slurp, err := ioutil.ReadAll(reader)
		reader.Close()
		if err != nil {
			log.Fatalf("Failed to read metadata: %v", err)
		}

		var ps message.PhotoSet
		if err := proto.Unmarshal(slurp, &ps); err != nil {
			log.Fatalf("Failed to decode proto: %v", err)
		}
		md.Data = append(md.Data, ps.Photo...)
	}

	log.Printf("Loaded %d metadata entries.", len(md.Data))

	//	md.SaveAll(ctx, client, md.Data)
}
