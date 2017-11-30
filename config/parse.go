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

func PhotoInfoBucket() string {
	parse()

	return instance.PhotoInfoBucket
}

func parse() {
	once.Do(func() {
		config, err := ioutil.ReadFile(filepath.Join(
			os.Getenv("HOME"), ".curator.pb.ascii"))
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
