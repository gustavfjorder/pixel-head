package model

import (
	"github.com/faiface/pixel"
	"math"
	"time"
	"github.com/pkg/errors"
)

type Action int

const (
	IDLE Action = iota
	SHOOT
	MELEE
	MOVE
	RELOAD
	BARREL
)

type Player struct {
	Id         string
	Pos        pixel.Vec
	Dir        float64
	WeaponType WeaponType
	WeaponList []WeaponI
	Action     Action
	Stats
}

func NewPlayer(id string, pos pixel.Vec) (player Player) {
	player.WeaponList = make([]WeaponI, nWeapon)
	player.WeaponType = KNIFE
	player.NewWeapon(NewWeapon(player.WeaponType), NewWeapon(RIFLE), NewWeapon(SHOTGUN))
	player.Id = id
	player.Stats = NewStats(HUMAN)
	player.Dir = 0
	player.Pos = pos

	return
}

func (player *Player) Move(dir float64, g *Game) {
	if !math.IsNaN(dir) {
		if player.TurnDelay(g) < g.State.Timestamp {
			player.Dir = dir
			player.Turn(g)
			return
		}
		newpos := player.Pos.Add(pixel.V(player.Stats.GetMoveSpeed(), 0).Rotated(player.Dir))
		if !g.CurrentMap.Bounds.Contains(newpos){
			return
		}
		for _, wall := range g.CurrentMap.Walls {
			if wall.Intersect(NewLine(PointFrom(player.Pos), PointFrom(newpos))) {
				return
			}
		}

		player.Pos = newpos
	}
}

func (player *Player) NewWeapon(weapons ...WeaponI) {
	for _, weapon := range weapons {
		player.WeaponList[weapon.Type()] = weapon
	}

}

func (player *Player) Weapon() (weapon WeaponI, e error) {
	if player.WeaponType >= nWeapon || player.WeaponType < 0 {
		weapon = &WeaponBase{}
		e = errors.New("Unable to change weapon")
		return
	}
	weapon = player.WeaponList[player.WeaponType]
	return weapon, e

}

func (player *Player) ChangeWeapon(weaponType WeaponType) {
	if player.IsAvailable(weaponType) {
		player.WeaponType = weaponType
	}
}

func (player *Player) IsAvailable(weaponType WeaponType) bool {
	return player.WeaponList[weaponType] != nil
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

func (player *Player) SetAction(action Action, game *Game){
	player.Action = action
	game.actionDelays[player.Id] = player.Delay(action) + game.State.Timestamp
}

func (player Player) Delay(action Action) time.Duration{
	wep, err := player.Weapon()
	if err != nil{
		return 0
	}
	switch action {
	case RELOAD:
		return wep.ReloadSpeed()
	case SHOOT, MELEE:
		return wep.ShootDelay()
	default:
		return 0
	}
}

func (player Player) Reload(game *Game) bool{
	wep, err := player.Weapon()
	if err != nil {
		return false
	}
	player.SetAction(RELOAD, game)
	return wep.RefillMag()
}

func (player Player) EmptyMag() bool{
	wep, err := player.Weapon()
	if err != nil {
		return true
	}
	return wep.GetMagazine() <= 0
}

func (player *Player) Shoot(g *Game) {
	weapon, err := player.Weapon()
	if err != nil {
		return
	}
	weapon.Shoot(*player, g)
	player.SetAction(SHOOT, g)
}

func (player *Player) Do(request Request, g *Game) {
	if g.State.Timestamp < player.ActionDelay(g){
		return
	}

	switch{
	case player.WeaponType != request.Weapon && player.IsAvailable(request.Weapon):
		player.ChangeWeapon(request.Weapon)
	case request.Reload() && player.Reload(g):
		player.SetAction(RELOAD, g)
	case request.Shoot() && (!player.EmptyMag() || player.WeaponType == KNIFE):
		player.Shoot(g)
	case request.Shoot() && player.Reload(g): // Has no ammo
		player.SetAction(RELOAD, g)
	case request.Melee():
		player.SetAction(MELEE, g)
		// todo: create melee attack
	case request.Barrel():
		player.SetAction(BARREL, g)
		g.Add(NewBarrel(player.Pos))
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

func (player Player) Turn(game *Game) {
	game.turnDelays[player.Id] = game.State.Timestamp + player.TurnSpeed()
}

func (player Player) ActionDelay(game *Game) time.Duration {
	return game.actionDelays[player.Id]
}

func (player Player) TurnDelay(game *Game) time.Duration {
	return game.turnDelays[player.Id]
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

func (player *Player) Regen(){
	player.Health+= MinInt(player.GetMaxHealth() - player.Health, 3)
}