package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:      "parse",
			ArgsUsage: "[filepath or - to use stdin]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "rules",
					Value: "fixtures/marc_rules.json",
					Usage: "Path to marc rules file",
				},
				cli.StringFlag{
					Name:  "consumer, c",
					Value: "es",
					Usage: "Consumer to use (es, json or title; default is es)",
				},
				cli.StringFlag{
					Name:  "type, t",
					Value: "marc",
					Usage: "Type of file to process (default is marc)",
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

				if c.String("consumer") == "json" {
					p.consumer = &JSONConsumer{out: os.Stdout}
				} else if c.String("consumer") == "title" {
					p.consumer = &TitleConsumer{out: os.Stdout}
				} else {
					client, err := elastic.NewSimpleClient()
					if err != nil {
						return err
					}
					es, err := client.BulkProcessor().Name("MyBackgroundWorker-1").Do(context.Background())
					if err != nil {
						return err
					}
					defer es.Close()
					p.consumer = &ESConsumer{Index: "timdex", RType: "marc", p: es}
				}

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

				out := p.Run()
				<-out
				return nil
			},
		},
		{
			Name: "create",
			Action: func(c *cli.Context) error {
				client, err := elastic.NewSimpleClient()
				if err != nil {
					return err
				}
				ctx := context.Background()
				created, err := client.CreateIndex("timdex").Do(ctx)
				if err != nil {
					return err
				}
				if !created.Acknowledged {
					fmt.Println("Elasticsearch couldn't create the index")
				}
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("Reserved for λ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
