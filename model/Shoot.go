package model

import (
	"github.com/faiface/pixel"
	"time"
)

type Shoot struct {
	Start     pixel.Vec
	Angle     float64
	StartTime int64
	Weapon    int
}

func (s Shoot) GetPos(t int64) (v pixel.Vec) {
	dt := float64((t - s.StartTime)) / float64(time.Second.Nanoseconds())
	delta := pixel.V(Weapons[s.Weapon].Speed, 0).Scaled(float64(dt)).Rotated(s.Angle)
	return s.Start.Add(delta)
}
