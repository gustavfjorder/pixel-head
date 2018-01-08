package model

import "github.com/faiface/pixel"

type Player struct {
	Id         string
	Pos        pixel.Vec
	Dir        float64
	Weapon     int
	WeaponList []Weapon
	Reload     bool
	Shoot      bool
	Melee      bool
	Moved      bool
	Stats
}

func NewPlayer(id string) Player {
	weaponList := make([]Weapon, 0, len(Weapons))
	weaponList = append(weaponList, Weapons[Knife])

	return Player{
		Id:         id,
		Pos:        pixel.V(200, 200),
		Dir:        0,
		Weapon:     Knife,
		WeaponList: weaponList,
		Stats:      NewStats(Human),
	}
}

func (player Player) Move(dir float64) (Player) {
	player.Dir = dir
	player.Pos = player.Pos.Add(pixel.V(player.Stats.MoveSpeed, 0).Rotated(dir))
	return player
}

func (player *Player) NewWeapon(weapon Weapon) {
	player.WeaponList = append(player.WeaponList, weapon)
}

func (player *Player) GetWeapon() *Weapon {
	return &player.WeaponList[player.Weapon]
}
