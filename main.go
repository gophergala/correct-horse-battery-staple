package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/shurcooL/go/gopherjs_http"
)

var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

var t *template.Template

func loadTemplates() error {
	var err error
	t = template.New("").Funcs(template.FuncMap{})
	t, err = t.ParseGlob("./assets/*.tmpl")
	return err
}

var state struct {
	mu sync.Mutex
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	if err := loadTemplates(); err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state.mu.Lock()
	err := t.ExecuteTemplate(w, "index.html.tmpl", &state)
	state.mu.Unlock()
	if err != nil {
		log.Println("t.Execute:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()

	err := loadTemplates()
	if err != nil {
		log.Fatalln("loadTemplates:", err)
	}

	http.Handle("/favicon.ico/", http.NotFoundHandler())
	http.HandleFunc("/", mainHandler)
	http.Handle("/assets/", http.FileServer(http.Dir("./")))
	http.Handle("/assets/script.go.js", gopherjs_http.GoFiles("./assets/script.go"))

	err = http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
