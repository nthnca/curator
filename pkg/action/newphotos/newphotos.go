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
	"os"
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
)

// Options is the set of valid options for the NewPhotos action.
type Options struct {
	DryRun bool
}

type action struct {
	ctx        context.Context
	client     *storage.Client
	store      *objectstore.ObjectStore
	numThreads int
	dryRun     bool
}

// Do performs the new photos action, we should improve this documentation.
func Do(opts *Options) error {
	var act action
	act.numThreads = 1
	act.dryRun = opts.DryRun
	act.ctx = context.Background()

	log.Printf("Dry Run: %v", act.dryRun)

	var err error
	act.client, err = storage.NewClient(act.ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	act.store, err = objectstore.New(act.ctx, act.client, config.MetadataBucket(), config.MetadataPath())
	if err != nil {
		return fmt.Errorf("New ObjectStore failed: %v", err)
	}

	return act.processPhotos()
}

func (act *action) processPhotos() error {
	var fatalError error
	var wg sync.WaitGroup
	chSet := make(chan []*storage.ObjectAttrs)

	for i := 0; i < act.numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for set := range chSet {
				// Someone hit a fatal error, stop processing.
				if fatalError != nil {
					break
				}

				err := act.processPhotoSet(set)
				if err != nil {
					fatalError = fmt.Errorf("Failed to process photo set:%v", err)
				}
			}
		}()
	}

	// Ordering is guaranteed: https://cloud.google.com/storage/docs/listing-objects
	log.Printf("Looking for photos in: %s", config.PhotoQueueBucket())
	set := []*storage.ObjectAttrs{}
	bkt := act.client.Bucket(config.PhotoQueueBucket())
	for it := bkt.Objects(act.ctx, nil); ; {
		obj, err := it.Next()
		if err == iterator.Done {
			if len(set) > 0 {
				chSet <- set
			}
			break
		}
		if err != nil {
			fatalError = fmt.Errorf("Failed to iterate through objects: %v", err)
		}
		if fatalError != nil {
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

	return fatalError
}

func (act *action) processPhotoSet(files []*storage.ObjectAttrs) error {
	media, err := act.createMediaProto(files)
	if err != nil {
		// log.Printf("Failed to process files: %v", err)
		// return nil
		return fmt.Errorf("Failed to process files: %v", err)
	}

	// log.Fatalf("PROTO OUTPUT:\n%v", proto.MarshalTextString(media))

	err = act.copyFiles(files, media)
	if err != nil {
		return fmt.Errorf("Failed to copy photos: %v", err)
	}

	err = act.insertMetadata(media)
	if err != nil {
		return fmt.Errorf("Failed to insert photo metadata: %v", err)
	}

	err = act.deleteFiles(files)
	if err != nil {
		return fmt.Errorf("Failed to delete files: %v", err)
	}
	return nil
}

func (act *action) copyFiles(files []*storage.ObjectAttrs, metadata *message.Media) error {
	dryMsg := ""
	if act.dryRun {
		dryMsg = "DRY_RUN: "
	}

	for _, a := range files {
		name := lookupSha256(a, metadata)
		_, err := act.client.Bucket(config.PhotoStorageBucket()).Object(name).Attrs(act.ctx)
		if err == nil {
			log.Printf("File already exists: %v", name)
			continue
		}
		if err != nil && err != storage.ErrObjectNotExist {
			log.Fatalf("Error checking for file: %v, %v", name, err)
		}
		log.Printf("%sCopying: %v/%v -> %v/%v\n",
			dryMsg, a.Bucket, a.Name, config.PhotoStorageBucket(), name)
		if !act.dryRun {
			src := act.client.Bucket(a.Bucket).Object(a.Name)
			dest := act.client.Bucket(config.PhotoStorageBucket()).Object(name)
			_, err = dest.CopierFrom(src).Run(act.ctx)
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}
	return nil
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

func (act *action) insertMetadata(metadata *message.Media) error {
	if !act.dryRun {
		data, err := proto.Marshal(metadata)
		if err != nil {
			log.Fatalf("Failed to marshal proto: %v", err)
		}

		err = act.store.Insert(act.ctx, hex.EncodeToString(metadata.Key), data)
		if err != nil {
			log.Fatalf("action isn't good!: %v", err)
		}
	}
	return nil
}

func (act *action) deleteFiles(files []*storage.ObjectAttrs) error {
	for _, a := range files {
		log.Printf("%sDeleting: %v/%v\n", "dryMsg", a.Bucket, a.Name)
		if !act.dryRun {
			if err := act.client.Bucket(a.Bucket).Object(a.Name).Delete(act.ctx); err != nil {
				log.Fatalf("Failed to delete: %v", err)
			}
		}
	}
	return nil
}

// attr is expected to contain exactly one JPG and zero or more other related files (RAWs for
// example.) A message.Media object will contain the EXIF information from act JPG and basic file
// information about the JPG and any other files included.
func (act *action) createMediaProto(attr []*storage.ObjectAttrs) (*message.Media, error) {
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

	tmpfile, err := ioutil.TempFile("", "jpgpic")
	if err != nil {
		return nil, fmt.Errorf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	jpginfo, err := act.getFile(jpg[0], tmpfile)
	if err != nil {
		log.Fatalf("Failed to retrieve file: %v", err)
	}

	var media message.Media
	media.File = append(media.File, jpginfo)

	mediainfo, err := exif.Parse(tmpfile.Name())
	if err != nil {
		log.Fatalf("Failed to get EXIF data from JPG: %v", err)
	}
	media.Photo = mediainfo

	for _, a := range other {
		info, err := act.getFile(a, nil)
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

func (act *action) getFile(attr *storage.ObjectAttrs, file *os.File) (*message.FileInfo, error) {
	rc, err := act.client.Bucket(attr.Bucket).Object(attr.Name).NewReader(act.ctx)
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

	if file != nil {
		if _, err := file.Write(slurp); err != nil {
			return nil, fmt.Errorf("Failed to write file: %v", err)
		}

		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("Failed to close file: %v", err)
		}
	}

	return &message.FileInfo{
		Filename:    name,
		Md5Sum:      md[:],
		Sha256Sum:   sha[:],
		SizeInBytes: attr.Size,
	}, nil
}
