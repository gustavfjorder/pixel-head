package server

import (
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
)

type Game struct {
	space        Space
	memory       Memory
	clientSpaces []Space
}

func NewGame(uri string, clientUris []string) Game {
	clientSpaces := make([]Space, len(clientUris))
	for i, clientUri := range clientUris {
		clientSpaces[i] = NewSpace(clientUri)
		clientSpaces[i].Put("done")
	}

	return Game{
		space:        setupSpace(uri),
		memory:       NewMemory(),
		clientSpaces: clientSpaces,
	}
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
	fmt.Println("Starting game")
	//fmt.Println("Starting game on uri '" + uri + "'")
	//fmt.Println("Players in game:", playerIds)

	for _, space := range g.clientSpaces {
		space.Put("map", model.MapTemplates["Test1"])
	}

	t := time.Tick(time.Second / 30)

	fmt.Println("Starting game loop")
	for {
		//g.space.Get("loop_lock")

		g.handleRequests()
		g.handleZombies()

		for _, space := range g.clientSpaces {
			if _, err := space.GetP("done"); err == nil {
				go g.putToSpaces(&space)
			}
		}

		//g.space.Put("loop_lock")
		<- t
	}
}

func (g *Game) putToSpaces(space *Space) {
	//space.Get("done")

	fmt.Println("Putting to client space")

	space.Put("players", g.memory.GetW("players").([]model.Player))
	space.Put("zombies", g.memory.GetW("zombies",  make([]model.Zombie, 0)).([]model.Zombie))
	space.Put("shoots", g.memory.GetW("shots", make([]model.Shoot, 0)).([]model.Shoot))

	space.Put("ready")
}

func setupSpace(uri string) Space {
	// Register models for encoding to space
	gob.Register(model.Request{})
	gob.Register(model.Player{})
	gob.Register([]model.Player{})
	gob.Register(model.Zombie{})
	gob.Register([]model.Zombie{})
	gob.Register(model.Shoot{})
	gob.Register([]model.Shoot{})
	gob.Register(model.Map{})
	gob.Register(model.Wall{})
	gob.Register(model.Line{})
	gob.Register(model.Point{})

	space := NewSpace(uri)

	// todo: pSpaces seems to need this to be able to Get/Query on clients
	space.QueryP(&model.Request{})
	space.QueryP(&model.Player{})
	space.QueryP(&[]model.Player{})
	space.QueryP(&model.Zombie{})
	space.QueryP(&[]model.Zombie{})
	space.QueryP(&model.Shoot{})
	space.QueryP(&[]model.Shoot{})
	space.QueryP(&model.Map{})
	space.QueryP(&model.Wall{})
	space.QueryP(&model.Line{})
	space.QueryP(&model.Point{})

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
		player.Weapon = request.CurrentWep

		if request.Move {
			// todo: check if move is doable in map
			player = player.Move(request.Dir)
		}

		if request.Reload {
			player.Reload = true
			player.GetWeapon().RefillMag()
		} else if request.Shoot {
			player.Shoot = true
			playerShoots := player.GetWeapon().GenerateShoots(request.Timestamp, player.Pos)
			g.memory.PutToArray("shots", playerShoots)
		} else if request.Melee {
			player.Melee = true
			// todo: create melee attack
		}

		//g.memory.Update("player." + request.PlayerId, player)

		players = append(players, player)
	}

	g.memory.Update("players", players)
}

func (g *Game) handleZombies() {
	zombies := g.memory.GetW("zombies", make([]model.Zombie, 0)).([]model.Zombie)

	shots := g.memory.GetW("shots", make([]model.Shoot, 0)).([]model.Shoot)

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
