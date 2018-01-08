package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"strconv"
	"github.com/gustavfjorder/pixel-head/Config"
	"github.com/gustavfjorder/pixel-head/model"
)

const MaxRooms = 2
const PlayersPerRoom = 3

var startPort = 31415

func main() {
	lounge := setupSpace(Config.Conf.LoungeUri)

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
				lounge.Put("join", id, uri, model.NewPlayer(model.Human))
			}

			go server.StartGame(uri, awaiting)

			startPort++
			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}

	lounge.Get("close")
}
