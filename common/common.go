package common

type ClientUpdate struct {
	Lat, Lng float64
}

type ClientState struct {
	Name     string
	Lat, Lng float64
}

type ServerUpdate struct {
	Clients []ClientState
}
