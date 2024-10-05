package posts

import "github.com/rubenvannieuwpoort/billygoblog/page"

var example page.RenderablePage = postTemplate.Instantiate(
	page.Header(page.Title("Example")),
	page.Section(
		page.Paragraph("This is an example blogpost. Of course extremely advanced features like ", page.Bold("bold text"), " are available."),
	),
	page.Section(
		page.H2("And there's more"),
		page.Paragraph("You can use ", page.InlineCode("inline code"), ", or code blocks:"),
		page.Code(
			`def foo():
while True:
    print("bar")

print("done!")`),
	),
	page.Section(
		page.H2("Math"),
		page.Paragraph("It's also possible to use inline math: ", page.InlineMath(`\int_0^1 \pi\ \text{d}x = \pi`), ". Or, well, not inline:"),
		page.Math(`\int_0^1 \pi\ \text{d}x = \pi`),
	),
)
