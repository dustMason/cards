package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			d, err := cardsDir()
			if err != nil {
				return err
			}
			Render(d)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "search",
				Usage:   "find and edit a card via grep-style search",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					d, err := cardsDir()
					if err != nil {
						return err
					}
					file, err := Search(d)
					if err != nil {
						return err
					}
					return editFile(file)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func cardsDir() (string, error) {
	me, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(me.HomeDir, ".cards"), nil
}
