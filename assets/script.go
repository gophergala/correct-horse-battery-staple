// +build js

package main

import (
	"fmt"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/gophergala/correct-horse-battery-staple/js/encoding/json"
	"github.com/gophergala/correct-horse-battery-staple/mapview"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

func setupMapView() error {
	mapView := mapview.New("map")
	marker := mapView.AddMarker(0, 0)

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		marker.SetLatLng(latlng.Get("lat").Float(), latlng.Get("lng").Float())
	})

	mapView.StartLocate()

	return nil
}

func run() error {
	err := setupMapView()
	if err != nil {
		return err
	}

	ws, err := websocket.Dial("ws://" + js.Global.Get("WebSocketHost").String() + "/websocket")
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
