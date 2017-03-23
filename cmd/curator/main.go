package main

import (
	"github.com/nthnca/curator/config"
	"github.com/nthnca/easybuild"
)

func main() {
	easybuild.Build(config.Path, config.ProjectID)
}
