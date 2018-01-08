package server

import (
	"github.com/faiface/pixel"
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
)

type Game struct {
	space  Space
	memory Memory
}

func NewGame(uri string) Game {
	return Game{setupSpace(uri), NewMemory()}
}

func (g *Game) AddPlayers(playerIds []string) {
	for _, id := range playerIds {
		g.AddPlayer(id)
	}
}

func (g *Game) AddPlayer(id string) {
	players := g.memory.GetW("players", make([]model.Player, 0)).([]model.Player)

	players = append(players, model.NewPlayer(id))

	g.memory.Update("players", players)
}

func (g *Game) Start() {
	//fmt.Println("Starting game on uri '" + uri + "'")
	//fmt.Println("Players in game:", playerIds)

	t := time.Tick(time.Second / 60)

	for {
		g.space.Get("loop_lock")

		g.handleRequests()
		g.handleZombies()

		for _, player := range g.memory.GetW("players").([]model.Player) {
			g.space.Put(player)
		}

		for _, zombie := range g.memory.GetW("zombies").([]model.Zombie) {
			g.space.Put(zombie)
		}

		for _, shoot := range g.memory.GetW("shots").([]model.Shoot) {
			g.space.Put(shoot)
		}

		g.space.Put("loop_lock")
		<- t
	}
}

func setupSpace(uri string) Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register(model.Zombie{})
	gob.Register(model.Shoot{})

	space := NewSpace(uri)

	space.Put("loop_lock")

	return space
}

func (g *Game) handleRequests() {
	// Load incoming requests
	rTuples, _ := g.space.GetAll(&model.Request{})

	players := g.memory.GetW("players").([]model.Player)

	for _, rTuple := range rTuples {
		request := rTuple.GetFieldAt(0).(model.Request)
		fmt.Println("Handling request:", request)

		// Load player
		//player := g.memory.GetW("player." + request.PlayerId).(model.Player)
		var player model.Player
		for i, p := range players {
			if p.Id == request.PlayerId {
				player = players[i]
				players = append(players[:i], players[i + 1:]...)
				break
			}
		}

		// Change weapon
		player.Weapon = model.Weapons[request.CurrentWep]

		if request.Move {
			// todo: check if move is doable in map
			player = player.Move(request.Dir)
		}

		if request.Reload {
			player.Weapon.RefillMag()
		} else if request.Shoot {
			playerShoots := player.Weapon.GenerateShoots(request.Timestamp, player.Pos)
			g.memory.PutToArray("shots", playerShoots)
		}

		//g.memory.Update("player." + request.PlayerId, player)

		players = append(players, player)
	}

	g.memory.Update("players", players)
}

func (g *Game) handleZombies() {
	mZombies, _ := g.memory.Get("zombies")
	zombies := mZombies.([]model.Zombie)

	mShots, _ := g.memory.Get("shots")
	shots := mShots.([]model.Shoot)

	for i, zombie := range zombies {
		// Any shoots hitting the zombie
		for i, shoot := range shots {
			if shoot.GetPos() == zombie.Pos {
				zombie.Stats.Health -= shoot.Weapon.Power
				shots = append(shots[:i], shots[i + 1:]...)
			}
		}

		if zombie.Stats.Health <= 0 {
			zombies = append(zombies[:i], zombies[i + 1:]...)
			continue
		}

		players := g.memory.GetW("players").([]model.Player)

		zombie.Move(players)
		zombie.Attack(players)
	}

	g.memory.Update("zombies", zombies)
	g.memory.Update("shots", shots)
}

func (g Game) players() []model.Player {
	players, _ := g.memory.Get("players")

	return players.([]model.Player)
}
