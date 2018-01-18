package server

import (
	"fmt"
	"github.com/gustavfjorder/pixel-head/config"
	"strconv"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/gustavfjorder/pixel-head/setup"
	. "github.com/pspaces/gospace/space"
)

const PlayersPerRoom = 2

func NewLounge(maxRooms int) {
	setup.RegisterModels()

	startPort := 31415

	ip, _ := config.GetIp()
	fmt.Println(ip)

	lounge := NewSpace("tcp://" + ip + ":" + strconv.Itoa(startPort) + "/lounge")

	active := make([]chan bool, 0)
	rooms := make([]string, 0, maxRooms)
	awaiting := make([]string, 0, PlayersPerRoom)

	for len(rooms) < cap(rooms) {
		var id string
		lounge.Get("request", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		fmt.Printf("Awaiting %d more players \n", cap(awaiting)-len(awaiting))

		if len(awaiting) == cap(awaiting) {
			clientSpaces := make([]ClientSpace, len(awaiting))
			for i, id := range awaiting {
				startPort++
				clientSpaces[i].Id = id
				clientSpaces[i].Uri = "tcp://" + ip + ":" + strconv.Itoa(startPort) + "/game/" + id
				clientSpaces[i].Space = SetupSpace(clientSpaces[i].Uri)

				lounge.Put("join", id, clientSpaces[i].Uri)
			}

			game := model.NewGame(awaiting, "Test1")
			active := append(active, make(chan bool, 1))
			go Start(&game, clientSpaces, active[len(active)-1])

			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}

	lounge.Get("close")
}
