package internal

import (
	"io"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

func NewMinifier(input io.Reader, mediaType string) Minifier {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/javascript", js.Minify)
	output := m.Reader(mediaType, input)

	return Minifier{output: output}
}

type Minifier struct {
	output io.Reader
}

func (m *Minifier) Reader() io.Reader {
	return m.output
}
