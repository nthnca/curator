package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/cmd/handler/analyze"
	"github.com/nthnca/curator/cmd/handler/cache"
	sy "github.com/nthnca/curator/cmd/handler/sync"
	"github.com/nthnca/curator/cmd/handler/update"
	"github.com/nthnca/curator/config"

	"github.com/nthnca/gobuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	gobuild.RegisterCommands(app, config.Path, config.ProjectID)
	app.Command("sync", "Update cloud data as needed").Action(
		func(_ *kingpin.ParseContext) error {
			sy.Handler()
			return nil
		})
	app.Command("cache", "Update datastore caches").Action(
		func(_ *kingpin.ParseContext) error {
			cache.Handler()
			return nil
		})
	app.Command("analyze", "list entries").Action(
		func(_ *kingpin.ParseContext) error {
			analyze.Handler()
			return nil
		})
	app.Command("update", "list entries").Action(
		func(_ *kingpin.ParseContext) error {
			update.Handler()
			return nil
		})
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
