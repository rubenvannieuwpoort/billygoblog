package page

import (
	"io"

	"github.com/graemephi/goldmark-qjs-katex/katex"
)

func Math(text string) *MathItem {
	return &MathItem{Text: text}
}

type MathItem struct {
	Text string
}

func (t *MathItem) Initialize(rp *RenderablePage) {
	rp.AppendStylesheet(Stylesheet{
		Uri: "katex.min.css",
	})
}

func (m *MathItem) Render(w io.Writer) error {
	return katex.RenderTo(w, []byte(m.Text), katex.Display)
}
