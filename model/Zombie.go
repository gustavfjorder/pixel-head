package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/rs/xid"
)

type Zombie struct {
	Id        string
	Pos       pixel.Vec
	Dir       float64
	Attacking bool
	Stats     Stats
}

func NewZombie(x, y float64) Zombie {
	return Zombie{
		Id:         xid.New().String(),
		Pos:        pixel.V(x, y),
		Dir:        0,
		Stats:      NewStats(ZOMBIE),
	}
}


func (zombie* Zombie) Move(players []Player) {
	closestPlayer := players[0]

	for _, player := range players {
		if player.Pos.Sub(zombie.Pos).Len() < closestPlayer.Pos.Sub(zombie.Pos).Len() {
			closestPlayer = player
		}
	}

	angle := angle(zombie.Pos, closestPlayer.Pos)

	zombie.Pos = zombie.Pos.Add(pixel.V(zombie.Stats.MoveSpeed, 0).Rotated(angle))
	zombie.Dir = angle
}

func (zombie *Zombie) Attack(players []Player) {
	for i := range players {
		player := &players[i]

		zombie.Attacking = false
		if zombie.Pos.Sub(player.Pos).Len() < 3 {
			zombie.Dir = angle(zombie.Pos, player.Pos)
			player.Stats.Health -= zombie.Stats.Power
			zombie.Attacking = true
			break
		}
	}
}

func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y - this.Y, other.X - this.X)
}