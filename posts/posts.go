package posts

import (
	_ "embed"

	"github.com/rubenvannieuwpoort/billygoblog/page"
)

//go:embed post.template
var postHTMLTemplate string
var postTemplate, _ = page.Template("templates", postHTMLTemplate)

var Posts []page.RenderablePage = []page.RenderablePage{
	example,
}
