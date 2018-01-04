package server

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
)

func StartGame(players []Player) {
	spc := NewSpace("tcp://localhost:31415/game1")

	// Save players into space
	for i, player := range players {
		players[i].Pos = player.Pos.Add(pixel.V(1, 1))
		spc.Put(player.Id, player.Pos.XY)
	}

	// Game loop
	for {
		var playerId string
		var playerPos pixel.Vec

		// Load players to handle
		loopPlayers, _ := spc.GetAll(&playerId, &playerPos)

		for _, tPlayer := range loopPlayers {
			player := Player{tPlayer.GetFieldAt(0).(string), tPlayer.GetFieldAt(2).(pixel.Vec)}

			spc.Put(player.Id, player.Pos.XY)
		}

		time.Sleep(time.Second / 60)
	}
}
