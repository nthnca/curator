package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config/internal"
)

// Config contains the users configuration settings for Curator.
type Config struct {
	instance     internal.CuratorConfig
	cameraModels map[string]string
}

// CameraModelAbbreviation returns the users preferred short string for a given camera name.
func (config *Config) CameraModelAbbreviation(name string) string {
	m, ok := config.cameraModels[name]
	if !ok {
		return "UNKNOWN"
	}
	return m
}

// PhotoQueueBucket is the GCS bucket where new photos are stored.
func (config *Config) PhotoQueueBucket() string {
	return config.instance.PhotoQueueBucket
}

// PhotoQueueBucket is the GCS bucket where photos are stored after they are processed.
func (config *Config) PhotoStorageBucket() string {
	return config.instance.PhotoStorageBucket
}

// ValidLabels are the list of labels that a user can use to tag their photos.
func (config *Config) ValidLabels() []string {
	return config.instance.ValidLabels
}

// MetadataBucket is the GCS bucket where photo metadata is stored.
func (config *Config) MetadataBucket() string {
	return config.instance.PhotoMetadataBucket
}

// MetadataBucket is the filename prefix for the photo metadata.
func (config *Config) MetadataPath() string {
	return config.instance.PhotoMetadataPath
}

// New parses the users configuration and returns a Config.
func New() *Config {
	configFile := os.Getenv("CONFIG_FILE")

	log.Printf("Loading config from: %s", configFile)
	data, err := ioutil.ReadFile(filepath.Join(configFile))
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	cfg := Config{}
	err = proto.UnmarshalText(string(data), &cfg.instance)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	cfg.cameraModels = make(map[string]string)
	for _, o := range cfg.instance.CameraModels {
		cfg.cameraModels[o.ExifModel] = o.Abbreviation
	}

	return &cfg
}
