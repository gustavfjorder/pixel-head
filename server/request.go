package server

import (
	"github.com/pspaces/gospace"
	"github.com/faiface/pixel"
)

func handleRequest(space gospace.Space, player Player) {
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

			player.Pos = player.Pos.Add(pixel.V(2, 0).Rotated(dir))
		}

	default:
		// ...
	}
}