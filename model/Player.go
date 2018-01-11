package model

import (
	"github.com/faiface/pixel"
	"math"
	"time"
	"github.com/pkg/errors"
	"math/rand"
)

type Player struct {
	Id          string
	Pos         pixel.Vec
	Dir         float64
	WeaponType  WeaponType
	WeaponList  []Weapon
	Reload      bool
	Shoot       bool
	Melee       bool
	Moved       bool
	Stats
	ActionDelay time.Duration
	TurnDelay   time.Duration
}

func NewPlayer(id string, pos ...pixel.Vec) (player Player) {
	player.WeaponList = make([]Weapon, nWeapon)
	player.WeaponType = HANDGUN
	player.NewWeapon(NewWeapon(KNIFE), NewWeapon(player.WeaponType))
	player.Id = id
	player.Stats = NewStats(HUMAN)
	player.Dir = 0
	if len(pos) > 0 {
		player.Pos = pos[0]
	} else {
		player.Pos = pixel.V(rand.Float64()*1000, rand.Float64()*1000)
	}
	return
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

func (player *Player) NewWeapon(weapons ...Weapon) {
	for _, weapon := range weapons{
		player.WeaponList[weapon.weaponType] = weapon
	}

}

func (player *Player) GetWeapon() (weapon *Weapon,e error) {
	if player.WeaponType >= nWeapon || player.WeaponType < 0{
		weapon = &Weapon{}
		e = errors.New("Unable to change weapon")
		return
	}
	weapon = &player.WeaponList[player.WeaponType]
	return weapon, e

}

func (player *Player) ChangeWeapon(weaponType WeaponType) {
	if player.IsAvailable(weaponType) {
		player.WeaponType = weaponType
	}
}

func (player *Player) IsAvailable(weaponType WeaponType) bool {
	return player.WeaponList[weaponType] != Weapon{}
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
