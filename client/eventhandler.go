package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"time"
	"github.com/gustavfjorder/pixel-head/config"
)

func HandleEvents(spc *space.Space, state *model.State, me *model.Player) {
	//Handle loop
	sec := time.Tick(time.Second)
	count := 0
	fmt.Println("Handling events")
	for {
		_, err := spc.Get("state", state)
		if err != nil {
			continue
		}

		for _,p := range state.Players{
			if config.Conf.Id == p.Id {
				*me = p
			}
		}

		count++

		select{
		case <-sec:
			fmt.Println("Handled:", count, "state updates")
			count = 0
		default:
			break
		}
	}
}
