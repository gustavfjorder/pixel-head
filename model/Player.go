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
	Id         string
	Pos        pixel.Vec
	Dir        float64
	WeaponType WeaponType
	WeaponList []Weapon
	Action     Action
	Stats
}

var actionDelays = make(map[string]time.Duration)
var turnDelays = make(map[string]time.Duration)

func NewPlayer(id string, pos ...pixel.Vec) (player Player) {
	player.WeaponList = make([]Weapon, nWeapon)
	player.WeaponType = KNIFE
	player.NewWeapon(NewWeapon(player.WeaponType), NewWeapon(RIFLE), NewWeapon(SHOTGUN))
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
		if player.TurnDelay() < Timestamp {
			player.Dir = dir
			player.Turn()
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

func (player *Player) Weapon() (weapon *Weapon, e error) {
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

func (player *Player) SetAction(action Action){
	player.Action = action
	actionDelays[player.Id] = player.Delay(action) + Timestamp
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

func (player Player) Reload() bool{
	wep, err := player.Weapon()
	if err != nil {
		return false
	}
	player.SetAction(RELOAD)
	return wep.RefillMag()
}

func (player Player) EmptyMag() bool{
	wep, err := player.Weapon()
	if err != nil {
		return true
	}
	return wep.MagazineCurrent <= 0
}

func (player *Player) Shoot(g *Game) {
	weapon, err := player.Weapon()
	if err != nil {
		return
	}
	playerShoots := weapon.GenerateShoots(*player)
	for _, shot := range playerShoots {
		g.Add(shot)
	}

	player.SetAction(SHOOT)
}

func (player *Player) Do(request Request, g *Game) {
	if Timestamp < player.ActionDelay(){
		return
	}

	switch{
	case player.WeaponType != request.Weapon && player.IsAvailable(request.Weapon):
		player.ChangeWeapon(request.Weapon)
	case request.Reload() && player.Reload():
		player.SetAction(RELOAD)
	case request.Shoot() && !player.EmptyMag():
		player.Shoot(g)
	case request.Shoot() && player.Reload(): // Has no ammo
		player.SetAction(RELOAD)
	case request.Melee():
		player.SetAction(MELEE)
		// todo: create melee attack
	default:
		if request.Moved() {
			player.Action = MOVE
		} else {
			player.Action = IDLE
		}
	}
}

func (player Player) GetHitbox() float64{
	return 50
}

func (player Player) Turn() {
	turnDelays[player.Id] = Timestamp + player.TurnSpeed()
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

func (player Player) ID() string {
	return player.Id
}

func (player Player) EntityType() EntityType {
	return PlayerE
}

func (player Player) GetPos() pixel.Vec{
	return player.Pos
}

func (player Player) GetDir() float64{
	return player.Dir
}