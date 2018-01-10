package model

import (
	"github.com/faiface/pixel"
	"math"
)

type Request struct {
	Timestamp int64
	PlayerId  string
	Weapon    int
	Dir       float64
	Move      bool
	Shoot     bool
	Melee     bool
	Reload    bool
}

func (request Request) GetRotation() pixel.Matrix {
	return pixel.IM.Rotated(pixel.V(0,0), request.Dir)
}

func (request Request) Merge(other Request) (merged Request) {
	merged.Move = request.Move || other.Move
	merged.PlayerId = request.PlayerId
	switch { // Handles direction
	case request.Dir == math.NaN():
		merged.Dir = other.Dir
	case other.Dir == math.NaN():
		merged.Dir = request.Dir
	default:
		if request.Timestamp > other.Timestamp {
			merged.Dir = request.Dir
		}else {
			merged.Dir = other.Dir
		}
	}
	switch {
	case request.Reload || other.Reload:
		merged.Reload = true
	case request.Shoot || other.Shoot:
		merged.Shoot = true
	case request.Melee || other.Melee:
		merged.Melee = true
	}
	if request.Timestamp > other.Timestamp {
		merged.Weapon = request.Weapon
	}else{
		merged.Weapon = other.Weapon
	}
	return merged
}