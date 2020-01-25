package main

import (
	"os/user"
	"path/filepath"

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
	me, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(me.HomeDir, ".cards"), nil
}

func showUI(dir string) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	events := ui.NewEvents()

	events.On("show:Rename", func() {
		pages.ShowPage("Rename")
	})
	events.On("show:Search", func() {
		pages.SwitchToPage("Search")
	})
	events.On("hide:Search", func() {
		pages.SwitchToPage("Browse")
	})

	bp := ui.NewBrowsePage(app, dir, events)
	pages.AddPage("Browse", bp.Page, true, true)

	rp := ui.NewRenamePage(app, dir, events, func() string {
		return bp.SelectedFile
	})
	pages.AddPage("Rename", *rp.Page, true, false)
	events.On("hide:Rename", func() {
		pages.HidePage("Rename")
		_ = bp.Draw()
	})

	sp := ui.NewSearchPage(app, dir, events)
	pages.AddPage("Search", sp.Page, true, false)

	app.SetRoot(pages, true).SetFocus(pages)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
