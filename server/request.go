package server

import (
	"github.com/gustavfjorder/pixel-head/server/model"
)

type Request struct {
	PlayerId string
	CurrentWep model.Weapon
	Dir float64
	Move bool
	Shoot bool
	Reload bool
}
