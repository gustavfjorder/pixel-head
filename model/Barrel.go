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
}

func NewBarrel(pos pixel.Vec) Barrel {
	return Barrel{
		Id:  xid.New().String(),
		Pos: pos,
	}
}


func (barrel Barrel) Explode(s *State) (exploded []Barrel){
	exploded= make([]Barrel,0)
	if barrel.Exploded{return}
	exploded=append(exploded, barrel)
	barrel.Exploded = true

	for index,player := range s.Players{
		if distanceBetween(player.Pos,barrel.Pos)<barrel.GetRange(){
			s.Players[index].Health-=int(barrel.GetPower()*(barrel.GetRange()-distanceBetween(player.Pos,barrel.Pos))/barrel.GetRange())
		}
	}
	for index,zombie := range s.Zombies{
		if distanceBetween(zombie.Pos,barrel.Pos)<barrel.GetRange(){
			s.Zombies[index].Stats.Health-=int(barrel.GetPower()*(barrel.GetRange()-distanceBetween(zombie.Pos,barrel.Pos))/barrel.GetRange())
			}
	}
	for index,b := range s.Barrels{
		if distanceBetween(b.Pos,barrel.Pos)<barrel.GetRange() && distanceBetween(b.Pos,barrel.Pos)!=0{
			//s.Barrels[barrelIndex]=s.Barrels[len(s.Barrels)-1]
			//s.Barrels=s.Barrels[:len(s.Barrels)-1]
			exploded=append(exploded,barrel.Explode(s)...)
			fmt.Println("barrel removed in Explode",index,b)
		}
	}
	return
}

func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}

func (barrel Barrel) GetHitBox() float64 {
	return 30
}

func (b Barrel) GetPower() float64{
	return 50
}

func (b Barrel) GetRange() float64{
	return 500
}

func (barrel Barrel) ID() string {
	return barrel.Id
}

func (barrel Barrel) EntityType()EntityType {
	return BarrelE
}