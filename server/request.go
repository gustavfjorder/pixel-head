package server

import (
	"github.com/pspaces/gospace"
)

func handleRequest(space gospace.Space, player Player) {
	var requestType string

	t, e := space.GetP(player.Id, "request", &requestType)
	if e != nil {
		return
	}

	requestType = t.GetFieldAt(2).(string)

	switch requestType {
	default:
		// ...
	}
}