package ui

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/dustmason/cards/actions"
)

type inputCapHandler func(event *tcell.EventKey) *tcell.EventKey

type BrowsePage struct {
	Page *tview.Flex

	textView *tview.TextView
	table    *tview.Table
	app      *tview.Application
	dir      string

	SelectedFile string
}

func NewBrowsePage(app *tview.Application, dir string, events *Events) *BrowsePage {
	table := tview.NewTable()
	textView := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	textView.SetBorderPadding(1, 1, 2, 2)
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(table)
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(table, 40, 1, true).
		AddItem(textView, 0, 1, false)

	bp := &BrowsePage{
		Page:     flex,
		textView: textView,
		table:    table,
		app:      app,
		dir:      dir,
	}

	table.SetBorders(false)
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)

	fileCount := bp.Draw()

	table.SetSelectable(true, true).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	}).SetSelectionChangedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		filename := filepath.Join(dir, cell.Text)
		bp.SelectedFile = cell.Text
		colored, err := ColorizedFileContents(filename)
		if err != nil {
			panic(err)
		}
		textView.SetText(tview.TranslateANSI(colored))
	}).Select(fileCount-1, 0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			err := actions.Archive(dir, bp.SelectedFile)
			if err != nil {
				panic(err)
			}
			fileCount = bp.Draw()
		case 'c':
			_ = actions.Pbcopy(filepath.Join(dir, bp.SelectedFile))
			app.Stop()
		case 'p':
			app.Stop()
			source, err := ioutil.ReadFile(filepath.Join(dir, bp.SelectedFile))
			if err != nil {
				panic(err)
			}
			actions.Present(string(source))
		case 'q':
			app.Stop()
		case 'n':
			app.Stop()
			_ = actions.Create(dir)
		case 'r':
			events.Emit("show:Rename")
		case '/':
			events.Emit("show:Search")
		}
		switch event.Key() {
		case tcell.KeyEnter:
			app.Stop()
			_ = actions.Edit(filepath.Join(dir, bp.SelectedFile))
		case tcell.KeyTab:
			app.SetFocus(textView)
		}
		return event
	})

	return bp
}

func (bp BrowsePage) Draw() int {
	files, _ := ioutil.ReadDir(bp.dir)
	fileCount := 0
	for _, file := range files {
		name := file.Name()
		if !strings.HasPrefix(name, ".") {
			cell := tview.NewTableCell(file.Name()).SetMaxWidth(40).SetExpansion(1)
			bp.table.SetCell(fileCount, 0, cell)
			fileCount++
		}
	}
	return fileCount
}
