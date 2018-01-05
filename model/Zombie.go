package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Zombie struct{
	Pos pixel.Vec
	Stats Stats
}


func (zombie* Zombie) Move(players []Player) {
	closestPlayer := players[0]

	for _, player := range players {
		if player.Pos.Sub(zombie.Pos).Len() < closestPlayer.Pos.Sub(zombie.Pos).Len() {
			closestPlayer = player
		}
	}

	move := pixel.V(zombie.Stats.MoveSpeed, 0).Rotated(angle(zombie.Pos, closestPlayer.Pos))

	zombie.Pos = zombie.Pos.Add(move)
}

func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y - this.Y, other.X - this.X)
}