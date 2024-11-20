package extension

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type asides struct {
}

// Asides is an extension that allow you to use asides like '< this text is placed in an aside' .
var Asides = &asides{}

func (e *asides) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewAsideParser(), 850),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewAsideRenderer(), 500),
	))
}

func (r *AsideRenderer) renderAside(
	w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<aside>")
	} else {
		_, _ = w.WriteString("</aside>\n")
	}
	return gast.WalkContinue, nil
}
