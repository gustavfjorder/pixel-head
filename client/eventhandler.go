package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"time"
)

func HandleEvents(spc *space.Space, state *model.State) {
	//Handle loop
	sec := time.Tick(time.Second)
	count := 0
	fmt.Println("Handling events")
	for {
		_, err := spc.Get("state", state)
		if err != nil {
			continue
		}

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

func GetPlayer(players []model.Player, player *model.Player){
	for _, p := range players {
		if player.Id == p.Id {
			*player = p
		}
	}
}