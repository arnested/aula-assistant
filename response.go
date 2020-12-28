package main

import (
	"io"
	"text/template"
)

func response(w io.Writer, message string) {
	tpl := `{
  "prompt": {
    "firstSimple": {
      "speech": "{{ . }}",
      "text": "{{ . }}"
    }
  }
}`
	t, _ := template.New("template").Parse(tpl)
	_ = t.Execute(w, message)
}
