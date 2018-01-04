package server

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/server/model"
)

func StartGame(players []model.Player) {
	spc := NewSpace("tcp://localhost:31415/game1")

	// Register Vec for encoding to space
	gob.Register(pixel.Vec{})

	// Save players into space
	for i, player := range players {
		players[i].Pos = player.Pos.Add(pixel.V(1, 1))
		spc.Put("player", player.Id, player.Pos)
	}

	t := time.Tick(time.Second / 60)

	// Game loop
	for {
		var playerId string
		var playerPos pixel.Vec

		// Load players to handle
		loopPlayers, _ := spc.GetAll("player", &playerId, &playerPos)

		for _, tPlayer := range loopPlayers {
			player := model.Player{
				Id:    tPlayer.GetFieldAt(1).(string),
				Pos:   tPlayer.GetFieldAt(2).(pixel.Vec),
				Stats: nil,
			}

			handleRequest(spc, player)

			spc.Put("player", player.Id, player.Pos)
		}

		<- t
	}
}
