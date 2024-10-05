package page

import (
	"io"
)

func Code(code string) *CodeItem {
	return &CodeItem{Code: code}
}

type CodeItem struct {
	Code string
}

func (c *CodeItem) Initialize(rp *RenderablePage) {}

func (c *CodeItem) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<pre><code>")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(c.Code)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</code></pre>")); err != nil {
		return err
	}
	return nil
}
