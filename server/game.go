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
	room.Put("rw_lock", 0)

	fmt.Println("Starting game on uri '" + uri + "'")
	fmt.Println("Players in game:", playerIds)

	addPlayerToRoom(room, playerIds)

	t := time.Tick(time.Second / 60)

	// Game loop
	for {
		room.Get("rw_lock", 0)

		players, newShoots := handleRequests(room)

		shoots := append(loadShoots(room), newShoots...)
		handleZombies(room, players, shoots)

		for _, player := range players {
			room.Put(player)
		}

		room.Put("rw_lock", 0)

		<-t
	}
}

func setupSpace(uri string) Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register(model.Zombie{})
	gob.Register(model.Shoot{})

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

func handleRequests(space Space) ([]model.Player, []model.Shoot) {
	// Load incoming requests
	rTuples, _ := space.GetAll(&model.Request{})

	players := make([]model.Player, 0)
	newShoots := make([]model.Shoot, 0)

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
			playerShoots := player.Weapon.GenerateShoots(request.Timestamp, player.Pos)
			newShoots = append(newShoots, playerShoots...)
		}

		players = append(players, player)
	}

	return players, newShoots
}

func handleZombies(room Space, players []model.Player, shoots []model.Shoot) {
	zTuples, _ := room.GetAll(&model.Zombie{})

	for _, zTuple := range zTuples {
		zombie := zTuple.GetFieldAt(0).(model.Zombie)

		// Any shoots hitting the zombie
		for i, shoot := range shoots {
			if shoot.GetPos() == zombie.Pos {
				zombie.Stats.Health -= shoot.Weapon.Power
				shoots = append(shoots[:i], shoots[i+1:]...)
			}
		}

		if zombie.Stats.Health <= 0 {
			continue
		}

		zombie.Move(players)
		zombie.Attack(players)

		room.Put(zombie)
	}

	room.Put(shoots)
}

func loadShoots(room Space) []model.Shoot {
	sTuples, _ := room.GetAll(&model.Shoot{})

	shoots := make([]model.Shoot, len(sTuples))
	for i, sTuple := range sTuples {
		shoots[i] = sTuple.GetFieldAt(0).(model.Shoot)
	}

	return shoots
}
