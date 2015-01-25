package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

func accuracyColor(radius float64) string {
	switch true {
	case radius < 10:
		return "#0055ff"
	case radius < 50:
		return "#00ff55"
	case radius < 100:
		return "#ff5500"
	}
	return "#ff0000"
}

type Avatar struct {
	js.Object
	latlng *LatLng
	circle *Circle
	marker *Marker
	popup  *Popup
}

func NewAvatar(lat, lng float64) *Avatar {
	group := L.Call("layerGroup")
	latlng := NewLatLng(lat, lng)
	circle := NewCircle(latlng)
	marker := NewMarker(latlng)
	popup := NewPopup(latlng)
	group.Call("addLayer", circle)
	group.Call("addLayer", marker)

	return &Avatar{
		Object: group,
		latlng: latlng,
		circle: circle,
		marker: marker,
		popup:  popup,
	}
}

func (a *Avatar) Update(lat, lng, accuracy float64, message string) {
	a.latlng.Set(lat, lng)

	a.marker.SetLatLng(a.latlng)

	a.circle.SetRadius(accuracy)
	a.circle.SetLatLng(a.latlng)

	color := accuracyColor(accuracy)
	a.circle.SetStyle(js.M{
		"color":       color,
		"opacity":     0.5,
		"fillColor":   color,
		"fillOpacity": 0.25,
	})

	a.popup.SetContent("<span class=\"popup\">" + message + "</span>")
	a.popup.SetLatLng(a.latlng)

	if message == "" {
		a.Call("removeLayer", a.popup)
	} else {
		a.Call("addLayer", a.popup)
	}
}
