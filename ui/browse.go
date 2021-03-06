package ui

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/dustmason/cards/actions"
)

type inputCapHandler func(event *tcell.EventKey) *tcell.EventKey

var linkPattern = regexp.MustCompile(`\[\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]\]`)

type BrowsePage struct {
	Page *tview.Flex

	textView *tview.TextView
	table    *tview.Table
	app      *tview.Application
	dir      string
	files    []string

	SelectedFile string
}

func NewBrowsePage(app *tview.Application, dir string, events *Events) *BrowsePage {
	table := tview.NewTable()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorderPadding(1, 1, 2, 2)

	rotateLinkSelection := func(dir int) {
		currentSelection := "0"
		if len(textView.GetHighlights()) > 0 {
			currentSelection = textView.GetHighlights()[0]
			regionId, _ := strconv.Atoi(currentSelection)
			regionId += dir
			if textView.GetRegionText(strconv.Itoa(regionId)) == "" {
				currentSelection = "0"
			} else {
				currentSelection = strconv.Itoa(regionId)
			}
		}
		textView.Highlight(currentSelection).ScrollToHighlight()
	}

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'J':
			rotateLinkSelection(1)
		case 'K':
			rotateLinkSelection(-1)
		case 'q':
			app.Stop()
		}

		switch event.Key() {
		case tcell.KeyEnter:
			if len(textView.GetHighlights()) > 0 {
				currentSelection := textView.GetHighlights()[0]
				link := textView.GetRegionText(currentSelection)
				events.Emit("show:Search", link)
			}
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
	bp.Reload()

	table.SetBorders(false)
	table.SetBorder(true)
	table.SetBorderPadding(0, 0, 1, 1)

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
		withRegions := addRegions(tview.Escape(colored))
		textView.SetText(tview.TranslateANSI(withRegions))
	}).Select(bp.table.GetRowCount()-1, 0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			err := actions.Archive(dir, bp.SelectedFile)
			if err != nil {
				panic(err)
			}
			bp.Reload()
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
			events.Emit("show:Rename", bp.SelectedFile)
		case '/':
			events.Emit("show:Search", "")
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

func addRegions(text string) string {
	var count int
	return linkPattern.ReplaceAllStringFunc(text, func(f string) string {
		result := fmt.Sprintf(`["%d"]%s[""]`, count, f)
		count++
		return result
	})
}

func (bp *BrowsePage) Reload() {
	files, _ := ioutil.ReadDir(bp.dir)
	bp.files = []string{}
	bp.table.Clear()
	for _, file := range files {
		name := file.Name()
		if !strings.HasPrefix(name, ".") {
			bp.files = append(bp.files, name)
			cell := tview.NewTableCell(name).SetMaxWidth(40).SetExpansion(1)
			bp.table.SetCell(bp.table.GetRowCount(), 0, cell)
		}
	}
}

func (bp *BrowsePage) Select(filename string) {
	for i, name := range bp.files {
		if name == filename {
			bp.table.Select(i, 0)
			return
		}
	}
}
