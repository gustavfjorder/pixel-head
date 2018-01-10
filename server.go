package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"strconv"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/model"
	"encoding/gob"
)

const MaxRooms = 10
const PlayersPerRoom = 1

var startPort = 31415

func main() {
	lounge := NewSpace(config.Conf.LoungeUri)
	active := make([]chan bool, 0)
	rooms := make([]string, 0, MaxRooms)
	awaiting := make([]string, 0, PlayersPerRoom)

	for len(rooms) < cap(rooms) {
		var id string
		lounge.Get("request", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		fmt.Printf("Awaiting %d more players \n", cap(awaiting)-len(awaiting))

		if len(awaiting) == cap(awaiting) {
			clientSpaces := make([]server.ClientSpace, len(awaiting))
			for i, id := range awaiting {
				clientSpaces[i].Id = id
				clientSpaces[i].Uri = "tcp://localhost:" + strconv.Itoa(startPort) + "/game/" + id
				clientSpaces[i].Space = SetupSpace(clientSpaces[i].Uri)

				lounge.Put("join", id, clientSpaces[i].Uri)
				startPort++
			}

			game := model.NewGame(awaiting, "Test1")
			active := append(active, make(chan bool, 1))
			go server.Start(&game, clientSpaces, active[len(active)-1])

			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}

	lounge.Get("close")
}

func SetupSpace(uri string) Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register([]model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Zombie{})
	gob.Register([]model.Zombie{})
	gob.Register(model.Shoot{})
	gob.Register([]model.Shoot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.State{})

	space := NewSpace(uri)

	// todo: pSpaces seems to need this to be able to Get/Query on clients
	space.QueryP(&model.Request{})
	space.QueryP(&[]model.Request{})
	space.QueryP(&model.Player{})
	space.QueryP(&[]model.Player{})
	space.QueryP(&model.Zombie{})
	space.QueryP(&[]model.Zombie{})
	space.QueryP(&model.Shoot{})
	space.QueryP(&[]model.Shoot{})
	space.QueryP(&model.Map{})
	space.QueryP(&model.Wall{})
	space.QueryP(&model.Segment{})
	space.QueryP(&model.Point{})
	space.QueryP(&model.State{})

	return space
}
