package main

import (
	"sync"

	"github.com/gophergala/correct-horse-battery-staple/urlgen"
)

var id int64
var idLock sync.Mutex

func generateRoomId() string {
	idLock.Lock()
	id++
	token := urlgen.GetTokenFromId(id)
	idLock.Unlock()

	return token
}
