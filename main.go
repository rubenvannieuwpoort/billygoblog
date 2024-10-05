package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rubenvannieuwpoort/billygoblog/fileserver"
	"github.com/rubenvannieuwpoort/billygoblog/page"
	"github.com/rubenvannieuwpoort/billygoblog/posts"
)

const PORT = 8080

// posts are defined in posts/<post-name>.go and mapped to a template in posts/template.go
// the templates themselves are defined in templates/<template-name>.template

func main() {
	log.Println("Starting server...")

	fileserver.Init()
	serveBySlug("/posts/", posts.Posts)

	log.Printf("Serving on port %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

func serveBySlug(root string, pages []page.RenderablePage) {
	for _, page := range pages {
		slug := makeSlug(page.Title)
		uri := fmt.Sprintf("%s%s", root, slug)
		http.HandleFunc(uri, func(w http.ResponseWriter, req *http.Request) {
			page.Render(w)
		})
	}
}

func makeSlug(title string) string {
	title = strings.ToLower(title)
	return strings.ReplaceAll(title, " ", "-")
}
