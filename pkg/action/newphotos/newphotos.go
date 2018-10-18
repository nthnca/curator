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
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo/exif"
	"github.com/nthnca/curator/pkg/mediainfo/message"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
	"google.golang.org/api/iterator"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const kMaxThreads = 5

func Register(app *kingpin.Application, actual *bool) {
	app.Command("new", "Process new photos").Action(
		func(_ *kingpin.ParseContext) error {
			this := This{dryRun: !*actual}
			err := this.handler()
			if err != nil {
				log.Fatalf("%v", err)
			}
			return nil
		})
}

type Options struct {
	DryRun bool
}

type This struct {
	ctx    context.Context
	client *storage.Client
	store  *objectstore.ObjectStore
	dryRun bool
	err    error
}

func Do(opts *Options) {
	this := This{dryRun: opts.DryRun}
	err := this.handler()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (this *This) handler() error {
	var err error
	this.ctx = context.Background()
	this.client, err = storage.NewClient(this.ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	this.store, err = objectstore.New(this.ctx, this.client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		return fmt.Errorf("New ObjectStore failed: %v", err)
	}

	chSet := make(chan []*storage.ObjectAttrs)

	var wg sync.WaitGroup
	for i := 0; i < kMaxThreads; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for set := range chSet {
				log.Printf("P %d", x)
				this.processPhotoSet(set, fmt.Sprintf("%d", x))
			}
		}(i)
	}

	// Ordering is guaranteed: https://cloud.google.com/storage/docs/listing-objects
	log.Printf("Looking for photos in: %s", config.PhotoQueueBucket())
	set := []*storage.ObjectAttrs{}
	bkt := this.client.Bucket(config.PhotoQueueBucket())
	for it := bkt.Objects(this.ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			if len(set) > 0 {
				chSet <- set
			}
			break
		}
		if err != nil {
			this.err = fmt.Errorf("Failed to iterate through objects: %v", err)
		}
		if this.err != nil {
			break
		}

		if len(set) > 0 && util.Base(set[0].Name) != util.Base(obj.Name) {
			chSet <- set
			// Reset "set", don't want to modify after passing it along.
			set = []*storage.ObjectAttrs{}
		}
		set = append(set, obj)
	}
	close(chSet)
	wg.Wait()

	return this.err
}

func (this *This) processPhotoSet(attr []*storage.ObjectAttrs, prefix string) {
	media, err := this.createMediaProto(attr, prefix)
	if err != nil {
		log.Printf("Skipping invalid photo set: %v", err)
		return
	}

	dryMsg := ""
	if this.dryRun {
		dryMsg = "DRY_RUN: "
	}

	// cp files
	for _, a := range attr {
		name := lookupSha256(a, media)
		_, err := this.client.Bucket(config.PhotoStorageBucket()).Object(name).Attrs(this.ctx)
		if err == nil {
			log.Printf("File already exists: %v", name)
			continue
		}
		if err != nil && err != storage.ErrObjectNotExist {
			log.Fatalf("Error checking for file: %v, %v", name, err)
		}
		log.Printf("%sCopying: %v/%v -> %v/%v\n",
			dryMsg, a.Bucket, a.Name, config.PhotoStorageBucket(), name)
		if !this.dryRun {
			src := this.client.Bucket(a.Bucket).Object(a.Name)
			dest := this.client.Bucket(config.PhotoStorageBucket()).Object(name)
			_, err = dest.CopierFrom(src).Run(this.ctx)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}

	// save metadata
	if !this.dryRun {
		data, err := proto.Marshal(media)
		if err != nil {
			log.Fatalf("Failed to marshal proto: %v", err)
		}

		err = this.store.Insert(this.ctx, hex.EncodeToString(media.Key), data)
		if err != nil {
			log.Fatalf("This isn't good!: %v", err)
		}
	}

	// delete files
	for _, a := range attr {
		log.Printf("%sDeleting: %v/%v\n", dryMsg, a.Bucket, a.Name)
		if !this.dryRun {
			if err := this.client.Bucket(a.Bucket).Object(a.Name).Delete(this.ctx); err != nil {
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
func (this *This) createMediaProto(attr []*storage.ObjectAttrs, prefix string) (*message.Media, error) {
	var jpg []*storage.ObjectAttrs
	var other []*storage.ObjectAttrs
	for _, a := range attr {
		if strings.ToLower(util.Suffix(a.Name)) == "jpg" ||
			strings.ToLower(util.Suffix(a.Name)) == "jpeg" {
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
	jpginfo, err := this.getFile(jpg[0], prefix+"tmpfile.jpg")
	if err != nil {
		log.Fatalf("Failed to retrieve file: %v", err)
	}
	media.File = append(media.File, jpginfo)

	mediainfo, err := exif.Parse(prefix + "tmpfile.jpg")
	if err != nil {
		log.Fatalf("Failed to get EXIF data from JPG: %v", err)
	}
	media.Photo = mediainfo

	for _, a := range other {
		info, err := this.getFile(a, "")
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

func (this *This) getFile(attr *storage.ObjectAttrs, localPath string) (*message.FileInfo, error) {
	rc, err := this.client.Bucket(attr.Bucket).Object(attr.Name).NewReader(this.ctx)
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
