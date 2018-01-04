package model

import (
	"github.com/faiface/pixel"
	"time"
)

type Shoot struct{
	Start pixel.Vec
	Angle float64
	StartTime time.Time
	Speed float64
}

func (s Shoot) GetPos() (v pixel.Vec) {
	dt := time.Now().Sub(s.StartTime).Seconds() / 1000
	return s.Start.Add(pixel.V(s.Speed,0).Scaled(dt).Rotated(s.Angle))
}


