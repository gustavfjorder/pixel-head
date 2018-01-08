package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"strconv"
	"github.com/gustavfjorder/pixel-head/Config"
)

const MaxRooms = 10
const PlayersPerRoom = 2

var startPort = 31415

func main() {
	lounge := NewSpace(config.Conf.LoungeUri)

	rooms := make([]string, 0, MaxRooms)
	awaiting := make([]string, 0, PlayersPerRoom)

	//lounge.Put("request", "a")
	//lounge.Put("request", "b")

	for len(rooms) < cap(rooms) {
		var id string
		lounge.Get("request", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		fmt.Printf("Awaiting %d more players \n", cap(awaiting) - len(awaiting))

		if len(awaiting) == cap(awaiting) {
			gameUri := "tcp://localhost:" + strconv.Itoa(startPort) + "/game"
			startPort++

			rooms = append(rooms, gameUri)

			clientUris := make([]string, len(awaiting))
			for i, id := range awaiting {
				clientUris[i] = "tcp://localhost:" + strconv.Itoa(startPort) + "/game/" + id
				startPort++
				fmt.Println("Client uri: " + clientUris[i])
			}

			game := server.NewGame(gameUri, clientUris)
			game.AddPlayers(awaiting)

			go game.Start()

			for i, id := range awaiting {
				lounge.Put("join", id, gameUri, clientUris[i])
			}

			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}

	lounge.Get("close")
}
