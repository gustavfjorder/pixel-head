package main

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
)

type Player struct {
	Id  string
	Pos pixel.Vec
}

func startGame(players []Player) {
	spc := NewSpace("tcp://localhost:31415/game1")
	for {
		for i, player := range players {
			players[i].Pos = player.Pos.Add(pixel.V(1, 1))
			spc.Put(player.Id, player.Pos.XY)
		}
		time.Sleep(time.Second / 60)
	}
}
