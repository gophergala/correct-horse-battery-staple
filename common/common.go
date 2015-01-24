package common

type ClientState struct {
	Name     string
	Lat, Lng float64

	ValidPosition bool `json:"-"`
}

type ServerUpdate struct {
	Clients []ClientState
}
