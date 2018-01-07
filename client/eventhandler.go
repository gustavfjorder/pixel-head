package client

import (
	"github.com/pspaces/gospace/space"
	"time"
	"github.com/gustavfjorder/pixel-head/Config"
	"github.com/gustavfjorder/pixel-head/model"
	"sync"
)

type StateLock struct {
	Mutex sync.Mutex
	State model.State
}

func HandleEvents(spc space.Space, stateLock *StateLock) {
	var (
		state  = getState(spc)
		ticker = time.NewTicker(time.Second / Config.Conf.HandleFrequency).C
	)
	stateLock.Mutex.Lock()
	stateLock.State = state
	stateLock.Mutex.Unlock()

	//Handle loop
	for {
		<-ticker //Stop from starving server too much
		state = getState(spc)
		stateLock.Mutex.Lock()
		stateLock.State = state
		stateLock.Mutex.Unlock()
	}
}

func getState(spc space.Space) (state model.State) {
	var n int
	//Could possibly starve server
	spc.Get("rw_lock", &n)
	spc.Put("rw_lock", n+1)
	playerTuples, err := spc.QueryAll(&model.Player{})
	if err != nil {
		panic(err)
	}
	state.Players = make([]model.Player, len(playerTuples))
	for i, playerTuple := range playerTuples {
		state.Players[i] = playerTuple.GetFieldAt(0).(model.Player)
	}

	zoombieTuples, err := spc.QueryAll(&model.Zombie{})
	if err != nil {
		panic(err)
	}
	state.Zombies = make([]model.Zombie, len(zoombieTuples))
	for i, zoombieTuple := range zoombieTuples {
		state.Zombies[i] = zoombieTuple.GetFieldAt(0).(model.Zombie)
	}
	shootsTuple, err := spc.QueryAll(&model.Shoot{})
	if err != nil {
		panic(err)
	}
	state.Shoots = make([]model.Shoot, len(shootsTuple))
	for i, shootTuple := range shootsTuple {
		state.Shoots[i] = shootTuple.GetFieldAt(0).(model.Shoot)
	}
	spc.Get("rw_lock", &n)
	spc.Put("rw_lock", n-1)
	return
}
