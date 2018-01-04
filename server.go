package main

import (
	"fmt"
	. "github.com/pspaces/gospace"
	"github.com/gustavfjorder/pixel-head/server"
)

const N = 1

func main() {
	lounge := NewSpace("tcp://localhost:31414/lounge")

	awaiting := make([]string, 0, N)

	for {
		var id string
		lounge.Get("client", &id)

		fmt.Println("Player '" + id + "' has connected")

		awaiting = append(awaiting, id)

		if len(awaiting) == cap(awaiting) {
			for _, id := range awaiting {
				lounge.Put(id, "ready")
			}

			go server.StartGame(awaiting)
		}
	}
}
