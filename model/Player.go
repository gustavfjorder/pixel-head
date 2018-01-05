package model

import "github.com/faiface/pixel"

type Player struct {
	Id  string
	Pos pixel.Vec
	Weapon Weapon
	WeaponList []Weapon
	Stats
}

func NewPlayer(id string) Player {
	var Weaponslist = make([]Weapon, len(Weapons))
	Weaponslist[Knife]=Weapons[Knife]
	return Player{
		Id: id,
		Pos: pixel.V(200,200),
		Weapon: Weapons[Knife],
		WeaponList: Weaponslist,
		Stats: NewStats(Human),
	}
}

func (player Player) Move(dir float64) (Player) {
	player.Pos = player.Pos.Add(pixel.V(2, 0).Rotated(dir))
	return player
}

func (player *Player) NewWeapon (weapon Weapon){
	player.WeaponList[weapon.Id]=weapon
}
