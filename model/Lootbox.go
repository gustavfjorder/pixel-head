package model

import (
	"github.com/faiface/pixel"
	"github.com/rs/xid"
	"math/rand"
)

type Lootbox struct {
	Id     string
	Pos    pixel.Vec
	Weapon Weapon
}

func NewLootbox(x, y float64) Lootbox {
	weapon := rand.Intn(int(nWeapon) - 1) + 1 // Knife not allowed in lootbox

	return Lootbox{
		Id:     xid.New().String(),
		Pos:    pixel.V(x, y),
		Weapon: NewWeapon(WeaponType(weapon)),
	}
}

func (player *Player) PickupLootbox(lootbox *Lootbox) {
	player.WeaponList[lootbox.Weapon.WeaponType].Bullets += lootbox.Weapon.MagazineCurrent
}

func (lootbox Lootbox) ID() string {
	return lootbox.Id
}

func (lootbox Lootbox) EntityType() EntityType {
	return LootboxE
}

