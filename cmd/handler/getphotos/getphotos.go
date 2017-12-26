package getphotos

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/data/mediainfo"
	"github.com/nthnca/curator/data/message"
	"github.com/nthnca/curator/util"
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

	photoData.Load(ctx, client)

	mi, err := mediainfo.NewMediaInfo(ctx, client, "c1410a-photo-storage-info")
	if err != nil {
		log.Fatalf("NewMediaInfo failed: %v", err)
	}

	for _, e := range photoData.Data {
		//if i >= 25 {
		//	break
		//}
		log.Printf("%v", e)
		ne := &message.Media{
			Key:  e.Sha256Sum,
			Name: e.Path,
			TimestampSecondsSinceEpoch: 0,
		}
		ne.Photo = &message.PhotoInfo{
			EpochInSeconds: e.Properties.EpochInSeconds,
			Make:           e.Properties.Make,
			Model:          e.Properties.Model,
			Aperture:       e.Properties.Aperture,
			ExposureTime:   e.Properties.ExposureTime,
			FocalLength:    e.Properties.FocalLength,
			Iso:            e.Properties.Iso,
			Width:          e.Properties.Width,
			Height:         e.Properties.Height,
		}
		ne.File = append(ne.File, &message.FileInfo{
			Filename:    e.Key + ".jpg",
			Md5Sum:      e.Md5Sum,
			Sha256Sum:   e.Sha256Sum,
			SizeInBytes: e.NumBytes,
			// Type:,
		})
		mi.Insert(ctx, client, ne)
		log.Printf("### %v", ne)
	}
	mi.Flush(ctx, client)
}

/*
	fmt.Printf("mkdir .pics\n")
	size := len(photoData.Data)
	for i, _ := range photoData.Data {
		if i >= 25 {
			break
		}

		iter := photoData.Data[size-i-1]
		fmt.Printf("gsutil cp gs://%s/%s .pics/%s\n",
			config.PhotoStorageBucket(), iter.GetPath(), hex.EncodeToString(iter.GetSha256Sum()))
		fmt.Printf("ln .pics/%s %s\n", hex.EncodeToString(iter.GetSha256Sum()),
			iter.GetPath())
	}
*/
