package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Barrel struct{
	Id string
	Pos pixel.Vec
}
func NewBarrel(id string, pos pixel.Vec) Barrel{
	return Barrel{
		Id: id,
		Pos: pos,

	}
}

func (s *State) Explode(id string){
	var(
		barrel Barrel
		Range=500.
	)
	for _,b := range s.Barrels{
		if b.Id==id{
			barrel=b
			break
		}
	}
	for _,player := range s.Players{
		if distanceBetween(player.Pos,barrel.Pos)<Range{
			player.Health-=50
		}
	}
	for _,zombie := range s.Zombies{
		if distanceBetween(zombie.Pos,barrel.Pos)<Range{
			zombie.Stats.Health-=50
		}
	}
	for _,b := range s.Barrels{
		if distanceBetween(b.Pos,barrel.Pos)<Range && distanceBetween(b.Pos,barrel.Pos)!=0{
			s.Explode(b.Id)
		}
	}
}

func distanceBetween(pos1 pixel.Vec, pos2 pixel.Vec) float64 {
	return math.Sqrt(math.Abs(pos1.X-pos2.X)*math.Abs(pos1.X-pos2.X) + math.Abs(pos1.Y-pos2.Y)*math.Abs(pos1.Y-pos2.Y))
}