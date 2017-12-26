package exif

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Parse populates a PhotoInfo protobuf with the exif data from the photo at 'path'.
func Parse(path string) (*message.PhotoInfo, error) {
	cmd := exec.Command("identify", "-format", "%[exif:*]", path)
	buffer, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ImageMagick failed: %v", err)
	}

	output := string(buffer[:])
	pi := &message.PhotoInfo{
		EpochInSeconds: getTime(output, reDateTime),
		Datetime:       getString(output, reDateTime),
		Make:           getString(output, reMake),
		Model:          getString(output, reModel),
		Aperture:       getFraction(output, reAperture),
		ExposureTime:   getFraction(output, reExposureTime),
		FocalLength:    getFraction(output, reFocalLength),
		Iso:            int32(getInt(output, reISOSpeedRating)),
		Width:          int32(getInt(output, reWidth)),
		Height:         int32(getInt(output, reHeight)),
	}

	return pi, nil
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
	a := getString(buffer, regex)
	if a == "" {
		return 0
	}

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
