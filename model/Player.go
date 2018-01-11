package model

import (
	"github.com/faiface/pixel"
	"math"
	"fmt"
	"time"
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
	ActionDelay time.Duration
	TurnDelay time.Duration
}

func NewPlayer(id string) Player {

	weaponList := make([]Weapon, 0)
	weaponList = append(weaponList, NewWeapon(KNIFE), NewWeapon(HANDGUN), NewWeapon(SHOTGUN))
	return Player{
		Id:         id,
		Pos:        pixel.V(200, 200),
		Dir:        0,
		Weapon:     HANDGUN,
		WeaponList: weaponList,
		Stats:      NewStats(HUMAN),
	}
}

func (player *Player) Move(dir float64, g *Game) {
	if dir != math.NaN() {
		if player.TurnDelay < g.State.Timestamp {
			player.Dir = dir
			player.TurnDelay = player.GetTurnDelay() + g.State.Timestamp
		}
		fmt.Println("moved player", player)
		newpos := player.Pos.Add(pixel.V(player.Stats.GetMoveSpeed(), 0).Rotated(player.Dir))
		for _, wall := range g.CurrentMap.Walls {
			if wall.Intersect(NewLine(PointFrom(player.Pos), PointFrom(newpos))) {
				fmt.Println("Invalid move")
				return
			}
		}
		player.Pos = newpos
	}
}

func (player *Player) NewWeapon(weapon Weapon) {
	if !player.IsAvailable(weapon.Id) {
		player.WeaponList = append(player.WeaponList, weapon)
	}
}

func (player *Player) GetWeapon() *Weapon {
	if player.Weapon < len(player.WeaponList) {
		return &player.WeaponList[player.Weapon]
	} else if len(player.WeaponList) > 0 {
		player.Weapon = 0
		return &player.WeaponList[player.Weapon]
	} else {
		player.WeaponList = append(player.WeaponList, NewWeapon(KNIFE))
		player.Weapon = 0
		return &player.WeaponList[player.Weapon]
	}
}

func (player *Player) ChangeWeapon(weaponNum int) {
	for i, weapon := range player.WeaponList {
		if weapon.Id == weaponNum {
			player.Weapon = i
			break
		}
	}
}

func (player *Player) IsAvailable(weaponNum int) bool {
	for _, weapon := range player.WeaponList {
		if weapon.Id == weaponNum {
			return true
		}
	}
	return false
}

func (player Player) GetTurnDelay() time.Duration {
	return time.Second / 15
}
