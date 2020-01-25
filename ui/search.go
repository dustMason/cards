package ui

import (
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/dustmason/cards/actions"
)

type SearchPage struct {
	Page *tview.Flex

	textView *tview.TextView
	table    *tview.Table
	app      *tview.Application
	dir      string
}

func NewSearchPage(app *tview.Application, dir string, events *Events) *SearchPage {
	// TODO
	// - highlight matching parts of text. might be hard given ANSI...

	sp := &SearchPage{
		textView: nil,
		table:    nil,
		app:      app,
		dir:      dir,
	}

	textView := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	textView.SetBorderPadding(1, 1, 2, 2)

	table := tview.NewTable()
	table.SetBorders(false)
	table.SetBorder(true)
	table.SetSelectable(true, true)
	table.SetBorderPadding(0, 0, 1, 1)
	changed := func(row int, column int) {
		cell := table.GetCell(row, column)
		filename := filepath.Join(dir, cell.Text)
		colored, err := ColorizedFileContents(filename)
		if err != nil {
			panic(err)
		}
		textView.SetText(tview.TranslateANSI(colored))
	}
	table.SetSelectionChangedFunc(changed)

	inputField := tview.NewInputField()
	cardResults, _ := actions.CardResults(dir)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			i, _ := table.GetSelection()
			if i > 0 {
				table.Select(i-1, 0)
			}
		case tcell.KeyDown:
			i, _ := table.GetSelection()
			table.Select(i+1, 0)
		case tcell.KeyEscape:
			events.Emit("hide:Search")
		case tcell.KeyEnter:
			row, column := table.GetSelection()
			cell := table.GetCell(row, column)
			filename := filepath.Join(dir, cell.Text)
			app.Stop()
			_ = actions.Edit(filename)
		default:
			matches := actions.Fuzzy(cardResults, inputField.GetText())
			table.Clear()
			for i, r := range matches {
				table.SetCellSimple(i, 0, cardResults[r.Index].Name)
			}
			table.Select(0, 0)
			changed(0, 0)
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(table, 0, 1, false).
			AddItem(inputField, 2, 1, true), 40, 1, false).
		AddItem(textView, 0, 1, false)

	events.On("show:Search", func() {
		app.SetFocus(inputField)
	})

	sp.Page = flex
	return sp
}
