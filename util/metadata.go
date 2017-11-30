package util

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/message"
	"google.golang.org/api/iterator"
)

const (
	masterFileName = "master.md"
)

type PhotoInfo struct {
	Data []*message.Photo
}

func (md *PhotoInfo) Lookup(sha256 []byte, photo *message.Photo) bool {
	for _, iter := range md.Data {
		if bytes.Equal(sha256[:], iter.GetSha256Sum()) {
			*photo = *iter
			return true
		}
	}
	return false
}

func (md *PhotoInfo) insert(photo *message.Photo) bool {
	var devnull message.Photo
	if !md.Lookup(photo.GetSha256Sum(), &devnull) {
		md.Data = append(md.Data, photo)
		return true
	}
	return false
}

func (md *PhotoInfo) insertAll(photos []*message.Photo) bool {
	yep := false
	for _, iter := range photos {
		yep = md.insert(iter) || yep
	}
	return yep
}

func (md *PhotoInfo) Save(ctx context.Context, client *storage.Client, photo *message.Photo) {
	set := message.PhotoSet{}
	set.Photo = append(set.Photo, photo)

	name := fmt.Sprintf("file.%s.%d.md", hex.EncodeToString(photo.GetSha256Sum()), time.Now().Unix())
	md.save(ctx, client, name, &set)
}

func (md *PhotoInfo) SaveAll(ctx context.Context, client *storage.Client, photos []*message.Photo) {
	set := message.PhotoSet{}
	for _, photo := range photos {
		set.Photo = append(set.Photo, photo)
	}

	md.save(ctx, client, masterFileName, &set)
}

func (md *PhotoInfo) save(ctx context.Context, client *storage.Client, filename string, photos *message.PhotoSet) {
	data, err := proto.Marshal(photos)
	if err != nil {
		log.Fatalf("Failed to marshal photo set proto: %v", err)
	}

	wc := client.Bucket(config.Get().PhotoInfoBucket).Object(filename).NewWriter(ctx)
	checksum := md5.Sum(data)
	wc.MD5 = checksum[:]
	if _, err := wc.Write(data); err != nil {
		log.Fatalf("Failed writing meta data: %v", err)
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("Failed closing meta data write: %v", err)
	}
}

func (md *PhotoInfo) Load(ctx context.Context, client *storage.Client) {
	ch := make(chan PhotoInfoFileInfo)
	ps := loadPhotoInfo(ctx, client, config.Get().PhotoInfoBucket, masterFileName)
	go loadExtras(ctx, client, ch)
	md.Data = append(md.Data, ps.Photo...)
	needSave := false
	for out := range ch {
		if !md.insertAll(out.photo.Photo) {
			file := *(out.file)
			go func() {
				if err := client.Bucket(file.Bucket).Object(file.Name).Delete(ctx); err != nil {
					log.Fatalf("Failed to delete: %v", err)
				}
			}()

		} else {
			needSave = true
		}
	}
	log.Printf("Loaded %d PhotoInfo entries.", len(md.Data))
	if needSave {
		md.SaveAll(ctx, client, md.Data)
	}
}

type PhotoInfoFileInfo struct {
	photo *message.PhotoSet
	file  *storage.ObjectAttrs
}

func loadPhotoInfo(ctx context.Context, client *storage.Client, bucketname, filename string) *message.PhotoSet {

	reader, err := client.Bucket(bucketname).Object(filename).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	slurp, err := ioutil.ReadAll(reader)
	reader.Close()
	if err != nil {
		log.Fatalf("Failed to read PhotoInfo: %v", err)
	}

	var ps message.PhotoSet
	if err := proto.Unmarshal(slurp, &ps); err != nil {
		log.Fatalf("Failed to decode proto: %v", err)
	}
	return &ps
}

func loadExtras(ctx context.Context, client *storage.Client, ch chan<- PhotoInfoFileInfo) {
	var wg sync.WaitGroup
	abc := make(chan *storage.ObjectAttrs)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for obj := range abc {
				ps := loadPhotoInfo(ctx, client, obj.Bucket, obj.Name)
				ch <- PhotoInfoFileInfo{ps, obj}
			}
		}()
	}

	meta := client.Bucket(config.Get().PhotoInfoBucket)
	for it := meta.Objects(ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects: %v", err)
		}
		if obj.Name == masterFileName {
			continue
		}

		abc <- obj
	}
	close(abc)
	wg.Wait()
	close(ch)
}
