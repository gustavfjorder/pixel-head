package server

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
)

func StartGame(uri string, playerIds []string) {
	room := setupSpace(uri)

	fmt.Println("Starting game on uri '" + uri + "'")
	fmt.Println("Players in game:", playerIds)

	addPlayerToRoom(room, playerIds)

	t := time.Tick(time.Second / 60)

	// Game loop
	for {
		room.Get("loop_lock")

		handleRequests(room)

		room.Put("loop_lock")

		<- t
	}
}

func setupSpace(uri string) Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})

	return NewSpace(uri)
}

func addPlayerToRoom(space Space, playerIds []string) {
	// Save players into space
	for _, id := range playerIds {
		player := model.Player{
			Id:     id,
			Pos:    pixel.V(0, 0),
			Weapon: model.Weapons[model.Handgun],
		}

		space.Put(player)
	}
}

func handleRequests(space Space) {
	// Load incoming requests
	rTuples, _ := space.GetAll(&model.Request{})
	for _, rTuple := range rTuples {
		request := rTuple.GetFieldAt(0).(model.Request)
		fmt.Println("Handling request:", request)

		// Load player who made the request
		t, _ := space.GetP(model.Player{Id: request.PlayerId})
		player := t.GetFieldAt(0).(model.Player)

		// Change weapon
		player.Weapon = model.Weapons[request.CurrentWep]

		if request.Move {
			// todo: check if move is doable in map
			player = player.Move(request.Dir)
		}

		if request.Reload {
			player.Weapon.RefillMag()
		} else if request.Shoot {
			shoot := model.Shoot{
				Start:     player.Pos,
				Angle:     player.Pos.Angle(),
				StartTime: request.Timestamp,
				Weapon:    player.Weapon,
			}

			space.Put(shoot)
		}

		space.Put(player)
	}
}
