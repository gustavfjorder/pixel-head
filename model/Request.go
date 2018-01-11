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

func (request *Request) Merge(other Request) {
	if request.PlayerId != other.PlayerId {
		return
	}
	request.Move = request.Move || other.Move
	switch { // Handles direction
	case request.Dir == math.NaN():
		request.Dir = other.Dir
	case other.Dir == math.NaN():
		request.Dir = request.Dir
	default:
		if request.Timestamp > other.Timestamp {
			request.Dir = request.Dir
		}else {
			request.Dir = other.Dir
		}
	}
	switch {
	case request.Reload || other.Reload:
		request.Reload = true
	case request.Shoot || other.Shoot:
		request.Shoot = true
	case request.Melee || other.Melee:
		request.Melee = true
	}
	if other.Timestamp > request.Timestamp {
		request.Weapon = other.Weapon
	}
}