package model

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"fmt"
	"math"
)

type Barrel struct {
	Id       string
	Pos      pixel.Vec
	Exploded bool
	Dir      float64
}

func NewBarrel(pos pixel.Vec) Barrel {
	return Barrel{
		Id:  xid.New().String(),
		Pos: pos,
		Dir: 0,
	}
}

func (barrel *Barrel) Explode(s *State) {
	if barrel.Exploded {
		return
	}
	barrel.Exploded = true
	fmt.Println(barrel.ID())

	for index, player := range s.Players {
		if distanceBetween(player.Pos, barrel.Pos) < barrel.GetRange() {
			s.Players[index].Health -= int(barrel.GetPower() * (barrel.GetRange() - distanceBetween(player.Pos, barrel.Pos)) / barrel.GetRange())
		}
	}
	for index, zombie := range s.Zombies {
		if distanceBetween(zombie.Pos, barrel.Pos) < barrel.GetRange() {
			s.Zombies[index].Stats.Health -= int(barrel.GetPower() * (barrel.GetRange() - distanceBetween(zombie.Pos, barrel.Pos)) / barrel.GetRange())
		}
	}
	for i := range s.Barrels {
		b := &s.Barrels[i]
		if distanceBetween(b.Pos, barrel.Pos) < barrel.GetRange() && distanceBetween(b.Pos, barrel.Pos) != 0 {
			b.Explode(s)
		}

	}
}
func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}
func (barrel Barrel) GetHitbox() float64 {
	return 30
}

func (b Barrel) GetPower() float64 {
	return 50
}

func (b Barrel) GetRange() float64 {
	return 500
}

func (barrel Barrel) ID() string {
	return barrel.Id
}

func (barrel Barrel) EntityType() EntityType {
	return BarrelE
}

func (barrel Barrel) GetPos() pixel.Vec {
	return barrel.Pos
}

func (barrel Barrel) GetDir() float64 {
	return barrel.Dir
}
