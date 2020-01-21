package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/dustmason/cards/browse"
	"github.com/dustmason/cards/search"
)

const DefaultEditor = "vim"

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Usage:   "create a new card",
				Aliases: []string{"n"},
				Action: func(c *cli.Context) error {
					d, err := cardsDir()
					if err != nil {
						return err
					}
					suffix := c.Args().Get(0)
					if suffix == "" {
						suffix = "md"
					}
					path := filepath.Join(d, fmt.Sprintf("%v.%v", time.Now().Format("2006-01-02-15-04-05"), suffix))
					return editFile(path)
				},
			},

			{
				Name:    "search",
				Usage:   "find and edit a card via grep-style search",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					d, err := cardsDir()
					if err != nil {
						return err
					}
					file, err := search.Search(d)
					if err != nil {
						return err
					}
					return editFile(file)
				},
			},

			{
				Name:    "browse",
				Usage:   "navigate, view and edit all cards",
				Aliases: []string{"b"},
				Action: func(c *cli.Context) error {
					d, err := cardsDir()
					if err != nil {
						return err
					}
					intent, file, err := browse.Browse(d)
					if err != nil {
						return err
					}
					switch intent {
					case browse.Edit:
						return editFile(file)
					case browse.Archive:
						fmt.Println("move to archive", file)
					case browse.Copy:
						pbcopy, err := exec.LookPath("pbcopy")
						if err != nil {
							return err
						}
						c := exec.Command(pbcopy)
						f, err := os.Open(file)
						if err != nil {
							return err
						}
						c.Stdin = f
						return c.Run()
					}
					return nil
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

func editFile(path string) error {
	if path == "" {
		return nil
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}
	return run(editor, []string{path})
}

func run(cmd string, args []string) error {
	executable, execErr := exec.LookPath(cmd)
	if execErr != nil {
		return execErr
	}
	c := []string{cmd}
	c = append(c, args...)
	return syscall.Exec(executable, c, os.Environ())
}
