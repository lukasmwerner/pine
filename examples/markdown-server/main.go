package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/lukasmwerner/pine"
	"github.com/lukasmwerner/pine/middlewares"
	"github.com/russross/blackfriday/v2"
)

//go:embed pages
var pagesFS embed.FS

//go:embed static
var staticFS embed.FS

type Post struct {
	Title   string
	Content string
}

func main() {

	posts := make(map[string]Post)
	results := find(pagesFS, ".", "md")
	fmt.Println(results)
	for _, result := range results {

		b, err := fs.ReadFile(pagesFS, result)
		if err != nil {
			panic(err)
		}
		parsed := blackfriday.Run(b, blackfriday.WithExtensions(blackfriday.CommonExtensions))

		url := "/" + strings.Replace(strings.ReplaceAll(result, ".md", ""), "pages/", "", 1)
		posts[url] = Post{
			Title:   url,
			Content: string(parsed),
		}
	}

	p := pine.New()
	fs := http.FileServer(http.FS(staticFS))
	p.Handle("/static", fs)
	p.Use(middlewares.HTTPLogger())
	p.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(r.URL.Path, "/") {
			path += "index"
		}
		log.Println(path)
		p, ok := posts[path]
		if !ok {
			http.Error(w, errors.New("post not found").Error(), 404)
			return
		}

		fmt.Fprint(w, p.Content)
	})

	srv := &http.Server{
		Handler: p,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on 0.0.0.0:8000")
	log.Fatal(srv.ListenAndServe())
}

// thank you: https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
func find(fsys fs.FS, root, ext string) []string {
	var a []string
	fs.WalkDir(fsys, root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if len(strings.Split(d.Name(), ".")) > 1 && strings.Split(d.Name(), ".")[1] == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}
