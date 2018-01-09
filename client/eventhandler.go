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
		state := getState(spc, stateLock.State)
		stateLock.Mutex.Lock()
		stateLock.State = state
		stateLock.Mutex.Unlock()
	}
}

func getState(spc space.Space, oldState model.State) (state model.State) {
	spc.Get("ready")

	playerTuples, err := spc.GetP("players", &[]model.Player{})
	if err == nil {
		state.Players = playerTuples.GetFieldAt(1).([]model.Player)
	} else {
		state.Players = oldState.Players
	}

	zombieTuples, err := spc.GetP("zombies", &[]model.Zombie{})
	if err == nil {
		state.Zombies = zombieTuples.GetFieldAt(1).([]model.Zombie)
	} else {
		state.Zombies = oldState.Zombies
	}

	shootTuples, err := spc.GetP("shoots", &[]model.Shoot{})
	if err == nil {
		state.Shoots = shootTuples.GetFieldAt(1).([]model.Shoot)
	} else {
		state.Shoots = oldState.Shoots
	}

	spc.Put("done")

	return
}
