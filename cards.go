package main

import (
	"errors"
	"os"

	"github.com/rivo/tview"

	"github.com/dustmason/cards/ui"
)

func main() {
	d, err := cardsDir()
	if err != nil {
		panic(err)
	}
	showUI(d)
}

func cardsDir() (string, error) {
	d, ok := os.LookupEnv("CARDS_DIR")
	if !ok {
		return "", errors.New("please set the CARDS_DIR env var to the directory you wish to store notes")
	}
	return d, nil
}

func showUI(dir string) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	events := ui.NewEvents()
	bp := &ui.BrowsePage{}

	events.On("show:Rename", func(s string) {
		pages.ShowPage("Rename")
	})
	events.On("show:Search", func(s string) {
		pages.SwitchToPage("Search")
	})
	events.On("hide:Search", func(s string) {
		pages.SwitchToPage("Browse")
		if s != "" {
			bp.Select(s)
		}
	})

	bp = ui.NewBrowsePage(app, dir, events)
	pages.AddPage("Browse", bp.Page, true, true)

	rp := ui.NewRenamePage(app, dir, events)
	pages.AddPage("Rename", *rp.Page, true, false)
	events.On("hide:Rename", func(s string) {
		pages.HidePage("Rename")
		bp.Reload()
	})

	sp := ui.NewSearchPage(app, dir, events)
	pages.AddPage("Search", sp.Page, true, false)

	app.SetRoot(pages, true).SetFocus(pages)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
