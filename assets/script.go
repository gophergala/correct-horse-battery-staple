// +build js

package main

import (
	"github.com/gophergala/correct-horse-battery-staple/mapview"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	mapView := mapview.New("map")
	marker := mapView.AddMarker(0, 0)

	mapView.OnLocFound(func(loc js.Object) {
		latlng := loc.Get("latlng")
		marker.SetLatLng(latlng.Get("lat").Float(), latlng.Get("lng").Float())
	})

	mapView.StartLocate()
}
