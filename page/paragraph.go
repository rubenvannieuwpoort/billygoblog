package page

import (
	"io"
)

func Paragraph(content ...interface{}) *ParagraphItem {
	return &ParagraphItem{Content: content}
}

type ParagraphItem struct {
	Content []interface{}
}

func (p *ParagraphItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, p.Content)
}

func (p *ParagraphItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<p>")); err != nil {
		return err
	}
	if err := renderItems(w, p.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</p>")); err != nil {
		return err
	}
	return nil
}
