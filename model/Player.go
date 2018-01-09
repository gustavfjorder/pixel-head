package model

import (
	"github.com/faiface/pixel"
	"math"
	"fmt"
)

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
	weaponList := make([]Weapon, len(Weapons))
	weaponList[Knife] = Weapons[Knife]
	weaponList[Handgun] = Weapons[Handgun]
	weaponList[Rifle] = Weapons[Rifle]
	weaponList[Shotgun] = Weapons[Shotgun]
	return Player{
		Id:         id,
		Pos:        pixel.V(200, 200),
		Dir:        0,
		Weapon:     Handgun,
		WeaponList: weaponList,
		Stats:      NewStats(HUMAN),
	}
}

func (player *Player) Move(dir float64, m Map) {
	if dir != math.NaN() {
		newpos := player.Pos.Add(pixel.V(player.Stats.MoveSpeed, 0).Rotated(dir))
		for _, wall := range m.Walls {
			if wall.Intersect(NewLine(PointFrom(player.Pos), PointFrom(newpos))){
				fmt.Println("Invalid move")
				return
			}
		}
		player.Dir = dir
		player.Pos = newpos
	}
}

func (player *Player) NewWeapon(weapon Weapon) {
	player.WeaponList[weapon.Id] = weapon
}

func (player Player) GetWeapon() *Weapon {
	return &player.WeaponList[player.Weapon]
}

func (player *Player) ChangeWeapon(weapon int) {
	for i := range player.WeaponList {
		if i == weapon {
			player.Weapon = weapon
			break
		}
	}
}

func (player *Player) IsAvailable(weapon int) bool{
	return weapon < len(player.WeaponList) && player.WeaponList[weapon] != (Weapon{})
}
