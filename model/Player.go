package model

import (
	"github.com/faiface/pixel"
	"math"
	"time"
	"github.com/pkg/errors"
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
		Weapon:     0,
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
		newpos := player.Pos.Add(pixel.V(player.Stats.GetMoveSpeed() , 0).Rotated(player.Dir))
		for _, wall := range g.CurrentMap.Walls {
			if wall.Intersect(NewLine(PointFrom(player.Pos), PointFrom(newpos))) {
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

func (player *Player) GetWeapon() (weapon *Weapon,e error) {
	if player.Weapon < len(player.WeaponList){
		weapon = &player.WeaponList[player.Weapon]
	} else {
		panic(errors.New("Wopsi"))
	}
	return weapon, e

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

func findPlayer(players []Player, id string) (p *Player,e error){
	p = &Player{}
	for i, player := range players {
		if id == player.Id{
			p = &players[i]
			return
		}
	}
	e = errors.New("Unable to find player")
	return
}
