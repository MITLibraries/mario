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
			Aliases: 		 []string{"u"},
			Value:       "http://127.0.0.1:9200",
			Usage:       "URL for the Elasticsearch cluster",
			Destination: &url,
		},
		&cli.BoolFlag{
			Name:        "v4",
			Usage:       "Use AWS v4 signing",
			Destination: &v4,
		},
	}

	app.Commands = []*cli.Command{
    // Elasticsearch commands
		{
			Name:     "aliases",
			Usage:    "List Elasticsearch aliases and their associated indexes",
      Category: "Elasticsearch actions",
			Action:   func(c *cli.Context) error {
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
			Usage:    "List all Elasticsearch indexes",
      Category: "Elasticsearch actions",
			Action:   func(c *cli.Context) error {
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
			Usage:    "Ping Elasticsearch",
      Category: "Elasticsearch actions",
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
			Usage:     "Parse and ingest the input file to a new or existing index",
			ArgsUsage: "[filepath, use format 's3://bucketname/objectname' for s3]",
      Category:  "Index actions",
			Flags:     []cli.Flag{
    		&cli.StringFlag{
      		Name:    "index",
    			Aliases: []string{"i"},
    			Usage:   "Name of the Elasticsearch index to ingest to. If not included, will default to a new index named with the source prefix plus timestamp (except for Aleph update files, which are always ingested into the current production Aleph index)",
    		},
				&cli.StringFlag{
					Name:  	  "source",
					Aliases:  []string{"s"},
					Usage: 	  "Source system of metadata file to process. Must be one of [aleph, aspace, dspace, mario]",
					Required: true,
				},
				&cli.StringFlag{
					Name:  	 "consumer",
					Aliases: []string{"c"},
					Value: 	 "es",
					Usage: 	 "Consumer to use. Must be one of [es, json, title, silent]",
				},
				&cli.BoolFlag{
					Name:  "auto",
					Usage: "Automatically promote / demote on completion",
				},
			},
			Action: func(c *cli.Context) error {
				var es *client.ESClient
				config := ingester.Config{
					Filename:  c.Args().Get(0),
					Consumer:  c.String("consumer"),
					Source:    c.String("source"),
					Index:     c.String("index"),
					Promote:   c.Bool("auto"),
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
			Name:     "promote",
			Usage:    "Promote an index to production",
      UsageText: "Demotes the existing production index for the provided prefix, if there is one",
			Category: "Index actions",
			Flags:    []cli.Flag{
    		&cli.StringFlag{
    			Name:     "index",
    			Aliases:  []string{"i"},
    			Usage:    "Name of the Elasticsearch index to promote",
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
			UsageText: "Use the Elasticsearch reindex API to copy one index to another. The doc source must be present in the original index.",
			Category:  "Index actions",
			Flags:     []cli.Flag{
    		&cli.StringFlag{
    			Name:     "index",
    			Aliases: 	[]string{"i"},
    			Usage:    "Name of the Elasticsearch index to copy",
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
			Name:      "delete",
			Usage:     "Delete an index",
			Category:  "Index actions",
			Flags:     []cli.Flag{
    		&cli.StringFlag{
    			Name:     "index",
    			Aliases: 	[]string{"i"},
    			Usage:    "Name of the Elasticsearch index to delete",
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
