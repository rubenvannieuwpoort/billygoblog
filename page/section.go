package page

import (
	"io"
)

func Section(content ...interface{}) *SectionItem {
	return &SectionItem{Content: content}
}

type SectionItem struct {
	Content []interface{}
}

func (i *SectionItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, i.Content)
}

func (i *SectionItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<section>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</section>")); err != nil {
		return err
	}
	return nil
}
