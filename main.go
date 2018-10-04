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
					Usage: "Consumer to use (es or json, default is es)",
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

				Process(file, c.String("rules"), c.String("consumer"))
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
