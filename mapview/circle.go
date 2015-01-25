package mapview

import (
	"github.com/gopherjs/gopherjs/js"
)

type Circle struct {
	js.Object
}

func NewCircle(latlng *LatLng) *Circle {
	return &Circle{
		Object: L.Call("circle", latlng),
	}
}

func (c *Circle) LatLng() *LatLng {
	return &LatLng{
		Object: c.Call("getLatLng"),
	}
}

func (c *Circle) SetLatLng(latlng *LatLng) {
	c.Call("setLatLng", latlng)
}

func (c *Circle) Radius() float64 {
	return c.Call("getRadius").Float()
}

func (c *Circle) SetRadius(radius float64) {
	c.Call("setRadius", radius)
}

func (c *Circle) SetStyle(style js.M) {
	c.Call("setStyle", style)
}
