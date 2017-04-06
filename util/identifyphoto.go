package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nthnca/curator/data/message"
)

var (
	reDateTime       = regexp.MustCompile(`DateTimeOriginal=(.*)`)
	reHeight         = regexp.MustCompile(`ExifImageLength=(.*)`)
	reWidth          = regexp.MustCompile(`ExifImageWidth=(.*)`)
	reMake           = regexp.MustCompile(`Make=(.*)`)
	reModel          = regexp.MustCompile(`Model=(.*)`)
	reAperture       = regexp.MustCompile(`FNumber=(.*)`)
	reExposureTime   = regexp.MustCompile(`ExposureTime=(.*)`)
	reFocalLength    = regexp.MustCompile(`FocalLength=(.*)`)
	reISOSpeedRating = regexp.MustCompile(`ISOSpeedRatings=(.*)`)
)

func getInt32(buffer string, regex *regexp.Regexp) *int32 {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return nil
	}

	v, err := strconv.Atoi(strings.Split(m[1], ",")[0])
	if err != nil {
		log.Fatalf("%v", err)
	}

	return proto.Int(v)
}

func getString(buffer string, regex *regexp.Regexp) *string {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return nil
	}

	return proto.String(m[1])
}

func getTime(buffer string, regex *regexp.Regexp) *int64 {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return nil
	}

	a := m[1]
	a = strings.Replace(a, " ", "T", -1)
	a = strings.Replace(a, ":", "-", 2)
	a += "Z"
	c, err := time.Parse(time.RFC3339, a)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return proto.Int64(c.Unix())
}

func getFraction(buffer string, regex *regexp.Regexp) *message.Fraction {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return nil
	}

	x := strings.Split(m[1], "/")
	a, err := strconv.Atoi(x[0])
	if err != nil {
		log.Fatalf("%v", err)
	}

	b := 1
	if len(x) == 2 {
		b, err = strconv.Atoi(x[1])
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	return &message.Fraction{
		Numerator:   proto.Int(a),
		Denominator: proto.Int(b)}
}

func IdentifyPhoto(path string) (*message.Photo, error) {
	base := filepath.Base(path)
	key := strings.TrimSuffix(base, ".jpg")
	if key == base {
		return nil, fmt.Errorf("Invalid photo name: %v", path)
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to stat file: %v", err)
	}

	cmd := exec.Command("identify", "-format", "%[exif:*]", path)
	buffer, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ImageMagick failed: %v", err)
	}

	photo := message.Photo{
		Properties: &message.Photo_PhotoProperties{}}

	photo.Key = proto.String(key)
	photo.Path = proto.String(path)
	photo.Bytes = proto.Int64(fi.Size())

	output := string(buffer[:])
	photo.Properties.OriginalEpoch = getTime(output, reDateTime)
	photo.Properties.Width = getInt32(output, reWidth)
	photo.Properties.Height = getInt32(output, reHeight)
	photo.Properties.Make = getString(output, reMake)
	photo.Properties.Model = getString(output, reModel)
	photo.Properties.Iso = getInt32(output, reISOSpeedRating)
	photo.Properties.Aperture = getFraction(output, reAperture)
	photo.Properties.ExposureTime = getFraction(output, reExposureTime)
	photo.Properties.FocalLength = getFraction(output, reFocalLength)
	return &photo, nil
}
