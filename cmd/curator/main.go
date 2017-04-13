package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/handler/analyze"
	sy "github.com/nthnca/curator/handler/sync"
	"github.com/nthnca/curator/handler/update"

	"github.com/nthnca/gobuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	gobuild.RegisterCommands(app, config.Path, config.ProjectID)
	app.Command("sync", "Update cloud data as needed").Action(sy.Handler)
	app.Command("analyze", "list entries").Action(analyze.Handler)
	app.Command("update", "list entries").Action(update.Handler)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
