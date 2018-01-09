package server

import (
	. "github.com/pspaces/gospace"
	"time"
	"encoding/gob"
	"github.com/gustavfjorder/pixel-head/model"
	"fmt"
	"strconv"
	"math/rand"
)

type ClientSpace struct {
	Id string
	Space
}

type Game struct {
	space        Space
	clientSpaces []ClientSpace
	state        *model.State
	currentMap   model.Map
	currentLevel int
}

func NewGame(uri string, clientUris, ids []string) Game {
	clientSpaces := make([]ClientSpace, len(clientUris))
	for i, clientUri := range clientUris {
		clientSpaces[i] = ClientSpace{
			Id:    ids[i],
			Space: NewSpace(clientUri),
		}
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
	}

	for _, space := range g.clientSpaces {
		space.Get("joined")
	}

	fmt.Println("Starting game loop")
	t := time.Tick(time.Second / 30)
	for g.currentLevel < len(model.Levels) {
		fmt.Println("Starting level " + strconv.Itoa(g.currentLevel))

		levelPrepared := make(chan bool)

		duration := time.Second * 1
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
			case <-levelPrepared:
				breakable = true
			default:
			}

			g.state.Timestamp = time.Now().UnixNano()

			g.handleRequests()
			g.handleZombies()
			g.handleShots()

			for i, player := range g.state.Players {
				fmt.Println(player.Stats.Health)
				if player.Stats.Health <= 0 {
					g.state.Players = append(g.state.Players[:i], g.state.Players[i+1:]...)
				}
			}

			for _, space := range g.clientSpaces {
				space.GetP("state", &model.State{})
				space.Put("state", g.state)
			}

			if breakable && len(g.state.Zombies) == 0 || len(g.state.Players) == 0 {
				break
			}

			<-t
		}

		if len(g.state.Players) == 0 {
			break
		}

		g.currentLevel++

	}
}

func (g *Game) prepareLevel(done chan bool) {
	level := model.Levels[g.currentLevel]

	for i := 0; i < level.NumberOfZombies; i++ {
		g.state.Zombies = append(g.state.Zombies, model.NewZombie(300+(rand.Float64()*100), 200+(rand.Float64()*100)))
	}

	close(done)
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
	requests:= make([]model.Request, len(g.clientSpaces))
	for i, space := range g.clientSpaces {
		rtuples, _ := space.GetAll(&model.Request{})
		for _, rtuple := range rtuples {
			request := rtuple.GetFieldAt(0).(model.Request)
			requests[i] = requests[i].Merge(request)
		}
		requests[i].PlayerId = space.Id
	}
	players := g.state.Players
	for _, request := range requests {

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
			playerShoots := player.GetWeapon().GenerateShoots(g.state.Timestamp, *player)
			g.state.Shoots = append(g.state.Shoots, playerShoots...)
		} else if request.Melee {
			player.Melee = true
			// todo: create melee attack
		}
	}
}

func (g *Game) handleZombies() {
	for i := len(g.state.Zombies) - 1; i >= 0; i-- {
		zombie := &g.state.Zombies[i]
		// Any shoots hitting the zombie
		for j := len(g.state.Shoots) - 1; j >= 0; j-- {
			shoot := g.state.Shoots[j]
			if shoot.GetPos(g.state.Timestamp).Sub(zombie.Pos).Len() <= zombie.GetHitbox() {
				zombie.Stats.Health -= model.Weapons[shoot.Weapon].Power
				g.state.Shoots[j] = g.state.Shoots[len(g.state.Shoots)-1]
				g.state.Shoots = g.state.Shoots[:len(g.state.Shoots)-1]
			}
		}

		if zombie.Stats.Health <= 0 {
			g.state.Zombies[i] = g.state.Zombies[len(g.state.Zombies)-1]
			g.state.Zombies = g.state.Zombies[:len(g.state.Zombies)-1]
			continue
		}

		zombie.Move(g.state.Players)
		zombie.Attack(g.state.Players)
	}
}

func (g *Game) handleShots() {
	for i := len(g.state.Shoots) - 1; i >= 0; i-- {
		shot := g.state.Shoots[i]
		fmt.Println(shot.GetPos(g.state.Timestamp).Sub(shot.Start).Len())
		if shot.GetPos(g.state.Timestamp).Sub(shot.Start).Len() > model.Weapons[shot.Weapon].Range {
			g.state.Shoots[i] = g.state.Shoots[len(g.state.Shoots)-1]
			g.state.Shoots = g.state.Shoots[:len(g.state.Shoots)-1]
			continue
		}
	}
}
