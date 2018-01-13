package model

import (
	"github.com/faiface/pixel"
	"math"
	"fmt"
)

type Barrel struct{
	Id string
	Pos pixel.Vec
	Range float64
}
func NewBarrel(id string, pos pixel.Vec) Barrel{
	return Barrel{
		Id: id,
		Pos: pos,
		Range: 200,
	}
}


func (barrel Barrel) Explode(barrelIndex int, s *State){


	for index,player := range s.Players{
		if distanceBetween(player.Pos,barrel.Pos)<barrel.Range{
			s.Players[index].Health-=int(barrel.GetPower()/distanceBetween(player.Pos,barrel.Pos))
		}
	}
	for index,zombie := range s.Zombies{
		if distanceBetween(zombie.Pos,barrel.Pos)<barrel.Range{
			s.Zombies[index].Stats.Health-=int(barrel.GetPower()*(barrel.Range-distanceBetween(zombie.Pos,barrel.Pos))/barrel.Range)
			}
	}
	for index,b := range s.Barrels{
		if distanceBetween(b.Pos,barrel.Pos)<barrel.Range && distanceBetween(b.Pos,barrel.Pos)!=0{
			/*s.Barrels[barrelIndex]=s.Barrels[len(s.Barrels)-1]
			s.Barrels=s.Barrels[:len(s.Barrels)-1]
			barrel.Explode(index,s)*/
			fmt.Println("barrel removed in Explode",index,b)
		}
	}
}

func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}

func (b Barrel) GetHitBox() float64{
	return 30
}

func (b Barrel) GetPower() float64{
	return 50
}