package main

import (
	"io"
	"text/template"
)

func response(w io.Writer, message string) {
	tpl := `{
  "google": {
    "expectUserResponse": true,
    "richResponse": {
      "items": [
	{
	  "simpleResponse": {
	    "textToSpeech": "{{ . }}"
	  }
	}
      ]
    }
  }
}`
	t, _ := template.New("template").Parse(tpl)
	_ = t.Execute(w, message)
}
