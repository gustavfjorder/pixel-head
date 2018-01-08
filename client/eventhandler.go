package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"sync"
)

type StateLock struct {
	Mutex sync.Mutex
	State model.State
}

func HandleEvents(spc space.Space, stateLock *StateLock) {
	//Handle loop
	for {
		state := getState(spc)
		stateLock.Mutex.Lock()
		stateLock.State = state
		stateLock.Mutex.Unlock()
	}
}

func getState(spc space.Space) (state model.State) {
	spc.Get("ready")

	playerTuples, _ := spc.GetP("players", &[]model.Player{})
	state.Players = playerTuples.GetFieldAt(1).([]model.Player)

	zombieTuples, _ := spc.GetP("zombies", &[]model.Zombie{})
	state.Zombies = zombieTuples.GetFieldAt(1).([]model.Zombie)

	shootTuples, _ := spc.GetP("shoots", &[]model.Shoot{})
	state.Shoots = shootTuples.GetFieldAt(1).([]model.Shoot)

	spc.Put("done")

	return
}
