package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nthnca/curator/pkg/action/fsckphotos"
	"github.com/nthnca/curator/pkg/action/getphotos"
	"github.com/nthnca/curator/pkg/action/mutatephotos"
	"github.com/nthnca/curator/pkg/action/newphotos"
	"github.com/nthnca/curator/pkg/action/statphotos"
	"github.com/nthnca/curator/pkg/config"
	objectstore "github.com/nthnca/object-store"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var actual bool
	app := kingpin.New(
		"curator",
		"Photo storage and organization tool")
	app.UsageWriter(os.Stdout)
	app.Flag("go", "Actually make modifications").BoolVar(&actual)

	{
		// Mutate Photos
		opts := mutatephotos.Options{}
		cmd := app.Command("mutate", "Mutate")
		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				opts.Ctx, opts.Storage, opts.ObjStore, opts.Cfg, err = setup()
				if err != nil {
					return err
				}
				opts.DryRun = !actual
				return mutatephotos.Do(&opts)
			})
		cmd.Flag("add", "Labels to add").Short('a').StringsVar(&opts.Tags.A)
		cmd.Flag("remove", "Labels to remove").Short('r').StringsVar(&opts.Tags.B)
	}

	{
		// Get Photos
		opts := getphotos.Options{}
		cmd := app.Command("get", "Create script for copying photos")
		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				opts.Ctx, opts.Storage, opts.ObjStore, opts.Cfg, err = setup()
				if err != nil {
					return err
				}
				return getphotos.Do(&opts)
			})
		cmd.Flag("filter", "description").StringVar(&opts.Filter)
		cmd.Flag("max", "The maximum number of results to return").IntVar(&opts.Max)
		cmd.Flag("has", "Has labels").StringsVar(&opts.Tags.A)
		cmd.Flag("not", "Not labels").StringsVar(&opts.Tags.B)
	}

	{
		// Get Stats
		opts := statphotos.Options{}
		app.Command("stats", "analyze curator data").Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				opts.Ctx, opts.Storage, opts.ObjStore, opts.Cfg, err = setup()
				if err != nil {
					return err
				}
				return statphotos.Do(&opts)
			})
	}

	{
		// Validate Photo System
		opts := fsckphotos.Options{}
		app.Command("fsck", "Validate photos are intact").Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				opts.Ctx, opts.Storage, opts.ObjStore, opts.Cfg, err = setup()
				if err != nil {
					return err
				}
				return fsckphotos.Do(&opts)
			})
	}

	{
		// Process New Photos
		opts := newphotos.Options{}
		app.Command("new", "Process new photos").Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				opts.Ctx, opts.Storage, opts.ObjStore, opts.Cfg, err = setup()
				if err != nil {
					return err
				}
				opts.DryRun = !actual
				return newphotos.Do(&opts)
			})
	}

	{
		// View the configuration
		app.Command("config", "Print the current configuration").Action(
			func(_ *kingpin.ParseContext) error {
				var err error
				fmt.Printf(config.New().ToString())
				return err
			})
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func setup() (context.Context, *storage.Client, *objectstore.ObjectStore, *config.Config, error) {
	ctx := context.Background()

	cfg := config.New()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create GCS client: %v", err)
	}

	os, err := objectstore.New(ctx, client, cfg.MetadataBucket(), cfg.MetadataPath())
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create ObjectStore client: %v", err)
	}

	return ctx, client, os, cfg, nil
}
