package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/shurcooL/go/gopherjs_http"
	"golang.org/x/net/websocket"
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

func websocketHandler(ws *websocket.Conn) {
	var msg = common.SampleMessage{
		X:       12,
		Y:       34,
		Message: "Hello from backend!",
	}

	err := json.NewEncoder(ws).Encode(msg)
	if err != nil {
		log.Println(err)
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
	http.Handle("/websocket", websocket.Handler(websocketHandler))
	http.Handle("/assets/", http.FileServer(http.Dir("./")))
	http.Handle("/assets/websocket.go.js", gopherjs_http.GoFiles("./assets/websocket.go"))

	err = http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
