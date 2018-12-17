package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

func main() {
	var debug bool
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
			ArgsUsage: "[filepath or - to use stdin]",
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
				cli.BoolFlag{
					Name:        "debug",
					Usage:       "Output debugging information",
					Destination: &debug,
				},
			},
			Action: func(c *cli.Context) error {
				var file *os.File
				var err error

				if index == "" {
					t := time.Now().UTC()
					ft := strings.ToLower(t.Format(time.RFC3339))
					index = fmt.Sprintf("%s-%s", "aleph", ft)
				}

				// if a file path is passed as a flag
				if c.Args().Get(0) != "-" {
					// Open the file.
					file, err = os.Open(c.Args().Get(0))
				} else {
					// otherwise try to use stdin
					file = os.Stdin
				}

				if err != nil {
					return err
				}

				defer file.Close()

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

					es, err := client.BulkProcessor().Name("IngestWorker-1").Do(context.Background())
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
				ctx := context.Background()
				indexes, err := client.CatIndices().Do(ctx)
				if err != nil {
					return err
				}
				for _, i := range indexes {
					fmt.Printf("Name: %s \n"+
						"  DocsCount: %d \n"+
						"  Health: %s \n"+
						"  Status: %s \n"+
						"  UUID: %s \n"+
						"  StoreSize: %s \n\n",
						i.Index, i.DocsCount, i.Health, i.Status, i.UUID, i.StoreSize)
				}
				if len(indexes) == 0 {
					fmt.Printf("No indexes found.")
				}
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
				ctx := context.Background()
				aliases, err := client.CatAliases().Do(ctx)
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
				ctx := context.Background()
				ping, code, err := client.Ping(url).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Response code: %d \n"+
					"Name: %s \n"+
					"Cluster Name: %s \n"+
					"Tag line: %s \n"+
					"Version: %s \n"+
					"BuildHash: %s \n"+
					"LuceneVersion: %s \n",
					code, ping.Name, ping.ClusterName, ping.TagLine, ping.Version.Number,
					ping.Version.BuildHash, ping.Version.LuceneVersion)
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
				ctx := context.Background()
				_, err = client.DeleteIndex(index).Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Index deleted")
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
				ctx := context.Background()
				_, err = client.Alias().Add(index, "production").Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Index %s promoted.", index)
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
				ctx := context.Background()
				_, err = client.Alias().Remove(index, "production").Do(ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Index %s demoted.", index)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
