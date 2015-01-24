package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

var state = struct {
	mu            sync.Mutex
	WebSocketHost string
	connections   map[*websocket.Conn]common.ClientState
}{connections: make(map[*websocket.Conn]common.ClientState)}

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

func webSocketHandler(ws *websocket.Conn) {
	state.mu.Lock()
	state.connections[ws] = common.ClientState{
		ValidPosition: false, // When a client first connects, its initial position is not valid.
	}
	state.mu.Unlock()

	dec := json.NewDecoder(ws)

	for {
		var msg common.ClientState
		err := dec.Decode(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println("Got an update:", msg)
		state.mu.Lock()
		clientState := state.connections[ws]
		clientState.ValidPosition = true
		clientState.Name = msg.Name
		clientState.Lat = msg.Lat
		clientState.Lng = msg.Lng
		state.connections[ws] = clientState
		state.mu.Unlock()
	}

	state.mu.Lock()
	delete(state.connections, ws)
	state.mu.Unlock()
}

func broadcastUpdates() {
	for {
		time.Sleep(1 * time.Second)

		var msg common.ServerUpdate
		var clients []*websocket.Conn // All clients to send an update message to.

		state.mu.Lock()
		for wc, clientState := range state.connections {
			// Only include clients with valid positions in the server update.
			if clientState.ValidPosition {
				msg.Clients = append(msg.Clients, clientState)
			}

			clients = append(clients, wc)
		}
		state.mu.Unlock()

		// Don't send empty update messages.
		if len(msg.Clients) == 0 {
			continue
		}

		for _, ws := range clients {
			err := json.NewEncoder(ws).Encode(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	flag.Parse()

	err := loadTemplates()
	if err != nil {
		log.Fatalln("loadTemplates:", err)
	}

	http.Handle("/favicon.ico/", http.NotFoundHandler())
	http.Handle("/", http.HandlerFunc(mainHandler))
	http.Handle("/websocket", websocket.Handler(webSocketHandler))
	http.Handle("/assets/", http.FileServer(http.Dir("./")))
	http.Handle("/assets/script.go.js", gopherjs_http.GoFiles("./assets/script.go"))

	go broadcastUpdates()

	err = http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
