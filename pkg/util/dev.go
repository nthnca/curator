package util

import (
	"os"
)

func IsDevAppServer() bool {
	return os.Getenv("RUN_WITH_DEVAPPSERVER") == "1"
}
