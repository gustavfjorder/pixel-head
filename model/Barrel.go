package model

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"math"
)

type Barrel struct {
	BarrelI
	Id       string
	Pos      pixel.Vec
	Exploded bool
	Dir      float64
}

type BarrelI interface {
	EntityI
	Explode(*State)
	GetPower() float64
	GetRange() float64
	SetPos(pixel.Vec)
	IsExploded() bool
}

func NewBarrel(pos pixel.Vec) BarrelI {
	return &Barrel{
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

	for index, player := range s.Players {
		if distanceBetween(player.Pos, barrel.Pos) < barrel.GetRange() {
			s.Players[index].Health -= int(barrel.GetPower() * (barrel.GetRange() - distanceBetween(player.Pos, barrel.Pos)) / barrel.GetRange())
		}
	}
	for index := range s.Zombies {
		zombie := s.Zombies[index]
		if distanceBetween(zombie.GetPos(), barrel.Pos) < barrel.GetRange() {
			s.Zombies[index].SubHealth(int(barrel.GetPower() * (barrel.GetRange() - distanceBetween(zombie.GetPos(), barrel.Pos)) / barrel.GetRange()))
		}
	}
	for _, b := range s.Barrels {
		if distanceBetween(b.GetPos(), barrel.Pos) < barrel.GetRange() && b.ID() != barrel.ID() {
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

func (barrel Barrel) GetPower() float64 {
	return 5
}

func (barrel Barrel) GetRange() float64 {
	return 100
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

func (barrel *Barrel) SetPos(vec pixel.Vec) {
	barrel.Pos = vec
}

func (barrel *Barrel) IsExploded() bool {
	return barrel.Exploded
}
