package model

import "github.com/faiface/pixel"

type Player struct {
	Id  string
	Pos pixel.Vec
	Weapon int
	WeaponList []Weapon
	Stats
}

func NewPlayer(id string) Player {
	var Weaponslist = make([]Weapon, len(Weapons))
	Weaponslist[Knife]=Weapons[Knife]
	return Player{
		Id: id,
		Pos: pixel.V(200,200),
		Weapon: Knife,
		WeaponList: Weaponslist,
		Stats: NewStats(Human),
	}
}

func (p Player) Move(dir float64) (Player) {
	p.Pos = p.Pos.Add(pixel.V(2, 0).Rotated(dir))
	return p
}

func (player *Player) NewWeapon (weapon Weapon){
	player.WeaponList[weapon.Id]=weapon
}
