package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"strconv"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/gustavfjorder/pixel-head/model"
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
				clientSpaces[i].Space = server.SetupSpace(clientSpaces[i].Uri)

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
