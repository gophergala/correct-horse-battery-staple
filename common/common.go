package common

type ClientState struct {
	Id       int64
	Name     string
	Lat, Lng float64
}

type ServerUpdate struct {
	Clients []ClientState
}
