package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/pkg/config/internal"
)

type Config struct {
	instance     internal.CuratorConfig
	cameraModels map[string]string
}

func (config *Config) CameraModelAbbreviation(name string) string {
	m, ok := config.cameraModels[name]
	if !ok {
		return "UNKNOWN"
	}
	return m
}

func (config *Config) PhotoQueueBucket() string {
	return config.instance.PhotoQueueBucket
}

func (config *Config) PhotoStorageBucket() string {
	return config.instance.PhotoStorageBucket
}

func (config *Config) ValidLabels() []string {
	return config.instance.ValidLabels
}

func (config *Config) MediaInfoBucket() string {
	return config.instance.PhotoInfoBucket
}

func (config *Config) MetadataBucket() string {
	return config.instance.PhotoMetadataBucket
}

func (config *Config) MetadataPath() string {
	return config.instance.PhotoMetadataPath
}

func New() *Config {
	config_file := os.Getenv("CONFIG_FILE")

	log.Printf("Loading config from: %s", config_file)
	data, err := ioutil.ReadFile(filepath.Join(config_file))
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
