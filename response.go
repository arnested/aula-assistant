package main

import (
	"io"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
)

type message string

func (msg message) Strip() string {
	p := bluemonday.NewPolicy()

	return p.Sanitize(string(msg))
}

func response(w io.Writer, msg string) {
	tpl := `{
  "prompt": {
    "firstSimple": {
      "speech": "{{ . }}",
      "text": "{{ .Strip }}"
    }
  }
}`
	t, _ := template.New("template").Parse(tpl)
	_ = t.Execute(w, message(msg))
}
