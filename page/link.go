package page

import (
	"fmt"
	"io"
)

func Link(url string, content ...interface{}) *LinkItem {
	return &LinkItem{URL: url, Content: content}
}

type LinkItem struct {
	URL     string
	Content []interface{}
}

func (l *LinkItem) Initialize(rp *RenderablePage) {
	initializeItems(rp, l.Content)
}

func (l *LinkItem) Render(w io.Writer) error {
	start := fmt.Sprintf("<a href=\"%s\">", l.URL)
	if _, err := w.Write([]byte(start)); err != nil {
		return err
	}
	if err := renderItems(w, l.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</a>")); err != nil {
		return err
	}
	return nil
}
