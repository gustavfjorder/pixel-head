package model

import (
	"github.com/faiface/pixel"
	"time"
	"github.com/gustavfjorder/pixel-head/config"
	"math"
	"github.com/rs/xid"
)

type Shot struct {
	Id         string
	Start      pixel.Vec
	Angle      float64
	StartTime  time.Duration
	WeaponType WeaponI
	Hit        bool
}

func NewShot(player Player, angleOffset ...float64) (shot Shot) {
	shot.Start = player.Pos.Add(pixel.V(config.GunPosX, config.GunPosY).Rotated(player.Dir - math.Pi/2))
	shot.Angle = player.Dir
	shot.WeaponType, _ = player.Weapon()
	shot.StartTime = Timestamp
	shot.Id = xid.New().String()
	if len(angleOffset) > 0 {
		shot.Angle += angleOffset[0]
	}
	return
}

func (s Shot) GetPos() (v pixel.Vec) {
	dt := float64(Timestamp-s.StartTime) / float64(time.Second.Nanoseconds())
	delta := pixel.V(s.WeaponType.ProjectileSpeed(), 0).Scaled(float64(dt)).Rotated(s.Angle)
	return s.Start.Add(delta)
}

func (s Shot) ID() string {
	return s.Id
}

func (s Shot) EntityType() EntityType {
	return ShotE
}

func (s Shot) GetHitbox() float64 {
	return 1
}

func (s Shot) GetDir() float64 {
	return s.Angle - math.Pi/2
}