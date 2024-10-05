package page

import "io"

func H1(content ...interface{}) *H1Item {
	return &H1Item{Content: content}
}

type H1Item struct {
	Content []interface{}
}

func (i *H1Item) Initialize(rp *RenderablePage) {
	initializeItems(rp, i.Content)
}

func (i *H1Item) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<h1>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</h1>")); err != nil {
		return err
	}
	return nil
}

func H2(content ...interface{}) *H2Item {
	return &H2Item{Content: content}
}

type H2Item struct {
	Content []interface{}
}

func (i *H2Item) Initialize(rp *RenderablePage) {
	initializeItems(rp, i.Content)
}

func (i *H2Item) Render(w io.Writer) error {
	if _, err := w.Write([]byte("<h2>")); err != nil {
		return err
	}
	if err := renderItems(w, i.Content); err != nil {
		return err
	}
	if _, err := w.Write([]byte("</h2>")); err != nil {
		return err
	}
	return nil
}
