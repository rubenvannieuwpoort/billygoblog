package page

import (
	"io"
)

func Italic(content ...interface{}) *ItalicItem {
	return &ItalicItem{Content: content}
}

type ItalicItem struct {
	Content []interface{}
}

func (i *ItalicItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, i.Content)
}

func (i *ItalicItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<i>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</i>")); err != nil {
		return err
	}
	return nil
}
