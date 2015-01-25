package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Marker struct {
	js.Object
	Message *Popup
	Lat     float64
	Lng     float64
}

func NewMarker(lat, lng float64) *Marker {
	return &Marker{
		Object:  L.Call("marker", NewLatLng(lat, lng)),
		Message: NewPopup(lat, lng),
		Lat:     lat,
		Lng:     lng,
	}
}

func (m *Marker) SetLatLng(lat, lng float64) {
	m.Lat = lat
	m.Lng = lng
	m.Call("setLatLng", NewLatLng(lat, lng))
	m.Message.Call("setLatLng", NewLatLng(lat, lng))
}

func (m *Marker) AddToMap(mapView *MapView) {
	m.Call("addTo", mapView)
	mapView.Call("addLayer", m.Message)
}

func (m *Marker) SetMessage(message string) {
	m.Message.SetContent(message)
}

func (m *Marker) RemoveFromMap(mapView *MapView) {
	mapView.Call("removeLayer", m.Message)
	mapView.Call("removeLayer", m)
}
