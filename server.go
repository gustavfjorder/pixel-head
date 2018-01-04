package main

import (
	"fmt"
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
	"github.com/gustavfjorder/pixel-head/server/model"
)

const N = 1

func main() {
	spc := NewSpace("tcp://localhost:31414/room1")
	var id string
	awaiting := make(map[string]bool, N)
	for {
		spc.Get("client", &id)
		fmt.Println("player", id, "connected")
		if _, ok := awaiting[id]; !ok {
			awaiting[id] = true
		}
		if len(awaiting) >= N {
			cpy := make([]model.Player, N)
			i := 0
			for k, _ := range awaiting {
				spc.Put(k, "ready")
				cpy[i] = model.Player{Id: k, Pos: pixel.V(0, 0)}
				i++
			}
			awaiting = make(map[string]bool, N)
			go server.StartGame(cpy)
		}
	}
}
