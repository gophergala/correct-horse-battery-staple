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

// TODO: Setting message doesn't really belong here, so need to clean up the API later!
func (m *Marker) SetLatLng(lat, lng float64, message string) {
	m.Call("setLatLng", NewLatLng(lat, lng)).Call("bindPopup", message).Call("openPopup")
}
