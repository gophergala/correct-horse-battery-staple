// +build js

package main

import (
	"fmt"
	"log"

	"github.com/gophergala/correct-horse-battery-staple/common"
	"github.com/gophergala/correct-horse-battery-staple/js/encoding/json"
	"github.com/gophergala/correct-horse-battery-staple/mapview"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var document = dom.GetWindow().Document().(dom.HTMLDocument)

func run() error {
	ws, err := websocket.Dial("ws://" + js.Global.Get("WebSocketHost").String() + "/websocket")
	if err != nil {
		return err
	}
	enc := json.NewEncoder(ws)
	dec := json.NewDecoder(ws)

	mapView := mapview.New("map")
	var markers [4]*mapview.Marker
	markers[0] = mapView.AddMarker(0, 0)
	markers[1] = mapView.AddMarker(0, 0)
	markers[2] = mapView.AddMarker(0, 0)
	markers[3] = mapView.AddMarker(0, 0)

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		clientUpdate := common.ClientUpdate{
			Lat: latlng.Get("lat").Float(),
			Lng: latlng.Get("lng").Float(),
		}
		err = enc.Encode(clientUpdate)
		if err != nil {
			log.Println("enc.Encode:", err)
		}
	})

	mapView.StartLocate()

	for {
		var msg common.ServerUpdate
		err = dec.Decode(&msg)
		if err != nil {
			return err
		}

		for i, clientState := range msg.Clients {
			markers[i].SetLatLng(clientState.Lat, clientState.Lng, clientState.Name)
		}

		document.GetElementByID("content").SetTextContent(document.GetElementByID("content").TextContent() + fmt.Sprintf("%#v\n", msg))
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
