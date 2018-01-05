package model

import (
	"github.com/faiface/pixel"
	"time"
)

type Request struct {
	Timestamp time.Time
	PlayerId string
	CurrentWep int
	Dir float64
	Move bool
	Shoot bool
	Melee bool
	Reload bool
}

func (r Request) MovementName() string{
	switch{
	case r.Reload:
		return "reload"
	case r.Shoot:
		return "shoot"
	case r.Melee:
		return "meleeattack"
	case r.Move:
		return "move"
	default:
		return "idle"
	}
}

func (r Request) WeaponName() string{
	return Weapons[r.CurrentWep].Name
}

func (r Request) GetRotation() pixel.Matrix {
	return pixel.IM.Rotated(pixel.V(0,0), r.Dir)
}