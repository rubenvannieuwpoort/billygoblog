package page

import (
	"io"
)

func Header(content ...interface{}) *HeaderItem {
	return &HeaderItem{Content: content}
}

type HeaderItem struct {
	Content []interface{}
}

func (i *HeaderItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, i.Content)
}

func (i *HeaderItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<header>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</header>")); err != nil {
		return err
	}
	return nil
}
