package common

type ClientState struct {
	Name     string
	Lat, Lng float64
}

type ServerUpdate struct {
	Clients []ClientState
}
