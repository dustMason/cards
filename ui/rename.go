package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/dustmason/cards/actions"
)

type RenamePage struct {
	Page *tview.Primitive
}

func NewRenamePage(app *tview.Application, dir string, events *Events, file func() string) *RenamePage {
	rp := &RenamePage{}
	inputField := tview.NewInputField().SetLabel("Enter a new filename: ")
	inputField.SetBorder(true)
	inputField.SetBorderPadding(0, 0, 1, 1)
	inputField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			n := inputField.GetText()
			err := actions.Rename(dir, file(), n)
			if err != nil {
				panic(err)
			}
		}
		events.Emit("hide:Rename")
	})

	events.On("show:Rename", func() {
		inputField.SetText("")
		app.SetFocus(inputField)
	})

	page := modal(inputField, 60, 3)
	rp.Page = &page
	return rp
}

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, false).
			AddItem(nil, 0, 1, false), width, 1, false).
		AddItem(nil, 0, 1, false)
}
