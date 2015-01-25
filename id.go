package main

import (
	"sync"

	"github.com/gophergala/correct-horse-battery-staple/urlgen"
)

var id int64
var idLock sync.Mutex

func getUniqueId() int64 {
	idLock.Lock()
	id++
	value := id
	idLock.Unlock()
	return value
}

func generateRoomId() string {
	return urlgen.GetTokenFromId(getUniqueId())
}
