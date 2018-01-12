package model

import (
	"github.com/faiface/pixel"
	"math"
	"time"
	"github.com/pkg/errors"
	"math/rand"
)

type Action int

const (
	RELOAD Action = iota
	SHOOT
	MELEE
	MOVE
	IDLE
)

type Player struct {
	Id          string
	Pos         pixel.Vec
	Dir         float64
	WeaponType  WeaponType
	WeaponList  []Weapon
	Action      Action
	Stats
}

var actionDelays = make(map[string]time.Duration)
var turnDelays = make(map[string]time.Duration)

func NewPlayer(id string, pos ...pixel.Vec) (player Player) {
	player.WeaponList = make([]Weapon, nWeapon)
	player.WeaponType = HANDGUN
	player.NewWeapon(NewWeapon(KNIFE), NewWeapon(player.WeaponType), NewWeapon(SHOTGUN))
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
	if !math.IsNaN(dir) {
		if player.TurnDelay() < g.State.Timestamp {
			player.Dir = dir
			player.Turn(g.State.Timestamp)
			return
		}
		newpos := player.Pos.Add(pixel.V(player.Stats.GetMoveSpeed(), 0).Rotated(player.Dir))
		for _, wall := range g.CurrentMap.Walls {
			if wall.Intersect(NewLine(PointFrom(player.Pos), PointFrom(newpos))) {
				return
			}
		}

		player.Pos = newpos
	}
}

func (player *Player) NewWeapon(weapons ...Weapon) {
	for _, weapon := range weapons {
		player.WeaponList[weapon.WeaponType] = weapon
	}

}

func (player *Player) GetWeapon() (weapon *Weapon, e error) {
	if player.WeaponType >= nWeapon || player.WeaponType < 0 {
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
	return time.Second / 8
}

func findPlayer(players []Player, id string) (p *Player, e error) {
	p = &Player{}
	for i, player := range players {
		if id == player.Id {
			p = &players[i]
			return
		}
	}
	e = errors.New("Unable to find player with ID: " + id)
	return
}

func (player *Player) SetAction(action Action, timestamp time.Duration){
	player.Action = action
	actionDelays[player.Id] = player.Delay(action) + timestamp
}

func (player Player) Delay(action Action) time.Duration{
	switch action {
	case RELOAD:
		return player.WeaponType.ReloadSpeed()
	case SHOOT, MELEE:
		return player.WeaponType.ShootDelay()
	default:
		return 0
	}
}



func (player Player) Turn(timestamp time.Duration) {
	turnDelays[player.Id] = timestamp + player.TurnSpeed()
}

func (player Player) ActionDelay() time.Duration {
	return actionDelays[player.Id]
}

func (player Player) TurnDelay() time.Duration {
	return turnDelays[player.Id]
}

func (player Player) TurnSpeed() time.Duration {
	return time.Second / 10
}

