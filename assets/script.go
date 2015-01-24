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

var serverMarker *mapview.Marker

func setupMapView() error {
	mapView := mapview.New("map")
	marker := mapView.AddMarker(0, 0)
	serverMarker = mapView.AddMarker(0, 0)

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		marker.SetLatLng(latlng.Get("lat").Float(), latlng.Get("lng").Float(), "Here I Am")
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
		var msg common.ServerUpdate
		err = dec.Decode(&msg)
		if err != nil {
			return err
		}

		document.GetElementByID("content").SetTextContent(document.GetElementByID("content").TextContent() + fmt.Sprintf("%#v\n", msg))
		serverMarker.SetLatLng(msg.Lat, msg.Lng, msg.Message)
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
