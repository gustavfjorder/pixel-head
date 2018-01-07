package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Zombie struct{
	Id string
	Pos pixel.Vec
	Dir float64
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

func (zombie *Zombie) Attack(players []Player) {
	for _, player := range players {
		if zombie.Pos.Sub(player.Pos).Len() < 1 {
			zombie.Pos = zombie.Pos.Rotated(angle(zombie.Pos, player.Pos))
			player.Stats.Health -= zombie.Stats.Power
			break
		}
	}
}


func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y - this.Y, other.X - this.X)
}