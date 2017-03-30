package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/config"
	"github.com/nthnca/curator/handler/analyze"
	"github.com/nthnca/curator/handler/sync"
	"github.com/nthnca/curator/handler/update"

	"github.com/nthnca/easybuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo organizational system that run in Google AppEngine")
	easybuild.RegisterCommands(app, config.Path, config.ProjectID)
	app.Command("sync", "list entries").Action(sync.Handler)
	app.Command("analyze", "list entries").Action(analyze.Handler)
	app.Command("update", "list entries").Action(update.Handler)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
