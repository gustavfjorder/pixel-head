package server

import (
	"time"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"strconv"
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/config"
	"encoding/gob"
)

type ClientSpace struct {
	Id  string
	Uri string
	space.Space
}

func Start(g *model.Game, clientSpaces []ClientSpace, finished <-chan bool) {
	fmt.Println("Starting game")
	start := time.Now()

	for _, spc := range clientSpaces {
		spc.Put("map", g.CurrentMap)
	}

	for _, spc := range clientSpaces {
		spc.Get("joined")
	}

	fmt.Println("Starting game loop")
	t := time.Tick(config.Conf.ServerHandleSpeed)
	for g.CurrentLevel < len(model.Levels) {
		fmt.Println("Starting level " + strconv.Itoa(g.CurrentLevel))

		duration := time.Second * 10
		if g.CurrentLevel == 0 {
			duration = 0
		}

		levelRdy := make(chan bool, 1)
		zombiesSpawned := false
		time.AfterFunc(duration, func() {
			g.PrepareLevel(levelRdy)
		})

		fmt.Println("after prepare")

		for {
			select {
			case <-levelRdy:
				zombiesSpawned = true
			default:
			}

			//Update game
			g.State.Timestamp = time.Since(start)
			g.HandleRequests(collectRequests(clientSpaces))
			g.HandleZombies()
			g.HandleShots()
			g.HandlePlayers()

			//Send new game state to clients
			if config.Conf.Online {
				for _, spc := range clientSpaces {
					spc.GetP("state", &model.State{})
					spc.Put("state", g.State)
				}
			}

			//If all players died end game
			if len(g.State.Players) == 0 {
				goto endgame
			}

			//If all zombies have been killed, go to next level
			if zombiesSpawned && len(g.State.Zombies) <= 0 {
				break
			}
			<-t
		}

		g.CurrentLevel++
	}
endgame:
	for _, spc := range clientSpaces {
		spc.Put("game over")
	}
	for _, spc := range clientSpaces {
		spc.Get("quit")
	}
	<-finished
	fmt.Println("Game ended")
}

func collectRequests(clientSpaces []ClientSpace) (requests []model.Request) {
	requests = make([]model.Request, len(clientSpaces))
	for i, spc := range clientSpaces {
		rtuples, _ := spc.GetAll(&model.Request{})
		for _, rtuple := range rtuples {
			request := rtuple.GetFieldAt(0).(model.Request)
			requests[i] = requests[i].Merge(request)
		}
		requests[i].PlayerId = spc.Id
	}
	return requests
}

func SetupSpace(uri string) space.Space {
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

	space := space.NewSpace(uri)

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