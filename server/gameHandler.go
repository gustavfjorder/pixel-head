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

		duration := time.Second * 1
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
			g.HandleRequests(collectRequests(clientSpaces, g.PlayerIds))
			g.HandleBarrels()
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

func collectRequests(clientSpaces []ClientSpace, playerIds map[string]bool) (requests []model.Request) {
	requests = make([]model.Request, 0)
	for _, spc := range clientSpaces {
		if v, ok := playerIds[spc.Id]; !ok || !v {
			continue
		}
		rtuples, _ := spc.GetAll(&model.Request{})
		if len(rtuples) <= 0 {
			continue
		}
		requests = append(requests, rtuples[0].GetFieldAt(0).(model.Request))
		last := len(requests) - 1
		for _, rtuple := range rtuples[1:] {
			request := rtuple.GetFieldAt(0).(model.Request)
			if request.Timestamp > requests[last].Timestamp {
				requests[last] = request
			}
		}
		requests[len(requests)-1].PlayerId = spc.Id
	}
	return requests
}

func SetupSpace(uri string) space.Space {
	// Register models for encoding to spc
	gob.Register(model.Request{})
	gob.Register([]model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Zombie{})
	gob.Register([]model.Zombie{})
	gob.Register(model.Shot{})
	gob.Register([]model.Shot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.State{})

	spc := space.NewSpace(uri)

	// todo: pSpaces seems to need this to be able to Get/Query on clients
	spc.QueryP(&model.Request{})
	spc.QueryP(&[]model.Request{})
	spc.QueryP(&model.Player{})
	spc.QueryP(&[]model.Player{})
	spc.QueryP(&model.Zombie{})
	spc.QueryP(&[]model.Zombie{})
	spc.QueryP(&model.Shot{})
	spc.QueryP(&[]model.Shot{})
	spc.QueryP(&model.Map{})
	spc.QueryP(&model.Wall{})
	spc.QueryP(&model.Segment{})
	spc.QueryP(&model.Point{})
	spc.QueryP(&model.State{})

	return spc
}
