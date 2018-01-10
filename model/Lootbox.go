package model

import (
	"github.com/faiface/pixel"
	"math/rand"
)

type Lootbox struct {
	Id string
	Pos pixel.Vec
	Weapon
}

func NewLootbox(id string, x float64, y float64, randPos bool, weapon int) Lootbox{
	if randPos{
		return Lootbox{
			Id: id,
			//todo: multiply random numbers by map bounds
			Pos: pixel.V(rand.Float64(),rand.Float64()),
			Weapon: NewWeapon(weapon),
		}
	}
	return Lootbox{
		Id: id,
		Pos: pixel.V(x,y),
		Weapon: NewWeapon(weapon),
	}
}

func (player *Player) PickupLootbox(lootbox Lootbox){
	player.WeaponList[lootbox.Weapon.Id].Bullets+=lootbox.Weapon.MagazineCurrent
}

