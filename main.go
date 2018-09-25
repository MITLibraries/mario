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
			Name:      "parse",
			ArgsUsage: "[filepath or - to use stdin]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "rules",
					Value: "fixtures/marc_rules.json",
					Usage: "Path to marc rules file",
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

				marc.Process(file, c.String("rules"))
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
