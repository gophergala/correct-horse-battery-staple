package main

import (
	"fmt"
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

// validateRoomId returns an error if id is of unexpected format.
func validateRoomId(id string) error {
	if len(id) < 3 || len(id) > 64 {
		return fmt.Errorf("id length is %v", len(id))
	}

	for _, b := range []byte(id) {
		ok := ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || ('0' <= b && b <= '9') || b == '-' || b == '_'
		if !ok {
			return fmt.Errorf("id contains unexpected character %+q", b)
		}
	}

	return nil
}
