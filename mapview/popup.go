package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Popup struct {
	js.Object
}

func NewPopup(latlng *LatLng) *Popup {
	options := make(map[string]interface{})
	options["offset"] = NewPoint(0, -24)
	options["closeButton"] = false
	popup := &Popup{
		Object: L.Call("popup", options),
	}
	popup.SetLatLng(latlng)
	return popup
}

func (popup *Popup) SetContent(msg string) {
	popup.Call("setContent", msg)
}

func (popup *Popup) SetLatLng(latlng *LatLng) {
	popup.Call("setLatLng", latlng)
}
