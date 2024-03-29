package main

import (
	"fmt"
	"github.com/mitlibraries/mario/pkg/client"
	"github.com/mitlibraries/mario/pkg/ingester"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var url string
	var v4 bool

	app := cli.NewApp()

	// Global options
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "url",
			Aliases:     []string{"u"},
			Value:       "http://127.0.0.1:9200",
			Usage:       "URL for the OpenSearch cluster",
			Destination: &url,
		},
		&cli.BoolFlag{
			Name:        "v4",
			Usage:       "Use AWS v4 signing",
			Destination: &v4,
		},
	}

	app.Commands = []*cli.Command{
		// OpenSearch commands
		{
			Name:     "aliases",
			Usage:    "List OpenSearch aliases and their associated indexes",
			Category: "OpenSearch actions",
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				aliases, err := es.Aliases()
				if err != nil {
					return err
				}
				for _, a := range aliases {
					fmt.Printf("Alias: %s\n\tIndex: %s\n\n", a.Alias, a.Index)
				}
				return nil
			},
		},
		{
			Name:     "indexes",
			Usage:    "List all OpenSearch indexes",
			Category: "OpenSearch actions",
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				indexes, err := es.Indexes()
				if err != nil {
					return err
				}
				for _, i := range indexes {
					fmt.Printf("Name: %s\n\tDocuments: %d\n\tHealth: %s\n\tStatus: %s\n\tUUID: %s\n\tSize: %s\n\n", i.Index, i.DocsCount, i.Health, i.Status, i.UUID, i.StoreSize)
				}
				return nil
			},
		},
		{
			Name:     "ping",
			Usage:    "Ping OpenSearch",
			Category: "OpenSearch actions",
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				res, err := es.Ping(url)
				if err != nil {
					return err
				}
				fmt.Printf("Name: %s\nCluster: %s\nVersion: %s\nLucene version: %s", res.Name, res.ClusterName, res.Version.Number, res.Version.LuceneVersion)
				return nil
			},
		},
		// Index-specific commands
		{
			Name:      "ingest",
			Usage:     "Parse and ingest the input file. By default, ingests into the current production index for the provided source.",
			ArgsUsage: "[filepath, use format 's3://bucketname/objectname' for s3]",
			Category:  "Index actions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "source",
					Aliases:  []string{"s"},
					Usage:    "Source system of metadata file to process. Must be one of [alma, aspace, dspace, mario]",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "consumer",
					Aliases: []string{"c"},
					Value:   "es",
					Usage:   "Consumer to use. Must be one of [es, json, title, silent]",
				},
				&cli.BoolFlag{
					Name:  "new",
					Usage: "Create a new index instead of ingesting into the current production index for the source",
				},
				&cli.BoolFlag{
					Name:  "auto",
					Usage: "Automatically promote / demote on completion",
				},
			},
			Action: func(c *cli.Context) error {
				var es *client.ESClient
				config := ingester.Config{
					Filename: c.Args().Get(0),
					Consumer: c.String("consumer"),
					Source:   c.String("source"),
					NewIndex: c.Bool("new"),
					Promote:  c.Bool("auto"),
				}
				log.Printf("Ingesting records from file: %s\n", config.Filename)
				stream, err := ingester.NewStream(config.Filename)
				if err != nil {
					return err
				}
				defer stream.Close()
				if config.Consumer == "es" {
					es, err = client.NewESClient(url, v4)
					if err != nil {
						return err
					}
				}
				ingest := ingester.Ingester{Stream: stream, Client: es}
				err = ingest.Configure(config)
				if err != nil {
					return err
				}
				count, err := ingest.Ingest()
				log.Printf("Total records ingested: %d\n", count)
				return err
			},
		},
		{
			Name:      "promote",
			Usage:     "Promote an index to production",
			UsageText: "Demotes the existing production index for the provided prefix, if there is one",
			Category:  "Index actions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "index",
					Aliases:  []string{"i"},
					Usage:    "Name of the OpenSearch index to promote",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				err = es.Promote(c.String("index"))
				return err
			},
		},
		{
			Name:      "reindex",
			Usage:     "Reindex one index to another index",
			UsageText: "Use the OpenSearch reindex API to copy one index to another. The doc source must be present in the original index.",
			Category:  "Index actions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "index",
					Aliases:  []string{"i"},
					Usage:    "Name of the OpenSearch index to copy",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "destination",
					Aliases:  []string{"d"},
					Usage:    "Name of new index",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				count, err := es.Reindex(c.String("index"), c.String("destination"))
				fmt.Printf("%d documents reindexed\n", count)
				return err
			},
		},
		{
			Name:     "delete",
			Usage:    "Delete an index",
			Category: "Index actions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "index",
					Aliases:  []string{"i"},
					Usage:    "Name of the OpenSearch index to delete",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				es, err := client.NewESClient(url, v4)
				if err != nil {
					return err
				}
				err = es.Delete(c.String("index"))
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
