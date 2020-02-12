package xmlapi

import (
	"fmt"
	"io"
	"text/template"
)

func generateText(w io.Writer, dir string, data interface{}, fn ...string) (err error) {
	var files []string
	for _, f := range fn {
		files = append(files, fmt.Sprintf("%s/templates/%s.tmpl", dir, f))
	}

	tmpls := template.Must(template.ParseFiles(files...))
	err = tmpls.ExecuteTemplate(w, "layout", data)

	return
}
