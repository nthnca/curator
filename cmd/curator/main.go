package main

import (
	"os"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/easybuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	easybuild.RegisterCommands(app, config.Path, config.ProjectID)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
