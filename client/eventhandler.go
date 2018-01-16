package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"time"
	"github.com/gustavfjorder/pixel-head/config"
)

func HandleEvents(spc *space.Space, state *model.State,  updates chan<- model.Updates, done chan<- bool) {
	//Handle loop
	sec := time.NewTicker(time.Second)
	delay := time.NewTicker(config.Conf.ServerHandleSpeed)
	defer sec.Stop()
	defer delay.Stop()
	var t time.Duration
	var tempState model.State
	count := 0
	fmt.Println("Handling events")
	for {
		_, err := spc.GetP("state", &t, &tempState)
		if err == nil {
			count++
			model.Timestamp = t
			*state = tempState
		}
		updateTuples, err := spc.GetAll("update", &t, &model.Updates{})
		for _, updateTuple := range updateTuples {
			updates <- updateTuple.GetFieldAt(2).(model.Updates)
		}
		if _, err := spc.GetP("game over"); err == nil {
			fmt.Println("Ending game")
			done <-true; done <- true
			break
		}
		<-delay.C
		select {
		case <-sec.C:
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
		if player.ID() == p.ID() {
			*player = p
		}
	}

}
