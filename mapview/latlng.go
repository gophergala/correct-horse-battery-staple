package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type LatLng struct {
	js.Object
	Lat float64 `js:lat`
	Lng float64 `js:lng`
}

func NewLatLng(lat, lng float64) *LatLng {
	return &LatLng{
		Object: L.Call("latLng", lat, lng),
	}
}

func (ll *LatLng) Set(lat, lng float64) {
	ll.Lat = lat
	ll.Lng = lng
}
