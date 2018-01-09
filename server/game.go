package server

import (
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"strconv"
)

type Game struct {
	space        Space
	clientSpaces []Space
	state        *model.State
	currentMap   model.Map
	currentLevel int
}

func NewGame(uri string, clientUris []string) Game {
	clientSpaces := make([]Space, len(clientUris))
	for i, clientUri := range clientUris {
		clientSpaces[i] = NewSpace(clientUri)
		clientSpaces[i].Put("done")
	}

	return Game{
		space:        setupSpace(uri),
		clientSpaces: clientSpaces,
		state:        &model.State{},
	}
}

func (g *Game) AddPlayers(playerIds []string) {
	for _, id := range playerIds {
		g.AddPlayer(id)
	}
}

func (g *Game) AddPlayer(id string) {
	g.state.Players = append(g.state.Players, model.NewPlayer(id))
}

func (g *Game) Start() {
	fmt.Println("Starting game")

	g.currentMap = model.MapTemplates["Test1"]

	for _, space := range g.clientSpaces {
		space.Put("map", g.currentMap)
		space.Put("ready")
	}

	time.Sleep(time.Second * 2)

	fmt.Println("Starting game loop")
	t := time.Tick(time.Second / 30)
	for g.currentLevel < len(model.Levels) {
		fmt.Println("Starting level " + strconv.Itoa(g.currentLevel))

		levelPrepared := make(chan bool)

		duration := time.Second * 10
		if g.currentLevel == 0 {
			duration = 0
		}
		time.AfterFunc(duration, func() {
			g.prepareLevel(levelPrepared)
		})

		fmt.Println("after prepare")

		breakable := false

		for {
			fmt.Println("Game looooooping")
			select {
			case <- levelPrepared:
				breakable = true
			default:
			}

			g.handleRequests()
			g.handleZombies()

			for _, space := range g.clientSpaces {
				if _, err := space.GetP("done"); err == nil {
					g.putToSpaces(&space)
				}
			}

			if breakable && len(g.state.Zombies) == 0 {
				break
			}

			<- t
		}

		g.currentLevel++

		<- t
	}
}

func (g *Game) prepareLevel(done chan bool) {
	level := model.Levels[g.currentLevel]

	for i := 0; i < level.NumberOfZombies; i++ {
		g.state.Zombies = append(g.state.Zombies, model.NewZombie(300, 200))
	}

	close(done)
}

func (g *Game) putToSpaces(space *Space) {
	//space.Get("done")

	fmt.Println("Putting to client space")

	space.Put("state", g.state)

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
	gob.Register(model.State{})

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
	space.QueryP(&model.State{})

	return space
}

func (g *Game) handleRequests() {
	// Load incoming requests
	rTuples, _ := g.space.GetAll(&model.Request{})

	players := g.state.Players
	for _, rTuple := range rTuples {
		request := rTuple.GetFieldAt(0).(model.Request)

		// Load player
		var player *model.Player
		for i, p := range players {
			if p.Id == request.PlayerId {
				player = &(players)[i]
				break
			}
		}

		// Change weapon
		player.ChangeWeapon(request.CurrentWep)

		if request.Move {
			fmt.Println("Handling move request from:", player.Pos)
			// todo: check if move is doable in map
			player.Move(request.Dir)
			fmt.Println("Moved to:", player.Pos)
		}

		player.Reload = false
		player.Shoot = false
		player.Melee = false

		if request.Reload {
			player.Reload = true
			player.GetWeapon().RefillMag()
		} else if request.Shoot {
			player.Shoot = true
			playerShoots := player.GetWeapon().GenerateShoots(request.Timestamp, player.Pos)
			g.state.Shoots = append(g.state.Shoots, playerShoots...)
		} else if request.Melee {
			player.Melee = true
			// todo: create melee attack
		}
	}
}

func (g *Game) handleZombies() {
	zombies := &g.state.Zombies
	shots := &g.state.Shoots

	for i, zombie := range *zombies {
		// Any shoots hitting the zombie
		for i, shoot := range *shots {
			if shoot.GetPos() == zombie.Pos {
				zombie.Stats.Health -= shoot.Weapon.Power
				*shots = append((*shots)[:i], (*shots)[i + 1:]...)
			}
		}

		if zombie.Stats.Health <= 0 {
			*zombies = append((*zombies)[:i], (*zombies)[i + 1:]...)
			continue
		}

		zombie.Move(&g.state.Players)
		zombie.Attack(&g.state.Players)
	}
}
