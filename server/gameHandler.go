package server

import (
	"time"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/config"
)

type ClientSpace struct {
	Id  string
	Uri string
	space.Space
}

func Start(g *model.Game, clientSpaces []ClientSpace, finished chan <- bool) {
	start := time.Now()

	for _, spc := range clientSpaces {
		spc.Put("map", g.CurrentMap)
	}

	for _, spc := range clientSpaces {
		spc.Get("joined")
	}

	t := time.Tick(config.Conf.ServerHandleSpeed)
	sec := time.Tick(time.Second)
	speed := config.Conf.ServerHandleSpeed
	for g.CurrentLevel < len(model.Levels) {

		duration := time.Second * 1
		if g.CurrentLevel == 0 {
			duration = 0
		}

		levelRdy := make(chan bool, 1)
		zombiesSpawned := false
		time.AfterFunc(duration, func() {
			g.PrepareLevel(levelRdy)
		})

		for {
			select {
			case <-levelRdy:
				zombiesSpawned = true
			default:
			}
			//Update game
			g.State.Timestamp = time.Since(start)
			g.Updates.Timestamp = g.State.Timestamp

			g.HandleRequests(collectRequests(clientSpaces, g.PlayerIds))
			g.HandleLoot()
			g.HandleBarrels()
			g.HandleZombies()
			g.HandleShots()
			g.HandlePlayers()
			g.HandleCorpses()

			//Send new game state to clients
			compressed := g.State.Compress()
			for _, spc := range clientSpaces {
				spc.GetP("state", &model.State{})
				spc.Put("state",compressed)
				if !g.Updates.Empty(){
					spc.Put("update", g.Updates)
				}
			}
			g.Updates.Clear()


			//If all players died end game
			if len(g.State.Players) == 0 {
				levelRdy <- true
				goto endgame
			}

			//If all zombies have been killed, go to next level
			if zombiesSpawned && len(g.State.Zombies) <= 0 {
				break
			}
			<-t
			select {
			case <-sec:
				for i := range g.State.Players {
					(&g.State.Players[i]).Regen()

				}
				config.Conf.ServerHandleSpeed = time.Second / speed
				speed = 0
			default:
				speed++
			}
		}

		//g.LevelDone(g.CurrentLevel)
		g.CurrentLevel++

	}
endgame:
	for _, spc := range clientSpaces {
		spc.Put("game over")
	}
	finished <- true
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
	spc.QueryP(&model.Updates{})
	spc.QueryP(&[]model.FastZombie{})
	spc.QueryP(&[]model.SlowZombie{})
	spc.QueryP(&[]model.BombZombie{})
	spc.QueryP(&[]model.Zombie{})
	spc.QueryP(&model.FastZombie{})
	spc.QueryP(&model.SlowZombie{})
	spc.QueryP(&model.BombZombie{})
	spc.QueryP(&model.Zombie{})
	var t time.Duration
	spc.QueryP(&t)

	return spc
}
