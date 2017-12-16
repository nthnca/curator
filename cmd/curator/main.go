package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/cmd/handler/cache"
	"github.com/nthnca/curator/cmd/handler/getphotos"
	"github.com/nthnca/curator/cmd/handler/newphotos"
	"github.com/nthnca/curator/cmd/handler/queue"
	"github.com/nthnca/curator/cmd/handler/stats"
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

	app.Command("new", "Process new photos").Action(
		func(_ *kingpin.ParseContext) error {
			newphotos.Handler()
			return nil
		})
	app.Command("get", "Create script for copying photos").Action(
		func(_ *kingpin.ParseContext) error {
			getphotos.Handler()
			return nil
		})
	app.Command("oldsync", "Sync photos on disk to the cloud").Action(
		func(_ *kingpin.ParseContext) error {
			update.Handler()
			return nil
		})
	app.Command("cache", "Update datastore caches").Action(
		func(_ *kingpin.ParseContext) error {
			cache.Handler()
			return nil
		})
	app.Command("queue", "queue more work items").Action(
		func(_ *kingpin.ParseContext) error {
			queue.Handler()
			return nil
		})
	app.Command("stats", "analyze curator data").Action(
		func(_ *kingpin.ParseContext) error {
			stats.Handler()
			return nil
		})
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
