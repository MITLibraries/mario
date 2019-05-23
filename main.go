package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	var debug bool
	var auto bool
	var url, index string
	var v4 bool

	app := cli.NewApp()

	//Global options
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url, u",
			Value:       "http://127.0.0.1:9200",
			Usage:       "URL for the Elasticsearch cluster",
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "index, i",
			Usage:       "Name of the Elasticsearch index",
			Destination: &index,
		},
		cli.BoolFlag{
			Name:        "v4",
			Usage:       "Use AWS v4 signing",
			Destination: &v4,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "ingest",
			Usage:     "Parse and ingest the input file",
			ArgsUsage: "[filepath, use format 's3://bucketname/objectname' for s3]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "rules",
					Value: "config/marc_rules.json",
					Usage: "Path to marc rules file",
				},
				cli.StringFlag{
					Name:  "consumer, c",
					Value: "es",
					Usage: "Consumer to use (es, json or title)",
				},
				cli.StringFlag{
					Name:  "type, t",
					Value: "marc",
					Usage: "Type of file to process",
				},
				cli.StringFlag{
					Name:  "prefix, p",
					Value: "aleph",
					Usage: "Index prefix to use: default is aleph",
				},
				cli.BoolFlag{
					Name:        "debug",
					Usage:       "Output debugging information",
					Destination: &debug,
				},
				cli.BoolFlag{
					Name:        "auto",
					Usage:       "Automatically promote / demote on completion",
					Destination: &auto,
				},
			},
			Action: func(c *cli.Context) error {
				var file io.ReadCloser
				var err error

				inputData := c.Args().Get(0)
				if len(inputData) == 0 {
					return cli.NewExitError("No filepath argument provided", 1)
				}

				if strings.Contains(inputData, "mit01_edsu1") {
					// this is an aleph update. Determine the index that is currently
					// associated with `production` with the aleph prefix and set that
					// as the index name.
					log.Printf("Update file detected.")
					client, err := esClient(url, index, v4)
					if err != nil {
						return err
					}

					old, err := previous(client, c.String("prefix"))
					if err != nil {
						return err
					}
					if len(old) != 1 {
						return errors.New("Multiple or zero indexes match. Unable to determine which index to update")
					}
					index = old[0]
					println("Using exisitng index:", index)

					// disable auto mode as we are using an existing index
					auto = false
				}

				if inputData[0:2] == "s3" {
					s3Info := strings.Split(inputData, "/")
					file, err = getS3Obj(s3Info[2], s3Info[3])
				} else {
					file, err = os.Open(inputData)
				}

				if err != nil {
					return err
				}

				defer file.Close()

				if index == "" {
					t := time.Now().UTC()
					ft := strings.ToLower(t.Format(time.RFC3339))
					index = fmt.Sprintf("%s-%s", c.String("prefix"), ft)
				}

				p := Pipeline{}

				//Configure the pipeline consumer
				if c.String("consumer") == "json" {
					p.consumer = &JSONConsumer{out: os.Stdout}
				} else if c.String("consumer") == "title" {
					p.consumer = &TitleConsumer{out: os.Stdout}
				} else {
					client, err := esClient(url, index, v4)
					if err != nil {
						return err
					}

					// check if requested index exists, if not, create it.
					exists, err := client.IndexExists(index).Do(context.Background())
					if err != nil {
						return err
					}
					if !exists {
						createRecordIndex(client, index)
					}

					es, err := client.BulkProcessor().After(after).Name("IngestWorker-1").Workers(2).Do(context.Background())
					if err != nil {
						return err
					}
					defer es.Close()
					p.consumer = &ESConsumer{Index: index, RType: "Record", p: es}
				}

				//Configure the pipeline input
				if c.String("type") == "marc" {
					p.generator = &MarcGenerator{
						marcfile:  file,
						rulesfile: c.String("rules"),
					}
				} else if c.String("type") == "json" {
					p.generator = &JSONGenerator{file: file}
				} else if c.String("type") == "archives" {
					p.generator = &ArchivesGenerator{archivefile: file}
				} else {
					log.Println("no valid type provided")
				}

				ctr := &Counter{}
				p.Next(ctr)

				out := p.Run()
				<-out

				if debug {
					log.Printf("Total records ingested: %d", ctr.Count)
				}

				if auto {
					client, err := esClient(url, index, v4)
					if err != nil {
						return err
					}

					log.Printf("Automatic mode detected")
					// retrieve old indexes with supplied prefix
					old, err := previous(client, c.String("prefix"))
					if err != nil {
						return err
					}

					// demote old indexes
					for _, i := range old {
						demote(client, i)
					}

					// promote new index
					promote(client, index)
				}

				return nil
			},
		},
		{
			Name:  "indexes",
			Usage: "List Elasticsearch indexes",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				indexes(client)
				return nil
			},
		},
		{
			Name:  "aliases",
			Usage: "List Elasticsearch aliases and associated indexes",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}

				aliases, err := aliases(client)
				if err != nil {
					return err
				}

				for _, a := range aliases {
					fmt.Printf("Alias: %s \n"+
						"  Index: %s \n\n", a.Alias, a.Index)
				}
				if len(aliases) == 0 {
					fmt.Printf("No aliases found.")
				}
				return nil
			},
		},
		{
			Name:  "ping",
			Usage: "Ping Elasticsearch",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				ping(client, url)
				return nil
			},
		},
		{
			Name:     "delete",
			Usage:    "Delete an Elasticsearch index",
			Category: "Index actions",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				delete(client, index)
				return nil
			},
		},
		{
			Name:     "promote",
			Usage:    "Promote Elasticsearch alias to prod",
			Category: "Index actions",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				promote(client, index)
				return nil
			},
		},
		{
			Name:     "demote",
			Usage:    "Demote Elasticsearch alias from prod",
			Category: "Index actions",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				demote(client, index)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
