package util

import (
	"fmt"
	"strings"

	"github.com/nthnca/curator/pkg/config"
	"github.com/nthnca/curator/pkg/mediainfo"
)

// GetCanonicalName creates a standardize photo file name that looks like
// <date>-<time>-<camera-model>-<filename>.jpg
func GetCanonicalName(config *config.Config, media *mediainfo.Media) string {
	if media.Photo.Datetime == "" {
		return media.Name
	}
	datetime := media.Photo.Datetime
	model := media.Photo.Model
	key := media.File[0].Filename
	key = strings.TrimSuffix(key, ".jpg")
	key = strings.TrimSuffix(key, ".JPG")
	split := strings.Split(strings.Replace(datetime, ":", "", -1), " ")
	date := "00000000"
	time := "000000"
	if len(split) == 2 {
		date = split[0]
		time = split[1]
	}

	m := config.CameraModelAbbreviation(model)

	return fmt.Sprintf("%s-%s-%s-%s.jpg", date, time, m, key)
}
