package page

import (
	"io"
)

func Bold(content ...interface{}) *BoldItem {
	return &BoldItem{Content: content}
}

type BoldItem struct {
	Content []interface{}
}

func (b *BoldItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, b.Content)
}

func (i *BoldItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<b>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</b>")); err != nil {
		return err
	}
	return nil
}
