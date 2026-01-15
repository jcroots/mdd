package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	http.HandleFunc("/", serveMarkdown)
	addr := "127.0.0.1:1980"
	log.Printf("Starting server http://%s/\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func serveMarkdown(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/README.md"
	}

	path = strings.TrimPrefix(path, "/")

	absPath, err := filepath.Abs(path)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if !strings.HasPrefix(absPath, cwd) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	if !strings.HasSuffix(path, ".md") {
		http.Error(w, "Only .md files allowed", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := markdown.Render(doc, renderer)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s</title>
<style>
body { max-width: 800px; margin: 40px auto; padding: 0 20px; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; line-height: 1.6; }
code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
pre { background: #f4f4f4; padding: 12px; border-radius: 5px; overflow-x: auto; }
pre code { background: none; padding: 0; }
</style>
</head>
<body>
%s
</body>
</html>`, path, htmlContent)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
	log.Printf("%s\n", absPath)
}
