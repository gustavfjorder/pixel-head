package model

import (
	"github.com/faiface/pixel"
)

type Lootbox struct {
	Id string
	Pos pixel.Vec
	Weapon
}

func NewLootbox(id string, x float64, y float64, weapon int) Lootbox{
	return Lootbox{
		Id: id,
		Pos: pixel.V(x,y),
		Weapon: Weapons[weapon],
	}
}

func (player *Player) PickupLootbox(lootbox Lootbox){
	player.WeaponList[lootbox.Weapon.Id].Bullets+=lootbox.Weapon.Magazine
}