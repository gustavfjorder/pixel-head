package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"sync"
	"fmt"
	"time"
)

type StateLock struct {
	Mutex sync.Mutex
	State model.State
}

func HandleEvents(spc space.Space, stateLock *StateLock) {
	//Handle loop
	sec := time.Tick(time.Second)
	count := 0
	fmt.Println("Handling events")
	for {
		state := getState(spc)
		stateLock.State = state
		count++;
		select{
		case <-sec:
			fmt.Println("Handled:", count, "state updates")
			count = 0
		default:
			break
		}
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
