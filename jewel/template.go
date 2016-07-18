package jewel

import (
	"html/template"
	"io"
)

func RenderTplFile(wr io.Writer, tplName string, model interface{}) {
	t, err := template.ParseFiles(tplName)
	if err != nil {
		panic(err)
	}
	t.Execute(wr, model)
}
