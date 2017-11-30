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

	"github.com/nthnca/curator/config"
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

func getCanonicalName(datetime, model, key string) string {
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

func IdentifyPhoto(path string, md5, sha256 []byte) (*message.Photo, error) {
	base := filepath.Base(path)
	key := strings.TrimSuffix(base, ".jpg")
	key = strings.TrimSuffix(key, ".JPG")
	sub := strings.Split(key, "-")
	key = sub[len(sub)-1]
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

	photo.Key = key
	photo.NumBytes = fi.Size()
	photo.Md5Sum = md5
	photo.Sha256Sum = sha256

	output := string(buffer[:])
	photo.Properties.EpochInSeconds = getTime(output, reDateTime)
	photo.Properties.Width = int32(getInt(output, reWidth))
	photo.Properties.Height = int32(getInt(output, reHeight))
	photo.Properties.Make = getString(output, reMake)
	photo.Properties.Model = getString(output, reModel)
	photo.Properties.Aperture = getFraction(output, reAperture)
	photo.Properties.ExposureTime = getFraction(output, reExposureTime)
	photo.Properties.FocalLength = getFraction(output, reFocalLength)
	photo.Properties.Iso = int32(getInt(output, reISOSpeedRating))

	photo.Path = getCanonicalName(getString(output, reDateTime),
		getString(output, reModel), key)
	if photo.Path == "" {
		return nil, fmt.Errorf("No date for photo")
	}

	photo.Properties.Model = getString(output, reModel)
	return &photo, nil
}

func getInt(buffer string, regex *regexp.Regexp) int {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return 0
	}

	v, err := strconv.Atoi(strings.Split(m[1], ",")[0])
	if err != nil {
		log.Fatalf("%v", err)
	}

	return v
}

func getString(buffer string, regex *regexp.Regexp) string {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return ""
	}

	return m[1]
}

func getTime(buffer string, regex *regexp.Regexp) int64 {
	m := regex.FindStringSubmatch(buffer)
	if len(m) != 2 {
		return 0
	}

	a := m[1]
	a = strings.Replace(a, " ", "T", -1)
	a = strings.Replace(a, ":", "-", 2)
	a += "Z"
	c, err := time.Parse(time.RFC3339, a)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return c.Unix()
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
		Numerator:   int32(a),
		Denominator: int32(b)}
}
