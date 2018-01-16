package model

import (
	"github.com/faiface/pixel"
	"math"
	"time"
)

type Request struct {
	Timestamp time.Duration
	PlayerId  string
	Weapon    WeaponType
	Dir       float64
	Move      bool
	Action    Action
}

func (request Request) GetRotation() pixel.Matrix {
	return pixel.IM.Rotated(pixel.V(0, 0), request.Dir)
}

func (request Request) Moved() bool {
	return !math.IsNaN(request.Dir)
}

func (request Request) Reload() bool {
	return request.Action == RELOAD
}

func (request Request) Shoot() bool {
	return request.Action == SHOOT
}

func (request Request) Melee() bool {
	return request.Action == MELEE
}

func (request Request) Idle() bool {
	return request.Action == IDLE
}

func (request Request) Barrel() bool {
	return request.Action == BARREL
}

func (request Request) Valid(me Player) bool {
	return request.Action != IDLE || me.WeaponType != request.Weapon
}