package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	tilesUrl = "http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
)

var L = js.Global.Get("L")

type MapView struct {
	js.Object
}

func New(id string) *MapView {
	mapView := L.Call("map", id)
	L.Call("tileLayer", tilesUrl).Call("addTo", mapView)

	return &MapView{
		Object: mapView,
	}
}

func (mv *MapView) SetView(lat, lng float64, zoom int) {
	mv.Call("setView", NewLatLng(lat, lng), zoom)
}

func (mv *MapView) AddMarker(lat, lng float64) *Marker {
	marker := NewMarker(lat, lng)
	marker.Call("addTo", mv)
	return marker
}

func (mv *MapView) RemoveMarker(marker *Marker) {
	mv.Call("removeLayer", marker)
}

func (mv *MapView) AddMarkerWithMessage(lat, lng float64, msg string) *Marker {
	marker := NewMarker(lat, lng)
	marker.AddToMap(mv)
	marker.SetMessage(msg)
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
