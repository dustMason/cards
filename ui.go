package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const DefaultEditor = "vim"

// TODO bring back borders

func Render(dir string) {
	app := tview.NewApplication()
	textView := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	table := tview.NewTable().SetBorders(false)

	flex := tview.NewFlex().
		AddItem(table, 40, 1, true).
		AddItem(textView, 0, 1, false)

	showBrowser := func() {
		app.SetRoot(flex, true).SetFocus(flex)
	}

	var selectedFile string

	renderFileList := func() {
		files, _ := ioutil.ReadDir(dir)
		for i, file := range files {
			name := file.Name()
			if !strings.HasPrefix(name, ".") {
				table.SetCell(i, 0, tview.NewTableCell(file.Name()).SetAlign(tview.AlignLeft))
			}
		}
	}
	renderFileList()

	showRenameUI := func(dir string, file string) {
		inputField := tview.NewInputField().SetLabel("Enter a new filename: ")
		inputField.SetDoneFunc(func(key tcell.Key) {
			n := inputField.GetText()
			err := rename(dir, file, n)
			if err != nil {
				panic(err)
			}
			showBrowser()
			renderFileList()
		})
		app.SetRoot(inputField, true)
	}

	table.SetSelectable(true, true).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			app.Stop()
			_ = editFile(filepath.Join(dir, selectedFile))
		}
	}).SetSelectionChangedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		filename := filepath.Join(dir, cell.Text)
		selectedFile = cell.Text
		colored, err := render(filename)
		if err != nil {
			panic(err)
		}
		textView.SetText(tview.TranslateANSI(colored))
	}).Select(0, 0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			err := archive(dir, selectedFile)
			if err != nil {
				panic(err)
			}
			renderFileList()
		case 'r':
			showRenameUI(dir, selectedFile)
		case 'c':
			_ = pbcopy(filepath.Join(dir, selectedFile))
			app.Stop()
		case 'n':
			app.Stop()
			_ = createFile(dir)
		}
		return event
	})

	showBrowser()
	err := app.Run()
	if err != nil {
		panic(err)
	}
}

func render(file string) (string, error) {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return "", nil
	}
	w := strings.Builder{}
	lexer := lexers.Match(file)
	lexer = chroma.Coalesce(lexer)
	style := styles.Get("fruity")
	formatter := formatters.Get("terminal")
	iterator, err := lexer.Tokenise(nil, string(source))
	if err != nil {
		return "", err
	}
	err = formatter.Format(&w, style, iterator)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}
