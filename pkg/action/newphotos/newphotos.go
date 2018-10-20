package newphotos

import (
	"context"
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
	"github.com/nthnca/curator/pkg/mediainfo"
	"github.com/nthnca/curator/pkg/util"
	objectstore "github.com/nthnca/object-store"
	"google.golang.org/api/iterator"
)

// Options allows you to modify the behavior of the NewPhotos action.
type Options struct {
	// Ctx is a valid context.Context to run this command under.
	Ctx context.Context

	// Storage is a Google Cloud Storage client.
	Storage *storage.Client

	// ObjStore is an ObjectStore client
	ObjStore *objectstore.ObjectStore

	// Cfg is the configuration settings.
	Cfg *config.Config

	// DryRun. If true don't actually make changes, print the changes you would have made.
	DryRun bool
}

type action struct {
	ctx        context.Context
	client     *storage.Client
	store      *objectstore.ObjectStore
	cfg        *config.Config
	numThreads int
	dryRun     bool
}

type file struct {
	attrs *storage.ObjectAttrs
	info  *mediainfo.FileInfo
}

// Do performs the new photos action, we should improve this documentation.
func Do(opts *Options) error {
	var act action
	act.numThreads = 1
	act.dryRun = opts.DryRun
	act.ctx = opts.Ctx
	act.client = opts.Storage
	act.store = opts.ObjStore
	act.cfg = opts.Cfg

	var fatalError error
	var wg sync.WaitGroup
	chSet := make(chan []*file)

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
	log.Printf("Looking for photos in: %s", act.cfg.PhotoQueueBucket())
	set := []*file{}
	bkt := act.client.Bucket(act.cfg.PhotoQueueBucket())
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

		if len(set) > 0 && util.Base(set[0].attrs.Name) != util.Base(obj.Name) {
			chSet <- set
			// Reset "set", don't want to modify after passing it along.
			set = []*file{}
		}

		f := &file{attrs: obj}

		// set has jpgs first, other files after.
		if isJpg(f.attrs.Name) {
			set = append([]*file{f}, set...)
		} else {
			set = append(set, f)
		}
	}
	close(chSet)
	wg.Wait()

	return fatalError
}

func isJpg(name string) bool {
	return strings.ToLower(util.Suffix(name)) == "jpg" ||
		strings.ToLower(util.Suffix(name)) == "jpeg"
}

func (act *action) dryRunMsg() string {
	if act.dryRun {
		return "[DRY RUN] "
	}
	return ""
}

func (act *action) processPhotoSet(files []*file) error {
	if !isJpg(files[0].attrs.Name) || (len(files) > 1 && isJpg(files[1].attrs.Name)) {
		log.Printf("Need exactly 1 jpg file: %s", files[0].attrs.Name)
		// This isn't an error, but we can't continue processing this set.
		return nil
	}

	media, err := act.createMediaProto(files)
	if err != nil {
		return fmt.Errorf("Failed to process files: %v", err)
	}

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
		log.Printf("Failed to delete files: %v", err)
	}
	return nil
}

func (act *action) copyFiles(files []*file, metadata *mediainfo.Media) error {
	for _, a := range files {
		// name := a.hex
		name := hex.EncodeToString(a.info.Sha256Sum)

		log.Printf("%sCopying: %v/%v -> %v/%v\n",
			act.dryRunMsg(), a.attrs.Bucket, a.attrs.Name, act.cfg.PhotoStorageBucket(), name)

		_, err := act.client.Bucket(act.cfg.PhotoStorageBucket()).Object(name).Attrs(act.ctx)
		if err == nil {
			log.Printf("No need to copy, file already exists: %v", name)
			continue
		}

		if err != storage.ErrObjectNotExist {
			return fmt.Errorf("Checking for file: %v, %v", name, err)
		}

		if !act.dryRun {
			src := act.client.Bucket(a.attrs.Bucket).Object(a.attrs.Name)
			dest := act.client.Bucket(act.cfg.PhotoStorageBucket()).Object(name)
			_, err = dest.CopierFrom(src).Run(act.ctx)
			if err != nil {
				return fmt.Errorf("Copying file: %v", err)
			}
		}
	}
	return nil
}

func (act *action) insertMetadata(metadata *mediainfo.Media) error {
	if !act.dryRun {
		data, err := proto.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("Failed to marshal proto: %v", err)
		}

		err = act.store.Insert(act.ctx, hex.EncodeToString(metadata.Key), data)
		if err != nil {
			return fmt.Errorf("Attempting insert: %v", err)
		}
	}
	return nil
}

func (act *action) deleteFiles(files []*file) error {
	for _, a := range files {
		log.Printf("%sDeleting: %v/%v\n", act.dryRunMsg(), a.attrs.Bucket, a.attrs.Name)
		if !act.dryRun {
			if err := act.client.Bucket(a.attrs.Bucket).Object(a.attrs.Name).Delete(act.ctx); err != nil {
				return fmt.Errorf("Failed to delete: %v", err)
			}
		}
	}
	return nil
}

// attr is expected to contain exactly one JPG and zero or more other related files (RAWs for
// example.) A mediainfo.Media object will contain the EXIF information from act JPG and basic file
// information about the JPG and any other files included.
func (act *action) createMediaProto(files []*file) (*mediainfo.Media, error) {
	tmpfile, err := ioutil.TempFile("", "jpgpic")
	if err != nil {
		return nil, fmt.Errorf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	fileinfo, err := util.GetFile(act.ctx, act.client, files[0].attrs, tmpfile)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve file: %v", err)
	}

	var media mediainfo.Media
	media.File = append(media.File, fileinfo)
	files[0].info = fileinfo

	media.Photo, err = mediainfo.ParseExif(tmpfile.Name())
	if err != nil {
		return nil, fmt.Errorf("Failed to get EXIF data from JPG: %v", err)
	}

	for i := 1; i < len(files); i++ {
		info, err := util.GetFile(act.ctx, act.client, files[1].attrs, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve file: %v", err)
		}
		media.File = append(media.File, info)
		files[i].info = info
	}

	media.Key = media.File[0].Sha256Sum
	media.TimestampSeconds = time.Now().Unix()
	media.Tags = []string{"new"}

	return &media, nil
}
