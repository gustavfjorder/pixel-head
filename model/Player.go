package model

import "github.com/faiface/pixel"

type Player struct {
	Id  string
	Pos pixel.Vec
	Weapon
	Stats
}

func NewPlayer(id string) Player {
	return Player{
		Id: id,
		Pos: pixel.V(200,200),
		Weapon: Weapons[Handgun],
		Stats: NewStats(Human),
	}
}

func (p Player) Move(dir float64) (Player) {
	p.Pos = p.Pos.Add(pixel.V(2, 0).Rotated(dir))
	return p
}

