# BillyGoBlog

BillyGoBlog is a simple blog-in-binary project written in Go and named in honor of [Bill Gosper](https://nl.wikipedia.org/wiki/Bill_Gosper). It uses the `goldmark-qjs-katex` project for server-side math rendering using KaTeX.


## Getting started

In the root of this repository, run `go run .`, wait until you see the server log `Serving on port 8080`, and go to `localhost:8080/posts/example` to see the example blog post.

Or, to actually build the binary, run `go build .`. Then, you can copy the binary to anywhere and run it.


## Architecture

BillyGoBlog consists of the following packages:
- `page` is a library, kind of, that allows you to define HTML pages with Go functions. These functions mostly mirror HTML structure, but also support, for example, math.
- `fileserver` is a simple file server which embeds files. You should add any files that your pages need here.
- `posts` is the package that contains the posts.

And of course there is the `main` package which mostly just starts the file server and serves the pages.


## TODO
  - Support for HTTPS
  - Think more about files, they seem a bit clunky now.
  - Should `RenderablePages` somehow know about where they are served from? It matters for relative paths.
  - Clean up, check if there are basic features missing.
  - There is a lot of code duplication in the `page` package; it's mostly boilerplate that is copy-pasted between every file. This can use a good refactoring.
  - Pre-render static pages. Especially the math seems a bit expensive to render server-side on every page load.
  - I need to go over my use of pointers, I might be copying too often
