package config

//go:generate protoc --go_out=. config.proto

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/protobuf/proto"
)

const (
	StorageBucket = ""
	ProjectID     = ""
	PhotoPath     = ""
	Path          = ""
)

var (
	instance     CuratorConfig
	once         sync.Once
	cameraModels map[string]string
)

func CameraModelAbbreviation(name string) string {
	parse()

	m, ok := cameraModels[name]
	if !ok {
		return "UNKNOWN"
	}
	return m
}

func PhotoQueueBucket() string {
	parse()

	return instance.PhotoQueueBucket
}

func PhotoStorageBucket() string {
	parse()

	return instance.PhotoStorageBucket
}

func ValidLabels() []string {
	parse()

	return instance.ValidLabels
}

func MediaInfoBucket() string {
	parse()

	return instance.PhotoInfoBucket
}

func parse() {
	once.Do(func() {
		path := filepath.Join(os.Getenv("HOME"), ".curator.pb.ascii")
		log.Printf("Loading config from: %s", path)
		config, err := ioutil.ReadFile(filepath.Join(path))
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		err = proto.UnmarshalText(string(config), &instance)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}

		cameraModels = make(map[string]string)
		for _, o := range instance.CameraModels {
			cameraModels[o.ExifModel] = o.Abbreviation
		}
	})
}
