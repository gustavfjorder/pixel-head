package server

import (
	"fmt"
	"github.com/gustavfjorder/pixel-head/config"
	"strconv"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/gustavfjorder/pixel-head/setup"
	. "github.com/pspaces/gospace/space"
	"net"
)

const PlayersPerRoom = 2
var startPort = 31415
var ip, _ = config.GetIp()


func NewLounge(maxRooms int) (spc *Space, port string) {
	config.Conf.Online = true
	setup.RegisterModels()

	fmt.Println(ip)
	port = NextValidPort()
	config.Conf.LoungeUri = "tcp://" + ip + ":" + port + "/lounge"
	lounge := NewSpace(config.Conf.LoungeUri)
	go newLounge(maxRooms, lounge)
	return &lounge, port
}

func newLounge(maxRooms int, lounge Space) {
	active := make([]chan bool, 0)
	rooms := make([]string, 0, maxRooms)
	awaiting := make([]string, 0, PlayersPerRoom)

	for len(rooms) < cap(rooms) {
		_, err := lounge.GetP("close")
		if err == nil{
			fmt.Println("Stopping server")
			return
		}
		var id string
		_, err = lounge.GetP("request", &id)
		if err != nil {
			continue
		}

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		fmt.Printf("Awaiting %d more players \n", cap(awaiting)-len(awaiting))

		if len(awaiting) == cap(awaiting) {
			clientSpaces := make([]ClientSpace, len(awaiting))
			for i, id := range awaiting {
				clientSpaces[i].Id = id
				clientSpaces[i].Uri = "tcp://" + ip + ":" + NextValidPort() + "/game/" + id
				clientSpaces[i].Space = SetupSpace(clientSpaces[i].Uri)

				lounge.Put("join", id, clientSpaces[i].Uri)
			}

			game := model.NewGame(awaiting, "Test1")
			active := append(active, make(chan bool, 1))
			go Start(&game, clientSpaces, active[len(active)-1])

			awaiting = make([]string, 0, PlayersPerRoom)
		}
	}
}

func NextValidPort() string{
	for { // Check that the port is valid before making the space
		conn, err := net.Listen("tcp", ":" + strconv.Itoa(startPort))
		if err == nil {
			conn.Close()
			break
		} else {
			fmt.Println(startPort, "is invalid retrying")
			startPort++
		}
	}
	return strconv.Itoa(startPort)
}
