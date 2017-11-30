package config

import (
	"context"
	"io/ioutil"
	"log"
	"net/url"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/data/message"
)

var (
	instance     message.Config
	once         sync.Once
	CameraModels map[string]string
)

func Get() *message.Config {
	once.Do(func() {
		u, err := url.Parse(curatorConfig)
		if err != nil {
			log.Fatalf("Unable to parse path: %v", err)
		}

		// TODO: Should handle file:// and local paths.
		if u.Scheme != "gs" {
			log.Fatalf("Only gs:// paths are supported: %v", curatorConfig)
		}

		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create storage client: %v", err)
		}

		rc, err := client.Bucket(u.Host).Object(u.Path[1:]).NewReader(ctx)
		if err != nil {
			log.Fatalf("Failed to create reader: %v", err)
		}

		slurp, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		err = proto.UnmarshalText(string(slurp), &instance)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}

		CameraModels = make(map[string]string)
		for _, o := range instance.CameraModels {
			CameraModels[o.ExifModel] = o.Abbreviation
		}
	})
	return &instance
}
