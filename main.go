package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
  mux := http.NewServeMux()

  // Standard root handler
  mux.Handle("/", http.FileServer(http.Dir("./static")))

  // MD to HTML Blog Handler
  mux.HandleFunc("/blog/entry/", mdToHtml)
  
  log.Println("Starting server on https://localhost:8080")
  err := http.ListenAndServeTLS(":8080", "./certs/server.crt", "./certs/server.key", mux)
  if err != nil {
    log.Fatal(err)
  }
}

func mdToHtml(w http.ResponseWriter, r *http.Request) {
  entryId := strings.TrimPrefix(r.URL.Path, "/blog/entry/")
  blogPath := fmt.Sprintf("./static/blog/entries/%v", entryId)
  // TODO: Check blogPath ends with .md file extension
  if _, err := os.Stat(blogPath); err != nil {
    if errors.Is(err, os.ErrNotExist) {
      fmt.Fprint(w, "Blog not found")
      return
    } else {
      log.Print(err)
      fmt.Fprint(w, "Error occured while reading blog")
      return
    }
  }
  md, err := os.ReadFile(blogPath)
  if err != nil {
      log.Print(err)
      fmt.Fprint(w, "Error occured while reading blog")
      return
  }
  // Combined using Bitwise OR
  extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
  p := parser.NewWithExtensions(extensions)
  doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
  opts := html.RendererOptions{Flags: htmlFlags}
  renderer := html.NewRenderer(opts)

  renderedBlog := markdown.Render(doc, renderer)

  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, string(renderedBlog))
}
