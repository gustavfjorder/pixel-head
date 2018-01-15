package setup

import (
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"time"
)

func RegisterModels() {
	// Register models for encoding to spc
	// Register models for encoding to spc
	gob.Register(model.Request{})
	gob.Register([]model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Shot{})
	gob.Register([]model.Shot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Segment{})
	gob.Register(model.Point{})
	gob.Register(model.State{})
	gob.Register(model.Updates{})
	gob.Register(&model.Barrel{})
	gob.Register(&model.WeaponBase{})
	gob.Register(&model.Shotgun{})
	gob.Register(&model.Handgun{})
	gob.Register(&model.Rifle{})
	gob.Register(&model.Knife{})
	gob.Register([]model.FastZombie{})
	gob.Register([]model.SlowZombie{})
	gob.Register([]model.BombZombie{})
	gob.Register([]model.Zombie{})
	gob.Register(&model.FastZombie{})
	gob.Register(&model.SlowZombie{})
	gob.Register(&model.BombZombie{})
	gob.Register(&model.Zombie{})
	gob.Register(model.Lootbox{})
	var t time.Duration
	gob.Register(t)
}