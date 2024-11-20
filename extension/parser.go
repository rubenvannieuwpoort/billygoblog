package extension

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type asideParser struct {
}

var defaultAsideParser = &asideParser{}

// NewAsideParser returns a new BlockParser that
// parses asides.
func NewAsideParser() parser.BlockParser {
	return defaultAsideParser
}

func (a *asideParser) process(reader text.Reader) bool {
	line, _ := reader.PeekLine()
	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w > 3 || pos >= len(line) || line[pos] != '<' {
		return false
	}
	pos++
	if pos >= len(line) || line[pos] == '\n' {
		reader.Advance(pos)
		return true
	}
	reader.Advance(pos)
	if line[pos] == ' ' || line[pos] == '\t' {
		padding := 0
		if line[pos] == '\t' {
			padding = util.TabWidth(reader.LineOffset()) - 1
		}
		reader.AdvanceAndSetPadding(1, padding)
	}
	return true
}

func (a *asideParser) Trigger() []byte {
	return []byte{'<'}
}

func (a *asideParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	if a.process(reader) {
		return NewAside(), parser.HasChildren
	}
	return nil, parser.NoChildren
}

func (a *asideParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	if a.process(reader) {
		return parser.Continue | parser.HasChildren
	}
	return parser.Close
}

func (a *asideParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

func (a *asideParser) CanInterruptParagraph() bool {
	return true
}

func (a *asideParser) CanAcceptIndentedLine() bool {
	return false
}
