package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Request struct {
	Timestamp int64
	PlayerId string
	CurrentWep int
	Dir float64
	Move bool
	Shoot bool
	Melee bool
	Reload bool
}

func (r Request) WeaponName() string{
	return Weapons[r.CurrentWep].Name
}

func (r Request) GetRotation() pixel.Matrix {
	return pixel.IM.Rotated(pixel.V(0,0), r.Dir)
}

func (this Request) Merge(other Request) (merged Request) {
	merged.Move = this.Move || other.Move
	merged.PlayerId = this.PlayerId
	switch { // Handles direction
	case this.Dir == math.NaN():
		merged.Dir = other.Dir
	case other.Dir == math.NaN():
		merged.Dir = this.Dir
	default:
		if this.Timestamp > other.Timestamp {
			merged.Dir = this.Dir
		}else {
			merged.Dir = other.Dir
		}
	}
	switch {
	case this.Reload || other.Reload:
		merged.Reload = true
	case this.Shoot || other.Shoot:
		merged.Shoot = true
	case this.Melee || other.Melee:
		merged.Melee = true
	}
	if this.Timestamp > other.Timestamp {
		merged.CurrentWep = this.CurrentWep
	}else{
		merged.CurrentWep = other.CurrentWep
	}
	return merged
}