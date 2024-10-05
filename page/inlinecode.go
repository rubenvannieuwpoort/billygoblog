package page

import (
	"io"
)

func InlineCode(code string) *InlineCodeItem {
	return &InlineCodeItem{Code: code}
}

type InlineCodeItem struct {
	Code string
}

func (c *InlineCodeItem) Initialize(rp *RenderablePage) {}

func (c *InlineCodeItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<code>")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(c.Code)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</code>")); err != nil {
		return err
	}
	return nil
}
