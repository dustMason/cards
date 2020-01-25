package actions

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/sahilm/fuzzy"
)

const bold = "\033[1m%s\033[0m"

type cardResult struct {
	Name     string
	Contents string
}

type cardResults []cardResult

func (cr cardResults) String(i int) string {
	return cr[i].Contents
}

func (cr cardResults) Len() int {
	return len(cr)
}

func CardResults(dir string) (cardResults, error) {
	var allContents cardResults

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() && !strings.HasPrefix(name, ".") {
			source, err := ioutil.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			allContents = append(allContents, cardResult{Name: name, Contents: string(source)})
		}
	}

	return allContents, nil
}

func Fuzzy(contents cardResults, needle string) fuzzy.Matches {
	return fuzzy.FindFrom(needle, contents)
}
