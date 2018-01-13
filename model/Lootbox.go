package model

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
)

type Lootbox struct {
	Id     string
	Pos    pixel.Vec
	Weapon Weapon
}

func NewLootbox(x, y float64, weapon WeaponType) Lootbox {
	return Lootbox{
		Id:     xid.New().String(),
		Pos:    pixel.V(x, y),
		Weapon: NewWeapon(weapon),
	}
}

func (player *Player) PickupLootbox(lootbox *Lootbox) {
	player.WeaponList[lootbox.Weapon.WeaponType].Bullets += lootbox.Weapon.MagazineCurrent
}

