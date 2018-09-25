package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitlibraries/mario/parsers"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "parse",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "rules",
					Value: "fixtures/marc_rules.json",
					Usage: "Path to marc rules file",
				},
			},
			Action: func(c *cli.Context) error {
				marc.Process(os.Stdin, c.String("rules"))
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
