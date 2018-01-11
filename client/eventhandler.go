package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"time"
	"sync"
)

func HandleEvents(spc *space.Space, state *model.State, lock *sync.Mutex) {
	//Handle loop
	sec := time.Tick(time.Second)
	count := 0
	fmt.Println("Handling events")
	for {

		stateTuple, err := spc.Get("state", &model.State{})
		if err != nil {
			continue
		}
		tempState := stateTuple.GetFieldAt(1).(model.State)
		lock.Lock()
		*state = tempState
		lock.Unlock()

		count++
		select {
		case <-sec:
			fmt.Println("Handled:", count, "state updates")
			count = 0
		default:
			break
		}
	}
}

func GetPlayer(players []model.Player, player *model.Player) {
	if player == nil {
		return
	}
	for _, p := range players {
		if player.Id == p.Id {
			*player = p
		}
	}

}
