package page

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
)

// TODO: factor out the common stuff between RenderablePage and PageTemplate?

type RenderablePage struct {
	Template *template.Template

	Meta        *Meta
	Stylesheets []Stylesheet
	Scripts     []Script
	Favicon     *Favicon

	Title   string
	Content []interface{}
}

type ContentItem interface {
	Initialize(*RenderablePage)
	Render(io.Writer) error
}

type PageTemplate struct {
	Template template.Template

	Meta        *Meta
	Stylesheets []Stylesheet
	Scripts     []Script
	Favicon     *Favicon
}

func Template(templatesFolder string, htmlTemplateContent string) (*PageTemplate, error) {
	template, err := template.New("").Funcs(map[string]any{
		"eval": eval,
	}).Parse(htmlTemplateContent)
	if err != nil {
		return nil, err
	}

	return &PageTemplate{
		Template: *template,
	}, nil
}

func eval(r ContentItem) (template.HTML, error) {
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	err := r.Render(writer)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}

	return template.HTML(buf.String()), err
}

func (t *PageTemplate) Instantiate(content ...interface{}) RenderablePage {
	stylesheets := make([]Stylesheet, 0)
	for _, stylesheet := range t.Stylesheets {
		stylesheets = append(stylesheets, stylesheet)
	}

	scripts := make([]Script, 0)
	for _, script := range t.Scripts {
		scripts = append(scripts, script)
	}

	var meta *Meta = nil
	if t.Meta != nil {
		metaCopy := *t.Meta
		meta = &metaCopy
	}

	var favicon *Favicon = nil
	if t.Favicon != nil {
		faviconCopy := *t.Favicon
		favicon = &faviconCopy
	}

	renderablePage := RenderablePage{
		Template: &t.Template,
		Content:  content,
		Favicon:  favicon,
		Meta:     meta,
	}

	initializeItems(&renderablePage, renderablePage.Content)

	return renderablePage
}

func initializeItems(rp *RenderablePage, content []interface{}) {
	for _, item := range content {
		if i, ok := item.(ContentItem); ok {
			i.Initialize(rp)
		}
	}
}

func (rt RenderablePage) Render(w io.Writer) error {
	return rt.Template.Execute(w, rt)
}

func (rt *RenderablePage) AppendStylesheet(s Stylesheet) {
	for _, existingStyle := range rt.Stylesheets {
		if s == existingStyle {
			return
		}
	}

	rt.Stylesheets = append(rt.Stylesheets, s)
}

func (rt RenderablePage) AppendScript(s Script) {
	for _, existingScript := range rt.Scripts {
		if s == existingScript {
			return
		}
	}

	rt.Scripts = append(rt.Scripts, s)
}

func renderItems(w io.Writer, items []interface{}) error {
	for _, item := range items {
		switch i := item.(type) {
		case (ContentItem):
			if err := i.Render(w); err != nil {
				return err
			}
		case string:
			if _, err := w.Write([]byte(i)); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Invalid type")
		}
	}
	return nil
}
