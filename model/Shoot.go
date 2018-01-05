package model

import (
	"github.com/faiface/pixel"
	"time"
)

type Shoot struct{
	Start pixel.Vec
	Angle float64
	StartTime time.Time
	Weapon *Weapon
}

func (s Shoot) GetPos() (v pixel.Vec) {
	dt := time.Now().Sub(s.StartTime).Seconds() / 1000
	newPos := pixel.V(s.Weapon.Speed,0).Scaled(dt).Rotated(s.Angle)

	return s.Start.Add(newPos)
}


