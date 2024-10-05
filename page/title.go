package page

import (
	"html/template"
	"io"
)

func Title(text string) *TitleItem {
	return &TitleItem{Text: text}
}

type TitleItem struct {
	Text string
}

func (t *TitleItem) Initialize(rp *RenderablePage) {
	rp.Title = t.Text
}

var titleTemplate = template.Must(template.New("paragraph").Parse(`<h1>{{.Text}}</h1>`))

func (t *TitleItem) Render(w io.Writer) error {
	return titleTemplate.Execute(w, t)
}
