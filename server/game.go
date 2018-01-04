package server

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/server/model"
	"fmt"
)

func StartGame(playerIds []string) {
	fmt.Println("Starting game with:", playerIds)

	spc := NewSpace("tcp://localhost:31415/game1")

	// Register models for encoding to space
	gob.Register(Request{})
	gob.Register(model.Player{})

	// Save players into space
	for _, id := range playerIds {
		player := model.Player{
			Id:     id,
			Pos:    pixel.V(0, 0),
			Weapon: model.Weapon{ // todo: Weapons should be defined somewhere
				Id:    "gun",
				Name:  "Magnum",
				Power: 1,
				Range: 8,
			},
		}

		spc.Put(player)
	}

	t := time.Tick(time.Second / 60)

	// Game loop
	for {
		spc.Get("loop_lock")

		// Load incoming requests
		rTuples, _ := spc.GetAll(&Request{})
		for _, rTuple := range rTuples {
			request := rTuple.GetFieldAt(0).(Request)
			fmt.Println("Handling request:", request)

			// Load player who made the request
			t, _ := spc.GetP(model.Player{Id: request.PlayerId})
			player := t.GetFieldAt(0).(model.Player)

			// Change weapon
			player.Weapon = request.CurrentWep

			if request.Move {
				// todo: check if move is doable in map
				player = player.Move(request.Dir)
			}

			if request.Reload {
				// todo: handle reload
			} else if request.Shoot {
				// todo: handle shoot
			}

			spc.Put(player)
		}

		spc.Put("loop_lock")

		<- t
	}
}
