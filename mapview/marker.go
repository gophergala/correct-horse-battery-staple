package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Marker struct {
	js.Object
}

func NewMarker(latlng *LatLng) *Marker {
	return &Marker{
		Object: L.Call("marker", latlng),
	}
}

func (m *Marker) SetLatLng(latlng *LatLng) {
	m.Call("setLatLng", latlng)
}
