package extension

import (
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

// AsideRenderer is a renderer.NodeRenderer implementation that
// renders aside nodes.
type AsideRenderer struct {
	html.Config
}

// NewAsideRenderer returns a new AsideRenderer.
func NewAsideRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &AsideRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *AsideRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindAside, r.renderAside)
}
