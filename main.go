package main

import (
	"context"
	"log"
	"os"

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
			Value:       "timdex",
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
			Name:      "parse",
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
			Name:  "create",
			Usage: "Create an Elasticsearch index",
			Action: func(c *cli.Context) error {
				client, err := esClient(url, index, v4)
				if err != nil {
					return err
				}
				createRecordIndex(client, index)
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "Delete an Elasticsearch index",
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
				log.Println("Index deleted")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
