// +build js

package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	tilesUrl = "http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
)

var L = js.Global.Get("L")

type Marker struct {
	js.Object
}

func NewMarker(lat, lng float64) *Marker {
	return &Marker{
		Object: L.Call("marker", L.Call("latLng", lat, lng)),
	}
}

func (m *Marker) SetLatLng(lat, lng float64) {
	m.Call("setLatLng", L.Call("latLng", lat, lng))
}

type MapView struct {
	js.Object
}

func NewMapView(id string) *MapView {
	mapView := L.Call("map", id)
	L.Call("tileLayer", tilesUrl).Call("addTo", mapView)

	return &MapView{
		Object: mapView,
	}
}

func (mv *MapView) SetView(lat, lng float64, zoom int) {
	mv.Call("setView", L.Call("latLng", lat, lng), zoom)
}

func (mv *MapView) AddMarker(lat, lng float64) *Marker {
	marker := NewMarker(lat, lng)
	marker.Call("addTo", mv)
	return marker
}

func (mv *MapView) StartLocate() {
	mv.Call("locate", js.M{
		"setView":            true,
		"watch":              true,
		"enableHighAccuracy": true,
	})
}

func (mv *MapView) StopLocate() {
	mv.Call("stopLocate")
}

func (mv *MapView) OnLocFound(cb func(js.Object)) {
	mv.Call("on", "locationfound", cb)
}

func main() {
	mapView := NewMapView("map")
	marker := mapView.AddMarker(0, 0)

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		marker.SetLatLng(latlng.Get("lat").Float(), latlng.Get("lng").Float())
	})

	mapView.StartLocate()
}
