package server

import (
	"github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server/model"
)

//type Request struct {
//	Id string
//	CurrentWep byte
//	Dir float64
//	Move bool
//	Shoot bool
//}

func handleRequest(space gospace.Space, player model.Player) {
	var requestType string

	t, e := space.GetP(player.Id, "request", &requestType)
	if e != nil {
		return
	}

	requestType = t.GetFieldAt(2).(string)

	switch requestType {
	case "move":
		var dir float64
		request, err := space.GetP(player.Id, "request", &dir)
		if err == nil {
			dir = request.GetFieldAt(2).(float64)

			player.Move(dir)
		}

	default:
		// ...
	}
}