package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"strconv"
)

const MaxRooms = 2
const PlayersPerRoom = 3

var startPort = 31415

func main() {
	lounge := NewSpace("tcp://localhost:31414/lounge")

	rooms := make([]string, 0, MaxRooms)
	awaiting := make([]string, 0, PlayersPerRoom)

	for len(rooms) < cap(rooms) {
		var id string
		lounge.Get("request", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		fmt.Printf("Awaiting %d more players \n", cap(awaiting) - len(awaiting))

		if len(awaiting) == cap(awaiting) {
			uri := "tcp://localhost:" + strconv.Itoa(startPort) + "/game"

			rooms = append(rooms, uri)

			for _, id := range awaiting {
				lounge.Put("join", id, uri)
			}

			game := server.NewGame(uri)
			game.AddPlayers(awaiting)

			go game.Start()

			startPort++
			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}

	lounge.Get("close")
}
