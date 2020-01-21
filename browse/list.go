package browse

import (
	"io/ioutil"
	"path/filepath"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Intent int

const (
	None    = 0
	Edit    = 1
	Archive = 2
	Copy    = 3
)

func Browse(dir string) (Intent, string, error) {
	if err := ui.Init(); err != nil {
		return None, "", err
	}

	l := widgets.NewList()
	l.Title = "cards"

	// TODO accept a flag that sorts by modified at stamp instead of name (created at)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return None, "", err
	}
	l.Rows = []string{}
	for _, file := range files {
		l.Rows = append([]string{file.Name()}, l.Rows...)
	}

	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	contents := widgets.NewParagraph()
	contents.Title = "contents"
	contents.Text = ""

	grid.Set(
		ui.NewCol(1.0/5, l),
		ui.NewCol(1.0-(1.0/5), contents),
	)

	path := func() string {
		return filepath.Join(dir, l.Rows[l.SelectedRow])
	}

	preview := func() {
		c, _ := raw(path()) // TODO syntax highlighting
		contents.Text = c
		ui.Render(contents)
	}

	ui.Render(grid)
	preview()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Enter>":
			return Edit, path(), nil
		case "a":
			ui.Close()
			return Archive, path(), nil
		case "c":
			ui.Close()
			return Copy, path(), nil
		case "q", "<C-c>":
			return None, "", nil
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "g", "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		preview()
		ui.Render(l)
	}
}

func raw(file string) (string, error) {
	out, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// func bat(filename string) (string, error) {
// 	d, err := cardsDir()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	var output bytes.Buffer
// 	bat, execErr := exec.LookPath("bat")
// 	if execErr != nil {
// 		return "", execErr
// 	}
//
// 	cmd := exec.Command(bat, filepath.Join(d, filename))
// 	cmd.Stdout = &output
// 	cmd.Stderr = os.Stderr
// 	err = cmd.Run()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return output.String(), nil
// }
