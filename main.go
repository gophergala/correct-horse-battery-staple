package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/shurcooL/go/gopherjs_http"
	"golang.org/x/net/websocket"
)

var httpFlag = flag.String("http", "localhost:8080", "Listen for HTTP connections on this address.")
var webSocketHostFlag = flag.String("websockethost", "localhost:8080", "Listen for WebSocket connections on this address.")

var t *template.Template

func loadTemplates() error {
	var err error
	t = template.New("").Funcs(template.FuncMap{})
	t, err = t.ParseGlob("./assets/*.tmpl")
	return err
}

var state struct {
	mu            sync.Mutex
	WebSocketHost string
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	if err := loadTemplates(); err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state.mu.Lock()
	state.WebSocketHost = *webSocketHostFlag
	err := t.ExecuteTemplate(w, "index.html.tmpl", &state)
	state.mu.Unlock()
	if err != nil {
		log.Println("t.Execute:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func websocketHandler(ws *websocket.Conn) {
	time.Sleep(5 * time.Second)

	var msg = common.ServerUpdate{
		Lat:     37.7740,
		Lng:     -122.4175,
		Message: "Starting Backend Pos",
	}

	err := json.NewEncoder(ws).Encode(msg)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(5 * time.Second)

	// Send another update.
	err = json.NewEncoder(ws).Encode(common.ServerUpdate{
		Lat:     37.7740,
		Lng:     -122.4165,
		Message: "New Backend Pos!",
	})
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
	http.Handle("/assets/script.go.js", gopherjs_http.GoFiles("./assets/script.go"))

	err = http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
