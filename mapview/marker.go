package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Marker struct {
	js.Object
}

func NewMarker(lat, lng float64) *Marker {
	return &Marker{
		Object: L.Call("marker", NewLatLng(lat, lng)),
	}
}

func (m *Marker) SetLatLng(lat, lng float64) {
	m.Call("setLatLng", NewLatLng(lat, lng)).Call("bindPopup", "Here I Am").Call("openPopup")
}
