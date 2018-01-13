package model

import (
	"math/rand"
	"fmt"
	"github.com/faiface/pixel"
)

type Game struct {
	PlayerIds    map[string]bool // Is true if player is active in game
	State        State
	Updates      Updates
	CurrentMap   Map
	CurrentLevel int
}

func NewGame(ids []string, mapName string) (game Game) {
	game.PlayerIds = make(map[string]bool)
	game.State.Players = make([]Player, len(ids))
	game.CurrentLevel = 0
	game.CurrentMap = MapTemplates[mapName]
	game.State.Barrels = make([]Barrel, 1)
	game.State.Barrels[0] = NewBarrel(pixel.V(500,500))
	for i, id := range ids {
		game.State.Players[i] = NewPlayer(id)
		game.PlayerIds[id] = true
	}
	return game
}

func (game *Game) PrepareLevel(end chan<- bool) {
	level := Levels[game.CurrentLevel]
	game.State.Zombies = make([]Zombie, level.NumberOfZombies)
	for i := range game.State.Zombies {
		game.State.Zombies[i] = NewZombie(rand.Float64()*900+100, rand.Float64()*900+100)
	}
	end <- true
}

func (game *Game) HandleRequests(requests []Request) {
	// Load incoming requests
	for _, request := range requests {
		// Load player
		player, err := findPlayer(game.State.Players, request.PlayerId)

		if err != nil {
			fmt.Println(err)
			continue
		}


		if request.Moved() {
			player.Move(request.Dir, game)
		}

		if Timestamp >= player.ActionDelay() {
			player.Action = IDLE
		}

		//Action priority is like so: weapon change > reload > shoot > melee
		player.Do(request, game)
	}
}

func (game *Game) HandleZombies() {
	for i := len(game.State.Zombies) - 1; i >= 0; i-- {
		zombie := &game.State.Zombies[i]

		// Any shoots hitting the zombie
		for j := len(game.State.Shots) - 1; j >= 0; j-- {
			shoot := game.State.Shots[j]
			if shoot.GetPos().Sub(zombie.Pos).Len() <= zombie.GetHitbox() {
				zombie.Stats.Health -= shoot.WeaponType.Power()
				game.Remove(Entry{shoot, j})
			}
		}

		//Remove all zombies at zero health
		if zombie.Stats.Health <= 0 {
			game.Remove(Entry{*zombie, i})
			continue
		}

		zombie.Move(game.State.Players)
		zombie.Attack(game.State)
	}
}

func (game *Game) HandleShots() {
	for i := len(game.State.Shots) - 1; i >= 0; i-- {
		shot := game.State.Shots[i]
		if shot.GetPos().Sub(shot.Start).Len() > shot.WeaponType.Range() {
			game.Remove(Entry{shot, i})
			continue
		}
	}
}

func (game *Game) HandlePlayers() {
	for i := len(game.State.Players) - 1; i >= 0; i-- {
		player := game.State.Players[i]
		if player.Stats.Health <= 0 {
			//Remove player from game
			game.PlayerIds[player.Id] = false
			game.Remove(Entry{player, i})
		}
	}
}

func (game *Game) HandleBarrels() {
	for i := len(game.State.Barrels) - 1; i >= 0; i-- {
		barrel := game.State.Barrels[i]
		for j := len(game.State.Shots) - 1; j >= 0; j-- {
			shot := game.State.Shots[j]
			if shot.GetPos().Sub(barrel.Pos).Len() < barrel.GetHitBox() {
				//Update objects
				barrel.Explode(&game.State)
				shot.Hit = true

				//Add to updates and remove from state
				game.Remove(Entry{barrel, i}, Entry{shot, j})
				break
			}
		}
	}
}

func (game *Game) Remove(entries ...Entry){
	for _, entry := range entries {
		switch entry.elem.(type){
		case Player:
			last := len(game.State.Players) - 1
			game.State.Players[entry.index] = game.State.Players[last]
			game.State.Players = game.State.Players[:last]
			game.Updates.Remove(entry.elem.(Player))
		case Shot:
			last := len(game.State.Shots) - 1
			game.State.Shots[entry.index] = game.State.Shots[last]
			game.State.Shots = game.State.Shots[:last]
			game.Updates.Remove(entry.elem.(Shot))
		case Zombie:
			last := len(game.State.Zombies) - 1
			game.State.Zombies[entry.index] = game.State.Zombies[last]
			game.State.Zombies = game.State.Zombies[:last]
			game.Updates.Remove(entry.elem.(Zombie))
		case Barrel:
			last := len(game.State.Barrels) - 1
			game.State.Barrels[entry.index] = game.State.Barrels[last]
			game.State.Barrels = game.State.Barrels[:last]
			game.Updates.Remove(entry.elem.(Barrel))
		}
	}
}
