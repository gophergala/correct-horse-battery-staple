package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Marker struct {
	js.Object
	Message *Popup
	Lat     float64
	Lng     float64
	mapView *MapView
}

func NewMarker(lat, lng float64) *Marker {
	return &Marker{
		Object:  L.Call("marker", NewLatLng(lat, lng)),
		Message: nil,
		Lat:     lat,
		Lng:     lng,
		mapView: nil,
	}
}

func (m *Marker) SetLatLng(lat, lng float64) {
	m.Lat = lat
	m.Lng = lng
	m.Call("setLatLng", NewLatLng(lat, lng))
	if m.Message != nil {
		m.Message.Call("setLatLng", NewLatLng(lat, lng))
	}
}

func (m *Marker) AddToMap(mapView *MapView) {
	m.mapView = mapView
	m.Call("addTo", mapView)
	if m.Message != nil {
		mapView.Call("addLayer", m.Message)
	}
}

func (m *Marker) SetMessage(message string) {
	if message != "" {
		if m.Message == nil {
			m.Message = NewPopup(m.Lat, m.Lng)
		}
		m.Message.SetContent(message)
	} else {
		if m.Message != nil && m.mapView != nil {
			m.mapView.Call("removeLayer", m.Message)
		}
	}
}

func (m *Marker) RemoveFromMap(mapView *MapView) {
	if m.Message != nil {
		mapView.Call("removeLayer", m.Message)
	}
	mapView.Call("removeLayer", m)
}
