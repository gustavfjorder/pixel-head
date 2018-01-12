package model

import (
	"github.com/faiface/pixel"
	"time"
	"github.com/gustavfjorder/pixel-head/config"
	"math"
)

type Shot struct {
	Start      pixel.Vec
	Angle      float64
	StartTime  time.Duration
	WeaponType WeaponType
}

func NewShot(player Player, timeStamp time.Duration, angleOffset ...float64) (shot Shot) {
	shot.Start = player.Pos.Add(pixel.V(config.GunPosX, config.GunPosY).Rotated(player.Dir - math.Pi/2))
	shot.Angle = player.Dir
	shot.WeaponType = player.WeaponType
	shot.StartTime = timeStamp
	if len(angleOffset) > 0 {
		shot.Angle += angleOffset[0]
	}
	return
}

func (s Shot) GetPos(t time.Duration) (v pixel.Vec) {
	dt := float64(t-s.StartTime) / float64(time.Second.Nanoseconds())
	delta := pixel.V(s.WeaponType.ProjectileSpeed(), 0).Scaled(float64(dt)).Rotated(s.Angle)
	return s.Start.Add(delta)
}
