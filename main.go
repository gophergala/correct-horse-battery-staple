package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/shurcooL/go/gopherjs_http"
	"golang.org/x/net/websocket"
)

var httpFlag = flag.String("http", "localhost:8080", "Listen for HTTP connections on this address.")
var webSocketHostFlag = flag.String("websockethost", "localhost:8080", "Listen for WebSocket connections on this address.")
var googleAnalyticsFlag = flag.String("ga", "", "Report to Google Analytics under this code")

var t *template.Template

func loadTemplates() error {
	var err error
	t = template.New("").Funcs(template.FuncMap{})
	t, err = t.ParseGlob("./assets/*.tmpl")
	return err
}

var state = struct {
	mu    sync.Mutex
	rooms map[string]*room
}{rooms: make(map[string]*room)}

type room struct {
	mu          sync.Mutex
	connections map[*websocket.Conn]serverClientState
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	if err := loadTemplates(); err != nil {
		log.Println("loadTemplates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roomId := req.URL.Path[1:]

	var pageVars = struct {
		WebSocketAddress    string
		GoogleAnalyticsCode string
	}{
		WebSocketAddress:    "ws://" + *webSocketHostFlag + "/websocket/" + roomId,
		GoogleAnalyticsCode: *googleAnalyticsFlag,
	}

	var roomExists bool
	state.mu.Lock()
	_, roomExists = state.rooms[roomId]
	state.mu.Unlock()
	if !roomExists {
		// Create room and redirect to it.
		if validateRoomId(roomId) != nil {
			roomId = generateRoomId()
		}
		state.mu.Lock()
		if _, ok := state.rooms[roomId]; !ok {
			state.rooms[roomId] = &room{connections: make(map[*websocket.Conn]serverClientState)}
		}
		state.mu.Unlock()
		go broadcastUpdates(roomId)
		http.Redirect(w, req, "/"+roomId, http.StatusFound)
		return
	}

	err := t.ExecuteTemplate(w, "index.html.tmpl", pageVars)
	if err != nil {
		log.Println("t.Execute:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type serverClientState struct {
	common.ClientState

	validPosition bool
}

func webSocketHandler(ws *websocket.Conn) {
	roomId := ws.Request().URL.Path[len("/websocket/"):]

	state.mu.Lock()
	room, roomExists := state.rooms[roomId]
	state.mu.Unlock()
	if !roomExists {
		return
	}

	room.mu.Lock()
	room.connections[ws] = serverClientState{
		ClientState: common.ClientState{
			Id: getUniqueId(),
		},
		validPosition: false, // When a client first connects, its initial position is not valid.
	}
	room.mu.Unlock()

	dec := json.NewDecoder(ws)

	for {
		var msg common.ClientState
		err := dec.Decode(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		room.mu.Lock()
		clientState := room.connections[ws]
		clientState.validPosition = true
		clientState.Name = msg.Name
		clientState.Lat = msg.Lat
		clientState.Lng = msg.Lng
		room.connections[ws] = clientState
		room.mu.Unlock()
	}

	room.mu.Lock()
	delete(room.connections, ws)
	room.mu.Unlock()
}

func broadcastUpdates(roomId string) {
	state.mu.Lock()
	room := state.rooms[roomId]
	state.mu.Unlock()

	for {
		time.Sleep(1 * time.Second)

		var msg common.ServerUpdate
		var clients []*websocket.Conn // All clients to send an update message to.

		room.mu.Lock()
		for wc, clientState := range room.connections {
			// Only include clients with valid positions in the server update.
			if clientState.validPosition {
				msg.Clients = append(msg.Clients, clientState.ClientState)
			}

			clients = append(clients, wc)
		}
		room.mu.Unlock()

		// If there are no connected clients, break out and remove room.
		if len(clients) == 0 {
			// TODO: Not yet, what if this happens before first client connects?
			//break
		}

		// Don't send empty update messages.
		if len(msg.Clients) == 0 {
			continue
		}

		sort.Sort(msg.Clients)

		for _, ws := range clients {
			err := json.NewEncoder(ws).Encode(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Remove room.
	state.mu.Lock()
	delete(state.rooms, roomId)
	state.mu.Unlock()
}

func main() {
	flag.Parse()

	err := loadTemplates()
	if err != nil {
		log.Fatalln("loadTemplates:", err)
	}

	http.Handle("/favicon.ico/", http.NotFoundHandler())
	http.Handle("/", http.HandlerFunc(mainHandler))
	http.Handle("/websocket/", websocket.Handler(webSocketHandler))
	http.Handle("/assets/", http.FileServer(http.Dir("./")))
	http.Handle("/assets/script.go.js", gopherjs_http.GoFiles("./assets/script.go"))

	err = http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
