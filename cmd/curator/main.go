package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nthnca/curator/pkg/action/fsckphotos"
	"github.com/nthnca/curator/pkg/action/getphotos"
	"github.com/nthnca/curator/pkg/action/mutatephotos"
	"github.com/nthnca/curator/pkg/action/newphotos"
	"github.com/nthnca/curator/pkg/action/statphotos"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	actual bool
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := kingpin.New(
		"curator",
		"Photo storage and organization tool")
	app.UsageWriter(os.Stdout)
	app.Flag("go", "Actually make modifications").BoolVar(&actual)

	// Mutate Photos
	{
		opts := mutatephotos.Options{}
		cmd := app.Command("mutate", "Mutate")
		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				opts.DryRun = !actual
				mutatephotos.Do(&opts)
				return nil
			})
		cmd.Flag("add", "Labels to add").Short('a').StringsVar(&opts.Tags.A)
		cmd.Flag("remove", "Labels to remove").Short('r').StringsVar(&opts.Tags.B)
	}

	// Get Photos
	{
		getOpt := getphotos.Options{}
		cmd := app.Command("get", "Create script for copying photos")
		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				getphotos.Do(&getOpt)
				return nil
			})
		cmd.Flag("filter", "description").StringVar(&getOpt.Filter)
		cmd.Flag("max", "The maximum number of results to return").IntVar(&getOpt.Max)
		cmd.Flag("has", "Has labels").StringsVar(&getOpt.Tags.A)
		cmd.Flag("not", "Not labels").StringsVar(&getOpt.Tags.B)
	}

	// Get Stats
	app.Command("stats", "analyze curator data").Action(
		func(_ *kingpin.ParseContext) error {
			statphotos.Do()
			return nil
		})

	// Validate Photo System
	app.Command("fsck", "Validate photos are intact").Action(
		func(_ *kingpin.ParseContext) error {
			fsckphotos.Do()
			return nil
		})

	// Process New Photos
	app.Command("new", "Process new photos").Action(
		func(_ *kingpin.ParseContext) error {
			opts := newphotos.Options{}
			opts.DryRun = !actual
			newphotos.Do(&opts)
			return nil
		})

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
