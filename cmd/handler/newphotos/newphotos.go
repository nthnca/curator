package newphotos

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/mediainfo/exif"
	"github.com/nthnca/curator/mediainfo/store"
	"github.com/nthnca/curator/util"
	"google.golang.org/api/iterator"
)

const (
	dryRun = false
)

var (
	photoData util.PhotoInfo
	ctx       context.Context
	client    *storage.Client
	mediaInfo *store.MediaInfo
)

func Handler() {
	ctx = context.Background()
	err := fmt.Errorf("") // Next line can't use :=
	client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	t := time.Now().UnixNano()
	log.Printf("Reading PhotoInfo: %s", config.PhotoInfoBucket())
	mi, err := store.New(ctx, client, config.PhotoInfoBucket())
	mediaInfo = mi
	if err != nil {
		log.Fatalf("NewMediaInfo failed: %v", err)
	}
	log.Printf("Count: %d", len(mediaInfo.All()))
	log.Printf("PhotoInfo read took %v seconds",
		float64(time.Now().UnixNano()-t)/1000000000.0)

	c := 0
	log.Printf("Looking for photos in: %s", config.PhotoQueueBucket())
	for it := client.Buckets(ctx, config.PhotoQueueBucket()); ; {
		bktiter, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through buckets: %v", err)
		}

		log.Printf("Processing files in bucket: %s", bktiter.Name)
		var set []*storage.ObjectAttrs
		bkt := client.Bucket(bktiter.Name)
		for it := bkt.Objects(ctx, nil); ; {
			obj, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate through objects: %v", err)
			}

			if len(set) > 0 && base(set[0].Name) != base(obj.Name) {
				processPhotoSet(set)
				set = set[:0]
			}
			set = append(set, obj)
			c++
			if c > 5 {
				log.Fatalf("Foo")
			}
		}
		if len(set) > 0 {
			processPhotoSet(set)
		}
	}
}

func base(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[:i]
	}
	return a
}

func suffix(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[i+1:]
	}
	return ""
}

func processPhotoSet(attr []*storage.ObjectAttrs) {
	media, err := convertToMedia(attr)
	if err != nil {
		log.Printf("convertToMedia: %v", err)
	}

	/*
		// Are these files already known?
		for _, k := range media.File {
			if true {
				fmt.Printf("%v\n", k)
			}
		}
	*/

	// cp files
	for _, a := range attr {
		name := getFileNameForAttr(a, media)
		_, err := client.Bucket(config.PhotoStorageBucket()).Object(name).Attrs(ctx)
		if err == nil {
			log.Printf("File already exists: %v", name)
			continue
		}
		if err != nil && err != storage.ErrObjectNotExist {
			log.Printf("File exists?, Error: %v", name)
			continue
		}
		log.Printf("Copying: %v/%v -> %v/%v\n",
			a.Bucket, a.Name, config.PhotoStorageBucket(), name)
		if !dryRun {
			src := client.Bucket(a.Bucket).Object(a.Name)
			dest := client.Bucket(config.PhotoStorageBucket()).Object(name)
			_, err = dest.CopierFrom(src).Run(ctx)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}

	// save meta
	mediaInfo.Insert(ctx, client, media)

	// delete files
	for _, a := range attr {
		log.Printf("Deleting: %v/%v\n", a.Bucket, a.Name)
		if !dryRun {
			if err := client.Bucket(a.Bucket).Object(a.Name).Delete(ctx); err != nil {
				log.Fatalf("Failed to delete: %v", err)
			}
		}
	}
}

func getFileNameForAttr(attr *storage.ObjectAttrs, m *message.Media) string {
	for _, mf := range m.File {
		if MD5(mf.Md5Sum) == MD5(attr.MD5) {
			return hex.EncodeToString(mf.Sha256Sum)
		}
	}
	log.Fatalf("Unable to find file")
	return ""
}

func convertToMedia(attr []*storage.ObjectAttrs) (*message.Media, error) {
	var jpg []*storage.ObjectAttrs
	var other []*storage.ObjectAttrs
	for _, a := range attr {
		if strings.ToLower(suffix(a.Name)) == "jpg" {
			jpg = append(jpg, a)
		} else {
			other = append(other, a)
		}
	}

	if len(jpg) != 1 {
		return nil, fmt.Errorf("Invalid set of JPGs found: %v", jpg)
	}

	var media message.Media
	jpginfo, err := getFile(jpg[0], "tmpfile.jpg")
	if err != nil {
		log.Fatalf("getfile %v", err)
	}
	media.File = append(media.File, jpginfo)

	mediainfo, err := exif.Parse("tmpfile.jpg")
	if err != nil {
		log.Fatalf("Oop %v", err)
	}
	media.Photo = mediainfo

	for _, a := range other {
		info, err := getFile(a, "")
		if err != nil {
			log.Fatalf("Oop 2 %v", err)
		}
		media.File = append(media.File, info)
	}

	media.Key = media.File[0].Sha256Sum
	media.TimestampSecondsSinceEpoch = time.Now().Unix()

	return &media, nil
}

func getFile(attr *storage.ObjectAttrs, localPath string) (*message.FileInfo, error) {
	rc, err := client.Bucket(attr.Bucket).Object(attr.Name).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create reader: %v", err)
	}

	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	md := md5.Sum([]byte(slurp))
	if !bytes.Equal(attr.MD5, md[:]) {
		return nil, fmt.Errorf("MD5 sum didn't match, file corrupted?")
	}

	sha := sha256.Sum256([]byte(slurp))
	sub := strings.Split(attr.Name, "/")
	name := sub[len(sub)-1]

	if localPath != "" {
		err = ioutil.WriteFile(localPath, slurp, 0644)
		if err != nil {
			return nil, fmt.Errorf("Failed to write file: %v", err)
		}
	}

	return &message.FileInfo{
		Filename:    name,
		Md5Sum:      md[:],
		Sha256Sum:   sha[:],
		SizeInBytes: attr.Size,
	}, nil
}

/*
	if err := os.Remove(name); err != nil {
		log.Fatalf("Attempting to remove file: %v", err)
	}
func removePhotoFromQueue(attr *storage.ObjectAttrs) {
	log.Printf("Deleting")
	if !dryRun {
		if err := client.Bucket(attr.Bucket).Object(attr.Name).Delete(ctx); err != nil {
			log.Fatalf("Failed to delete: %v", err)
		}
	}
}
*/

func MD5(a []byte) [16]byte {
	var rv [16]byte
	if len(a) != 16 {
		log.Fatalf("Byte array has incorrect size for MD5 (%d)", len(a))
	}

	copy(rv[:], a)
	return rv
}

func Sha256(a []byte) [32]byte {
	var rv [32]byte
	if len(a) != 32 {
		log.Fatalf("Byte array has incorrect size for SHA256 (%d)", len(a))
	}

	copy(rv[:], a)
	return rv
}
