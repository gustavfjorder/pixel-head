package model

import (
	"math/rand"
	"fmt"
	"github.com/faiface/pixel"
	"time"
	"sort"
	"github.com/gustavfjorder/pixel-head/helper"
	"github.com/rs/xid"
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
	game.Add(NewBarrel(pixel.V(500,500)), NewBarrel(pixel.V(600,600)), NewBarrel(pixel.V(700,700)), NewBarrel(pixel.V(900,900)), NewBarrel(pixel.V(1000,1000)))
	for _, id := range ids {
		game.Add( NewPlayer(id))
		game.PlayerIds[id] = true
	}
	return
}

func (g *Game) PrepareLevel(end chan<- bool) {
	level := Levels[g.CurrentLevel]
	g.State.Zombies = make([]ZombieI, 0)
	waveticker := time.NewTicker(level.TimeBetweenWaves)
	zombieticker := time.NewTicker(level.TimeBetweenZombies)
	fmt.Println("num zom", level.NumberOfZombiesPerWave, "num waves", level.NumberOfWaves)

	for i := 0; i < level.NumberOfWaves; i++ {
		fmt.Println("i:", i)
		<-waveticker.C
		for j := 0; j < level.NumberOfZombiesPerWave; j++ {
			fmt.Println("j:", j)
			fmt.Println("i*level.NumberOfZombiesPerWave+j", i*level.NumberOfZombiesPerWave+j)
			fmt.Println(len(g.State.Zombies))
			var ZOM Being
			switch rand.Intn(4){
			case 0:
				ZOM = ZOMBIE
			case 1:
				ZOM = BOMBZOMBIE
			case 2:
				ZOM = FASTZOMBIE
			case 3:
				ZOM = SLOWZOMBIE
			}
			g.NewZombie(g.CurrentMap.SpawnPoint[rand.Intn(len(g.CurrentMap.SpawnPoint))],ZOM)
			<-zombieticker.C
		}
	}
	end<-true
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

var lastTime = time.Now()
func (g *Game) HandleLoot() {
	for i := len(g.State.Players) - 1; i >= 0; i-- {
		player := &g.State.Players[i]

		for j := len(g.State.Lootboxes) - 1; j >= 0; j-- {
			lootbox := g.State.Lootboxes[j]

			fmt.Println("Distance", PointFrom(player.Pos).Dist(PointFrom(lootbox.Pos)))
			if PointFrom(player.Pos).Dist(PointFrom(lootbox.Pos)) < 30 {
				player.PickupLootbox(&lootbox)
				g.Remove(Entry{lootbox, j})
			}
		}
	}

	// Place lootboxes
	if lastTime.Add(time.Second * 10).Before(time.Now()) && float64(rand.Intn(100)) <= 3.8 {
		min := 0
		max := len(g.CurrentMap.LootPoints)

		lootPoint := rand.Intn(max - min) + min
		point := g.CurrentMap.LootPoints[lootPoint]

		if ! g.State.HasLootboxAt(point) {
			g.Add(NewLootbox(point.X, point.Y))

			lastTime = time.Now()
		}
	}
}

func (game *Game) HandleZombies() {
	for i := len(game.State.Zombies) - 1; i >= 0; i-- {
		zombie := game.State.Zombies[i]

		// Any shoots hitting the zombie
		for j := len(game.State.Shots) - 1; j >= 0; j-- {
			shoot := game.State.Shots[j]

			if shoot.GetPos().Sub(zombie.GetPos()).Len() <= zombie.GetHitbox() {
				zombie.SubHealth(shoot.WeaponType.Power())
				game.Remove(Entry{shoot, j})
			}
			//Remove all zombies at zero health
			if zombie.GetStats().Health <= 0 {
				game.Remove(Entry{zombie, i})
				goto endloop
			}
		}

		zombie.Move(game)
		zombie.Attack(game)
		endloop:
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

}

func (game *Game) HandleCorpses(){
	for i := len(game.State.Players) - 1; i >= 0; i-- {
		player := game.State.Players[i]
		if player.Stats.Health <= 0 {
			//Remove player from game
			game.PlayerIds[player.Id] = false
			game.Remove(Entry{player, i})
		}
	}
	for i := len(game.State.Zombies) - 1; i >= 0; i-- {
		zombie := game.State.Zombies[i]
		if zombie.GetStats().Health <= 0 {
			//Remove zombie from game
			game.Remove(Entry{zombie, i})
		}
	}
	for i := len(game.State.Barrels) - 1; i >= 0; i-- {
		barrel := game.State.Barrels[i]
		if  barrel.Exploded{
			//Remove zombie from game
			game.Remove(Entry{barrel, i})
		}
	}
}

func (game *Game) HandleBarrels() {
	for i := range game.State.Barrels {
		barrel := &game.State.Barrels[i]
		for j := len(game.State.Shots) - 1; j >= 0; j-- {
			shot := game.State.Shots[j]
			if shot.GetPos().Sub(barrel.Pos).Len() < barrel.GetHitbox() {
				//Update objects
				barrel.Explode(&game.State)
				game.Remove(Entry{shot, j})
				break
			}
		}
	}
	barrelEntries := make([]Entry, 0, len(game.State.Barrels))
	for i,barrel := range game.State.Barrels{
		if barrel.Exploded{
			barrelEntries = append(barrelEntries, Entry{elem:barrel, index: i})
		}
	}
	game.Remove(barrelEntries...)
}

func (game *Game) Add(entities ...EntityI) {
	for _, entity := range entities {
		fmt.Println(helper.RealType(entity))
		switch entity.EntityType() {
		case BarrelE: game.State.Barrels = append(game.State.Barrels, entity.(Barrel))
		case ShotE: game.State.Shots = append(game.State.Shots, entity.(Shot))
		case ZombieE: game.State.Zombies = append(game.State.Zombies, entity.(ZombieI))
		case PlayerE: game.State.Players = append(game.State.Players, entity.(Player))
		case LootboxE: game.State.Lootboxes = append(game.State.Lootboxes, entity.(Lootbox))
		}
	}
	game.Updates.Add(entities...)
}


func (game *Game) Remove(entries ...Entry){
	shots := make([]Entry, 0,minInt(len(entries), len(game.State.Shots)))
	players := make([]Entry, 0,minInt(len(entries), len(game.State.Players)))
	zombies := make([]Entry, 0,minInt(len(entries), len(game.State.Zombies)))
	barrels := make([]Entry, 0,minInt(len(entries), len(game.State.Barrels)))
	lootboxes := make([]Entry, 0,minInt(len(entries), len(game.State.Lootboxes)))
	for _, entry := range entries {
		switch entry.elem.EntityType(){
		case ShotE: shots = append(shots, entry)
		case PlayerE: players = append(players, entry)
		case ZombieE: zombies = append(zombies, entry)
		case BarrelE: barrels = append(barrels, entry)
		case LootboxE: lootboxes = append(lootboxes, entry)
		}
		game.Updates.Remove(entry.elem)
	}
	sort.Sort(ByIndexDescending(shots))
	sort.Sort(ByIndexDescending(players))
	sort.Sort(ByIndexDescending(zombies))
	sort.Sort(ByIndexDescending(barrels))
	sort.Sort(ByIndexDescending(lootboxes))
	for _, entry := range shots {
		last := len(game.State.Shots) - 1
		game.State.Shots[entry.index] = game.State.Shots[last]
		game.State.Shots = game.State.Shots[:last]
	}
	for _, entry := range players {
		last := len(game.State.Players) - 1
		game.State.Players[entry.index] = game.State.Players[last]
		game.State.Players = game.State.Players[:last]
	}
	for _, entry := range zombies {
		last := len(game.State.Zombies) - 1
		game.State.Zombies[entry.index] = game.State.Zombies[last]
		game.State.Zombies = game.State.Zombies[:last]
	}
	for _, entry := range barrels {
		last := len(game.State.Barrels) - 1
		game.State.Barrels[entry.index] = game.State.Barrels[last]
		game.State.Barrels = game.State.Barrels[:last]
	}
	for _, entry := range lootboxes {
		last := len(game.State.Lootboxes) - 1
		game.State.Lootboxes[entry.index] = game.State.Lootboxes[last]
		game.State.Lootboxes = game.State.Lootboxes[:last]
	}
}

func (game *Game) NewZombie(vec pixel.Vec, zombieType Being) ZombieI {
	var zom ZombieI

	zombie := Zombie{
		Id:    xid.New().String(),
		Pos:   vec,
		Dir:   0,
		Stats: NewStats(zombieType),
		Type:  zombieType,
	}

	switch zombieType{
	case FASTZOMBIE:
		zom= &FastZombie{
			zombie,
		}
	case BOMBZOMBIE:
		barrel := NewBarrel(vec)
		z:=BombZombie{
			zombie,
			&barrel,
		}
		game.Add(*z.Barrel)
		zom=&z
	case SLOWZOMBIE:
		zom=&SlowZombie{
			zombie,
	}
	case ZOMBIE:
		zom=&zombie
}
	game.Add(zom)
	return zom
}