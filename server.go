package main

import (
	"fmt"
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"github.com/gustavfjorder/pixel-head/model"
)

const N = 1

func main() {
	lounge := NewSpace("tcp://localhost:31414/lounge")

	awaiting := make([]string, 0, N)

	for {
		var id string
		lounge.Get("client", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		if len(awaiting) == cap(awaiting) {
			players := make([]model.Player,N)
			for i, id := range awaiting {
				players[i] = model.Player{
					Id:id,
					Pos: pixel.V(0,0),
					Weapon: model.Weapons[model.Handgun],
					Stats: model.NewStats(model.Human),
				}
				lounge.Put(id, "ready")
			}

			go server.StartGame(players)
		}
	}
}
