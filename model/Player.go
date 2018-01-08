package model

import "github.com/faiface/pixel"

type Player struct {
	Id         string
	Pos        pixel.Vec
	Dir        float64
	Weapon     Weapon
	WeaponList []bool
	Reload     bool
	Shoot      bool
	Melee      bool
	Moved      bool
	Stats
}

func NewPlayer(id string) Player {
	Weaponslist := make([]bool, len(Weapons))
	weapon := Knife
	Weaponslist[weapon] = true
	return Player{
		Id:         id,
		Pos:        pixel.V(200, 200),
		Weapon:     Weapon{},
		WeaponList: Weaponslist,
		Stats:      NewStats(Human),
	}
}

func (player Player) Move(dir float64) (Player) {
	player.Pos = player.Pos.Add(pixel.V(2, 0).Rotated(dir))
	return player
}

func (player *Player) NewWeapon(weapon Weapon) {
	player.WeaponList[weapon.Id] = true
}
