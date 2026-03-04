package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"hexlet-boilerplates/gopackage/internal/code"
)

func main() {
	app := &cli.App{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory",
		UsageText: "hexlet-path-size [path] [options]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Usage:   "human-readable sizes (auto-select unit)",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden files and directories",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("DEBUG - Args().Slice(): %#v\n", c.Args().Slice())
			fmt.Printf("DEBUG - Флаг 'human': %v\n", c.Bool("human"))
			fmt.Printf("DEBUG - Флаг 'all': %v\n", c.Bool("all"))

			if c.NArg() == 0 {
				return fmt.Errorf("path is required")
			}

			path := c.Args().Get(0)
			human := c.Bool("human")
			all := c.Bool("all")

			size, err := code.GetSize(path, all)
			if err != nil {
				return err
			}

			formattedSize := code.FormatSize(size, human)
			fmt.Printf("%s\t%s\n", formattedSize, path)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}