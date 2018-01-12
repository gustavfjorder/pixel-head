package model

import (
	"math/rand"
	"github.com/faiface/pixel"
)

type Game struct {
	PlayerIds    []string
	State        State
	CurrentMap   Map
	CurrentLevel int
}

func NewGame(ids []string, mapName string) (game Game) {
	game.PlayerIds = ids
	game.State.Players = make([]Player, len(ids))
	game.CurrentLevel = 0
	game.CurrentMap = MapTemplates[mapName]
	for i, id := range ids {
		game.State.Players[i] = NewPlayer(id)
	}
	game.State.Barrels=[]Barrel {NewBarrel("1",pixel.Vec{500,500})}
	return game
}

func (g *Game) PrepareLevel(end <-chan bool) {
	level := Levels[g.CurrentLevel]
	g.State.Zombies = make([]Zombie, level.NumberOfZombies)
	for i := range g.State.Zombies {
		g.State.Zombies[i] = NewZombie(rand.Float64()*900+100, rand.Float64()*900+100)
	}
	<-end
}

func (g *Game) HandleRequests(requests []Request) {
	// Load incoming requests

	players := g.State.Players
	for _, request := range requests {

		// Load player
		var player *Player
		for i, p := range players {
			if p.Id == request.PlayerId {
				player = &(players)[i]
				break
			}
		}
		if g.State.Timestamp >= player.ActionDelay {
			player.Reload = false
			player.Shoot = false
			player.Melee = false
		}

		if request.Move {
			// todo: check if move is doable in map
			player.Move(request.Dir, g)
		}

		//Action priority is like so: weapon change > reload > shoot > melee
		switch {
		case g.State.Timestamp < player.ActionDelay:
			break
		case player.GetWeapon().Id != request.Weapon && player.IsAvailable(request.Weapon):
			player.ChangeWeapon(request.Weapon)
		case request.Reload && player.GetWeapon().RefillMag():
			player.Reload = true
			player.ActionDelay = player.GetWeapon().GetReloadSpeed() + g.State.Timestamp
		case request.Shoot && player.GetWeapon().MagazineCurrent > 0:
			playerShoots := player.GetWeapon().GenerateShoots(g.State.Timestamp, *player)
			player.Shoot = len(playerShoots) > 0
			g.State.Shots = append(g.State.Shots, playerShoots...)
			player.ActionDelay = player.GetWeapon().GetShootDelay() + g.State.Timestamp
		case request.Shoot && player.GetWeapon().RefillMag(): // Has no ammo
			player.Reload = true
			player.ActionDelay = player.GetWeapon().GetReloadSpeed() + g.State.Timestamp
		case request.Melee:
			player.Melee = true
			// todo: create melee attack
		}
	}
}

func (g *Game) HandleZombies() {
	for i := len(g.State.Zombies) - 1; i >= 0; i-- {
		zombie := &g.State.Zombies[i]

		// Any shoots hitting the zombie
		for j := len(g.State.Shots) - 1; j >= 0; j-- {
			shoot := g.State.Shots[j]
			if shoot.GetPos(g.State.Timestamp).Sub(zombie.Pos).Len() <= zombie.GetHitbox() {
				zombie.Stats.Health -= GetWeaponRef(shoot.Weapon).GetPower()
				g.State.Shots[j] = g.State.Shots[len(g.State.Shots)-1]
				g.State.Shots = g.State.Shots[:len(g.State.Shots)-1]
			}
		}

		//Remove all zombies at zero health
		if zombie.Stats.Health <= 0 {
			g.State.Zombies[i] = g.State.Zombies[len(g.State.Zombies)-1]
			g.State.Zombies = g.State.Zombies[:len(g.State.Zombies)-1]
			continue
		}

		zombie.Move(g.State.Players)
		zombie.Attack(g.State)
	}
}

func (g *Game) HandleShots() {
	for i := len(g.State.Shots) - 1; i >= 0; i-- {
		shot := g.State.Shots[i]
		if shot.GetPos(g.State.Timestamp).Sub(shot.Start).Len() > GetWeaponRef(shot.Weapon).GetRange() {
			g.State.Shots[i] = g.State.Shots[len(g.State.Shots)-1]
			g.State.Shots = g.State.Shots[:len(g.State.Shots)-1]
			continue
		}
	}
}

func (g *Game) HandlePlayers() {
	for i, player := range g.State.Players {
		if player.Stats.Health <= 0 {
			g.State.Players = append(g.State.Players[:i], g.State.Players[i+1:]...)
		}
	}
}

func (g *Game) HandleBarrels(){
	for i:=len(g.State.Barrels)-1; i>=0;i--{
		barrel:=g.State.Barrels[i]
		for j:=len(g.State.Shots)-1;j>=0;j--{
			shot:=g.State.Shots[j]
			if shot.GetPos(g.State.Timestamp).Sub(barrel.Pos).Len()<barrel.GetHitBox(){
				barrel.Explode(&g.State)

				g.State.Barrels[i]=g.State.Barrels[len(g.State.Barrels)-1]
				g.State.Barrels=g.State.Barrels[:len(g.State.Barrels)-1]

				g.State.Shots[j]=g.State.Shots[len(g.State.Shots)-1]
				g.State.Shots=g.State.Shots[:len(g.State.Shots)-1]

				break
			}
		}
	}
}