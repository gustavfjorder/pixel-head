package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/rs/xid"
	"fmt"
)

type Barrel struct {
	Id       string
	Pos      pixel.Vec
	Exploded bool
	Dir float64
}

func NewBarrel(pos pixel.Vec) Barrel {
	return Barrel{
		Id:  xid.New().String(),
		Pos: pos,
		Dir: 0,
	}
}

func (barrel *Barrel) Explode(s *State){
	if barrel.Exploded {
		return
	}
	barrel.Exploded = true
	fmt.Println(barrel.ID())

	for index, player := range s.Players {
		if distanceBetween(player.Pos, barrel.Pos) <= barrel.GetRange() {
			s.Players[index].Health -= barrel.GetPower()
		}
	}
	for index, zombie := range s.Zombies {
		if distanceBetween(zombie.Pos, barrel.Pos) < barrel.GetRange() {
			s.Zombies[index].Stats.Health -= barrel.GetPower()
		}
	}
	for i := range s.Barrels {
		b := &s.Barrels[i]
		if distanceBetween(b.Pos, barrel.Pos) < barrel.GetRange() && b.ID() != barrel.ID(){
			b.Explode(s)
		}
	}
}

func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}

func (barrel Barrel) GetRange() float64{
	return 200
}

func (barrel Barrel) GetHitbox() float64 {
	return 30
}

func (barrel Barrel) GetPower() int {
	return 1
}

func (barrel Barrel) ID() string {
	return barrel.Id
}

func (barrel Barrel) EntityType()EntityType {
	return BarrelE
}

func (barrel Barrel) GetPos() pixel.Vec{
	return barrel.Pos
}

func (barrel Barrel) GetDir() float64{
	return barrel.Dir
}