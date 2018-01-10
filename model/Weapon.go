package model

import (
	"time"
	"math"
	"github.com/faiface/pixel"
	"github.com/gustavfjorder/pixel-head/config"
)

type Weapon struct {
	Id             int
	Name           string
	Power          int
	Range          float64
	Speed          float64
	Magazine       int
	Capacity       int
	Bullets        int
	BulletsPerShot int
	Spread         int
}

const (
	KNIFE            = iota
	RIFLE
	SHOTGUN
	HANDGUN
)

var Weapons = map[int]Weapon{
	RIFLE:
	{
		Id:             RIFLE,
		Name:           "rifle",
		Power:          10,
		Range:          300,
		Speed:          1000,
		Magazine:       30,
		Capacity:       150,
		Bullets:        0,
		BulletsPerShot: 1,
	},
	KNIFE:
	{
		Id:       KNIFE,
		Name:     "knife",
		Power:    100,
		Range:    20,
		Speed:    1000,
		Magazine: -1,
		Capacity: -1,
		Bullets:  -1,
	},
	SHOTGUN:
	{
		Id:             SHOTGUN,
		Name:           "shotgun",
		Power:          5,
		Range:          300,
		Speed:          1000,
		Magazine:       3,
		Capacity:       24,
		Bullets:        0,
		BulletsPerShot: 5,
		Spread:         5,
	},
	HANDGUN:
	{
		Id:             HANDGUN,
		Name:           "handgun",
		Power:          5,
		Range:          300,
		Speed:          1000,
		Magazine:       10,
		Capacity:       50,
		Bullets:        0,
		BulletsPerShot: 1,
	},
}

func (weapon *Weapon) RefillMag() {
	weapon.Magazine = Weapons[weapon.Id].Magazine
	weapon.Bullets -= Weapons[weapon.Id].Magazine
}

func (weapon *Weapon) GenerateShoots(timestamp int64, player Player) []Shoot {
	shotsPerSideOfDirection := int(math.Floor(float64(weapon.BulletsPerShot / 2)))
	angle := -(shotsPerSideOfDirection * weapon.BulletsPerShot)

	shoots := make([]Shoot, weapon.BulletsPerShot)

	for i := 0; i < weapon.BulletsPerShot; i++ {
		shoots[i] = Shoot{
			Start:     player.Pos.Add(pixel.V(config.GUNPOSX, config.GUNPOSY).Rotated(player.Dir - math.Pi/2)),
			Angle:     player.Dir + (float64(angle) * (math.Pi / 180)),
			StartTime: time.Now().UnixNano(),
			Weapon:    weapon.Id,
		}

		angle += weapon.Spread
		weapon.Magazine--
	}

	return shoots
}
