package page

import (
	"io"

	"github.com/graemephi/goldmark-qjs-katex/katex"
)

func InlineMath(text string) *InlineMathItem {
	return &InlineMathItem{Text: text}
}

type InlineMathItem struct {
	Text string
}

func (i *InlineMathItem) Initialize(rp *RenderablePage) {
	rp.AppendStylesheet(Stylesheet{
		Uri: "katex.min.css",
	})
}

func (i *InlineMathItem) Render(w io.Writer) error {
	return katex.RenderTo(w, []byte(i.Text), katex.Inline)
}
