package fileserver

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var embedded embed.FS

var fileMap map[string]string = map[string]string{
	"/posts/katex.min.css": "assets/stylesheets/katex.min.css",
	"/posts/style.css":     "assets/stylesheets/style.css",
}

var directoryMap map[string]string = map[string]string{
	"/posts/katex_fonts/": "assets/katex_fonts",
}

func Init() {
	for k, v := range fileMap {
		serveFile(k, v)
	}

	for k, v := range directoryMap {
		serveDirectory(k, v)
	}
}

func serveFile(uri string, file string) {
	http.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, embedded, file)
	})
}

func serveDirectory(uri string, directory string) {
	fsys, err := fs.Sub(embedded, directory)
	if err != nil {
		panic(err)
	}
	handler := http.StripPrefix(uri, http.FileServer(http.FS(fsys)))
	http.Handle(uri, handler)
}
