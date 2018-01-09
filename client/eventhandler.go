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

	s, err := spc.GetP("state", &model.State{})
	if err == nil {
		state = s.GetFieldAt(1).(model.State)
	} else {
		state = oldState
	}

	spc.Put("done")

	return
}
