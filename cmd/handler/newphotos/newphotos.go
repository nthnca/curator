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
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/exif"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/mediainfo/store"
	"github.com/nthnca/curator/pkg/util"
	"google.golang.org/api/iterator"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	ctx       context.Context
	client    *storage.Client
	dryRun    = false
	mediaInfo *store.MediaInfo
)

func Register(app *kingpin.Application, actual *bool) {
	app.Command("new", "Process new photos").Action(
		func(_ *kingpin.ParseContext) error {
			dryRun = !*actual
			handler()
			return nil
		})
}

func handler() {
	ctx = context.Background()
	err := fmt.Errorf("") // Next line can't use :=
	client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mi, err := store.New(ctx, client, config.MediaInfoBucket())
	mediaInfo = mi
	if err != nil {
		log.Fatalf("NewMediaInfo failed: %v", err)
	}

	c := 0
	log.Printf("Looking for photos in: %s", config.PhotoQueueBucket())
	//log.Printf("Processing files in bucket: %s", bktiter.Name)
	var set []*storage.ObjectAttrs
	bkt := client.Bucket(config.PhotoQueueBucket())
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

		// This should probably disappear
		c++
		if c > 500 {
			log.Fatalf("Foo")
		}
	}
	if len(set) > 0 {
		processPhotoSet(set)
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
		log.Printf("Skipping invalid photo set: %v", err)
		return
	}

	dryMsg := ""
	if dryRun {
		dryMsg = "DRY_RUN: "
	}

	// cp files
	for _, a := range attr {
		name := lookupSha256(a, media)
		_, err := client.Bucket(config.PhotoStorageBucket()).Object(name).Attrs(ctx)
		if err == nil {
			log.Printf("File already exists: %v", name)
			continue
		}
		if err != nil && err != storage.ErrObjectNotExist {
			log.Fatalf("Error checking for file: %v, %v", name, err)
		}
		log.Printf("%sCopying: %v/%v -> %v/%v\n",
			dryMsg, a.Bucket, a.Name, config.PhotoStorageBucket(), name)
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
	// TODO: Can't tell if this fails... DANGER
	if !dryRun {
		mediaInfo.Insert(ctx, client, media)
	}

	// delete files
	for _, a := range attr {
		log.Printf("%sDeleting: %v/%v\n", dryMsg, a.Bucket, a.Name)
		if !dryRun {
			if err := client.Bucket(a.Bucket).Object(a.Name).Delete(ctx); err != nil {
				log.Fatalf("Failed to delete: %v", err)
			}
		}
	}
}

func lookupSha256(attr *storage.ObjectAttrs, m *message.Media) string {
	for _, mf := range m.File {
		if util.MD5(mf.Md5Sum) == util.MD5(attr.MD5) {
			return hex.EncodeToString(mf.Sha256Sum)
		}
	}
	log.Fatalf("Unable to find file")
	return ""
}

// attr is expected to contain exactly one JPG and zero or more other related files (RAWs for
// example.) A message.Media object will contain the EXIF information from this JPG and basic file
// information about the JPG and any other files included.
func convertToMedia(attr []*storage.ObjectAttrs) (*message.Media, error) {
	var jpg []*storage.ObjectAttrs
	var other []*storage.ObjectAttrs
	for _, a := range attr {
		if strings.ToLower(suffix(a.Name)) == "jpg" ||
			strings.ToLower(suffix(a.Name)) == "jpeg" {
			jpg = append(jpg, a)
		} else {
			other = append(other, a)
		}
	}

	if len(jpg) != 1 {
		return nil, fmt.Errorf("Invalid set of JPGs found %s: %v",
			attr[0].Name, jpg)
	}

	var media message.Media
	jpginfo, err := getFile(jpg[0], "tmpfile.jpg")
	if err != nil {
		log.Fatalf("Failed to retrieve file: %v", err)
	}
	media.File = append(media.File, jpginfo)

	mediainfo, err := exif.Parse("tmpfile.jpg")
	if err != nil {
		log.Fatalf("Failed to get EXIF data from JPG: %v", err)
	}
	media.Photo = mediainfo

	for _, a := range other {
		info, err := getFile(a, "")
		if err != nil {
			log.Fatalf("Failed to retrieve file: %v", err)
		}
		media.File = append(media.File, info)
	}

	media.Key = media.File[0].Sha256Sum
	media.TimestampSeconds = time.Now().Unix()
	media.Tags = []string{"new"}

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
