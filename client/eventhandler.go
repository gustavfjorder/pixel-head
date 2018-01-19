package client

import (
	"github.com/pspaces/gospace/space"
	"github.com/gustavfjorder/pixel-head/model"
	"time"
	"github.com/gustavfjorder/pixel-head/config"
)

func HandleEvents(spc *space.Space, state *model.State,  updates chan<- model.Updates, done chan<- bool) {
	//Handle loop
	sec := time.NewTicker(time.Second)
	delay := time.NewTicker(config.Conf.ServerHandleSpeed)
	defer sec.Stop()
	defer delay.Stop()
	var tempState model.State
	count := 0
	for {
		_, err := spc.GetP("state", &tempState)
		if err == nil {
			count++
			*state = tempState
			model.Timestamp = state.Timestamp
		}
		updateTuples, err := spc.GetAll("update", &model.Updates{})
		for _, updateTuple := range updateTuples {
			updates <- updateTuple.GetFieldAt(1).(model.Updates)
		}
		if _, err := spc.GetP("game over"); err == nil {
			done <-true; done <- true
			break
		}
		<-delay.C
		select {
		case <-sec.C:
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
