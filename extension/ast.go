package extension

import (
	"github.com/yuin/goldmark/ast"
)

// A Aside struct represents an aside with Markdown text.
type Aside struct {
	ast.BaseBlock
}

// Dump implements Node.Dump .
func (n *Aside) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindAside is a NodeKind of the Aside node.
var KindAside = ast.NewNodeKind("Aside")

// Kind implements Node.Kind.
func (n *Aside) Kind() ast.NodeKind {
	return KindAside
}

// NewAside returns a new Blockquote node.
func NewAside() *Aside {
	return &Aside{
		BaseBlock: ast.BaseBlock{},
	}
}
