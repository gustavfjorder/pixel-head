package model

import "github.com/faiface/pixel"

type Player struct {
	Id  string
	Pos pixel.Vec
	Stats Stats
}

func (p Player) Move(dir float64) (Player) {
	p.Pos = p.Pos.Add(pixel.V(2, 0).Rotated(dir))
	return p
}

