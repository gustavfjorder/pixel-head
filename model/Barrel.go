package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Barrel struct {
	Id       string
	Pos      pixel.Vec
	Exploded bool
}

func NewBarrel(id string, pos pixel.Vec) Barrel {
	return Barrel{
		Id:  id,
		Pos: pos,
	}
}

func (barrel *Barrel) Explode(s *State) {
	var (
		Range = 500.
	)

	for index, player := range s.Players {
		if distanceBetween(player.Pos, barrel.Pos) < Range {
			s.Players[index].Health -= barrel.GetPower()
		}
	}
	for index, zombie := range s.Zombies {
		if distanceBetween(zombie.Pos, barrel.Pos) < Range {
			s.Zombies[index].Stats.Health -= barrel.GetPower()
		}
	}
	for _, b := range s.Barrels {
		if distanceBetween(b.Pos, barrel.Pos) < Range && distanceBetween(b.Pos, barrel.Pos) != 0 {
			barrel.Explode(s)
		}
	}
	barrel.Exploded = true
}

func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}

func (barrel Barrel) GetHitBox() float64 {
	return 30
}

func (barrel Barrel) GetPower() int {
	return 50
}
