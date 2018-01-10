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
			//todo: instead of 900+100 put map dimensions
			Pos: pixel.V(rand.Float64()*900+100,rand.Float64()*900+100),
			Weapon: Weapons[weapon],
		}
	}
	return Lootbox{
		Id: id,
		Pos: pixel.V(x,y),
		Weapon: Weapons[weapon],
	}
}

func (player *Player) PickupLootbox(lootbox Lootbox){
	player.WeaponList[lootbox.Weapon.Id].Bullets+=lootbox.Weapon.Magazine
}

