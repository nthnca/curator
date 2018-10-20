package util

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/pkg/mediainfo"
)

// Base returns the filename without the suffix.
func Base(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[:i]
	}
	return a
}

// Suffix returns the suffix of the filename.
func Suffix(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[i+1:]
	}
	return ""
}

// GetFile retrieves a file and returns the MD5, SHA256, file name, and file size. If a file
// was passed in it stores the file there.
func GetFile(ctx context.Context, client *storage.Client, attrs *storage.ObjectAttrs, file *os.File) (*mediainfo.FileInfo, error) {
	rc, err := client.Bucket(attrs.Bucket).Object(attrs.Name).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create reader: %v", err)
	}

	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	if int64(len(slurp)) != attrs.Size {
		return nil, fmt.Errorf("File size didn't match.")
	}

	md := md5.Sum([]byte(slurp))
	if !bytes.Equal(attrs.MD5, md[:]) {
		return nil, fmt.Errorf("MD5 sum didn't match, file corrupted?")
	}

	sha := sha256.Sum256([]byte(slurp))
	sub := strings.Split(attrs.Name, "/")
	name := sub[len(sub)-1]

	if file != nil {
		if _, err := file.Write(slurp); err != nil {
			return nil, fmt.Errorf("Failed to write file: %v", err)
		}

		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("Failed to close file: %v", err)
		}
	}

	return &mediainfo.FileInfo{
		Filename:    name,
		Md5Sum:      md[:],
		Sha256Sum:   sha[:],
		SizeInBytes: attrs.Size,
	}, nil
}
