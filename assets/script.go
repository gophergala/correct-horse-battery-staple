// +build js

package main

import (
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
	ws, err := websocket.Dial(js.Global.Get("WebSocketAddress").String())
	if err != nil {
		return err
	}
	enc := json.NewEncoder(ws)
	dec := json.NewDecoder(ws)

	mapView := mapview.New("map")
	var markers []*mapview.Marker
	var bounds *mapview.LatLngBounds

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		clientState := common.ClientState{
			Name: "Me: (" + latlng.Get("lat").String() + ", " + latlng.Get("lng").String() + ")",
			Lat:  latlng.Get("lat").Float(),
			Lng:  latlng.Get("lng").Float(),
		}
		err = enc.Encode(clientState)
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

		for _, marker := range markers {
			marker.RemoveFromMap(mapView)
		}

		markers = nil

		for i, clientState := range msg.Clients {
			markers = append(markers,
				mapView.AddMarkerWithMessage(clientState.Lat, clientState.Lng, clientState.Name))
			if i == 0 {
				bounds = mapview.NewLatLngBounds(
					mapview.NewLatLng(clientState.Lat, clientState.Lng),
					mapview.NewLatLng(clientState.Lat, clientState.Lng))
			} else {
				bounds.Extend(mapview.NewLatLng(clientState.Lat, clientState.Lng))
			}
		}

		if bounds != nil {
			log.Printf("Fit bounds called! %#v %#v , %#v %#v\n", bounds.Call("getNorth").Float(), bounds.Call("getEast").Float(), bounds.Call("getSouth").Float(), bounds.Call("getWest").Float())
			mapView.FitBounds(bounds)
		}

		log.Printf("%#v\n", msg)
	}

	return nil
}

func main() {
	document.AddEventListener("DOMContentLoaded", false, func(_ dom.Event) {
		go func() {
			err := run()
			if err != nil {
				log.Println(err)
			}
		}()
	})
}
