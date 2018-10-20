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
	"github.com/nthnca/curator/pkg/mediainfo/message"
)

func Base(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[:i]
	}
	return a
}

func Suffix(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[i+1:]
	}
	return ""
}

func GetFile(client *storage.Client, ctx context.Context, attrs *storage.ObjectAttrs, file *os.File) (*message.FileInfo, error) {
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

	return &message.FileInfo{
		Filename:    name,
		Md5Sum:      md[:],
		Sha256Sum:   sha[:],
		SizeInBytes: attrs.Size,
	}, nil
}
