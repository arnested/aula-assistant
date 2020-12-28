package main

import (
	"io"
	"text/template"
)

func response(w io.Writer, message string) {
	tpl := `{
  "fulfillment_response": {
    "messages": [
      {
	"text": {
	  "text": [
	    "{{ . }}"
	  ]
	}
      }
    ]
  }
}`
	t, _ := template.New("template").Parse(tpl)
	_ = t.Execute(w, message)
}
