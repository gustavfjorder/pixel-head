package model

import (
	"fmt"
	"github.com/faiface/pixel"
	"math/rand"
	"time"
	/*"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
*/
	)

type Game struct {
	PlayerIds    map[string]bool // Is true if player is active in game
	State        State
	CurrentMap   Map
	CurrentLevel int
}

func NewGame(ids []string, mapName string) (game Game) {
	game.PlayerIds = make(map[string]bool)
	game.State.Players = make([]Player, len(ids))
	game.CurrentLevel = 0
	game.CurrentMap = MapTemplates[mapName]
	for i, id := range ids {
		game.State.Players[i] = NewPlayer(id)
		game.PlayerIds[id] = true
	}
	game.State.Barrels=[]Barrel {NewBarrel("1",pixel.Vec{500,500}),
		NewBarrel("1",pixel.Vec{700,500}),
		NewBarrel("1",pixel.Vec{400,500})}
	return game
}

func (g *Game) PrepareLevel(end chan<- bool) {
	level := Levels[g.CurrentLevel]
	g.State.Zombies = make([]Zombie, 0)
	waveticker := time.NewTicker(level.TimeBetweenWaves)
	zombieticker:=time.NewTicker(level.TimeBetweenZombies)
	fmt.Println("num zom",level.NumberOfZombiesPerWave,"num waves",level.NumberOfWaves)

	for i:=0;i<level.NumberOfWaves;i++ {
		fmt.Println("i:",i)
		<-waveticker.C
		for j:=0;j<level.NumberOfZombiesPerWave;j++ {
			fmt.Println("j:",j)
			fmt.Println("i*level.NumberOfZombiesPerWave+j",i*level.NumberOfZombiesPerWave+j)
			fmt.Println(len(g.State.Zombies))
			g.State.Zombies = append(g.State.Zombies,NewZombie(g.CurrentMap.SpawnPoint[rand.Intn(len(g.CurrentMap.SpawnPoint))]))
			<-zombieticker.C
			}


	}
	//todo print level complete

	end<-true
}
/*
func (g *Game) LevelDone(level int){
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	levelText := text.New(pixel.V(g.CurrentMap.SpawnPoint[0].X,g.CurrentMap.SpawnPoint[0].X), basicAtlas)
	fmt.Fprintln(levelText, "Level",level, "complete")
	levelText.Draw(,pixel.IM.Moved(pixel.Vec{500,500}))
}
*/

func (g *Game) HandleRequests(requests []Request) {
	// Load incoming requests
	for _, request := range requests {
		timestamp := g.State.Timestamp
		// Load player
		player, err := findPlayer(g.State.Players, request.PlayerId)

		if err != nil {
			fmt.Println(err)
			continue
		}

		weapon, err := player.GetWeapon()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if request.Moved() {
			player.Move(request.Dir, g)
		}

		if timestamp >= player.ActionDelay() {
			player.Action = IDLE
		}

		//Action priority is like so: weapon change > reload > shoot > melee
		switch {
		case timestamp < player.ActionDelay():
			break
		case weapon.WeaponType != request.Weapon && player.IsAvailable(request.Weapon):
			player.ChangeWeapon(request.Weapon)
		case request.Reload() && weapon.RefillMag():
			player.SetAction(RELOAD,timestamp)
		case request.Shoot() && weapon.MagazineCurrent > 0:
			playerShoots := weapon.GenerateShoots(g.State.Timestamp, *player)
			g.State.Shots = append(g.State.Shots, playerShoots...)
			player.SetAction(SHOOT, timestamp)
		case request.Shoot() && weapon.RefillMag(): // Has no ammo
			player.SetAction(RELOAD, timestamp)
		case request.Melee():
			player.SetAction(MELEE, timestamp)
			// todo: create melee attack
		default:
			if request.Moved(){
				player.Action = MOVE
			} else {
				player.Action = IDLE
			}
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
				zombie.Stats.Health -= shoot.WeaponType.Power()
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
		if shot.GetPos(g.State.Timestamp).Sub(shot.Start).Len() > shot.WeaponType.Range() {
			g.State.Shots[i] = g.State.Shots[len(g.State.Shots)-1]
			g.State.Shots = g.State.Shots[:len(g.State.Shots)-1]
			continue
		}
	}
}

func (g *Game) HandlePlayers() {
	for i := len(g.State.Players) -1 ; i >= 0; i-- {
		player := g.State.Players[i]
		if player.Stats.Health <= 0 {
			for i := 0; i < len(g.PlayerIds); i++ {
				//Remove player from game
				g.PlayerIds[player.Id] = false
			}
			g.State.Players[i] = g.State.Players[len(g.State.Players)-1]
			g.State.Players = g.State.Players[:len(g.State.Players)-1]
		}
	}
}

func (g *Game) HandleBarrels(){
	for i:=len(g.State.Barrels)-1; i>=0;i--{
		barrel:=g.State.Barrels[i]
		for j:=len(g.State.Shots)-1;j>=0;j--{
			shot:=g.State.Shots[j]
			if shot.GetPos(g.State.Timestamp).Sub(barrel.Pos).Len()<barrel.GetHitBox(){
				barrel.Explode(i,&g.State)
				b :=Barrel{}
				if g.State.Barrels[i] != b {
					g.State.Barrels[i] = g.State.Barrels[len(g.State.Barrels)-1]
					g.State.Barrels = g.State.Barrels[:len(g.State.Barrels)-1]
					fmt.Println("barrel removed in handleBarrels")
					g.State.Shots[j] = g.State.Shots[len(g.State.Shots)-1]
					g.State.Shots = g.State.Shots[:len(g.State.Shots)-1]
				}
				break
			}
		}
	}
}