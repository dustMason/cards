package ui

import (
	"io/ioutil"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func ColorizedFileContents(file string) (string, error) {
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
