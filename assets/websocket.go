// +build js

package main

import (
	"fmt"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/gophergala/correct-horse-battery-staple/js/encoding/json"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

func main() {
	document.AddEventListener("DOMContentLoaded", false, func(_ dom.Event) {
		go func() {
			err := run()
			if err != nil {
				fmt.Println(err)
			}
		}()
	})
}

func run() error {
	ws, err := websocket.Dial("ws://localhost:8080/websocket") // TODO: Pass the server websocket endpoint URL from backend.
	if err != nil {
		return err
	}

	dec := json.NewDecoder(ws)

	for {
		var msg common.SampleMessage
		err = dec.Decode(&msg)
		if err != nil {
			return err
		}

		msg.Message += " And frontend!"

		document.GetElementByID("content").SetTextContent(fmt.Sprintf("%#v", msg))
	}
}
