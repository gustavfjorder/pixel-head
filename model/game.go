package model

import (
	"math/rand"
	"fmt"
	"github.com/faiface/pixel"
	"time"
	"sort"
	"github.com/rs/xid"
	"os"
	"reflect"
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
	game.Add(NewBarrel(pixel.V(500, 500)), NewBarrel(pixel.V(600, 600)), NewBarrel(pixel.V(700, 700)), NewBarrel(pixel.V(900, 900)), NewBarrel(pixel.V(1000, 1000)))
	for _, id := range ids {
		game.Add(NewPlayer(id))
		game.PlayerIds[id] = true
	}
	return
}

func (game *Game) PrepareLevel(end chan<- bool) {
	level := Levels[game.CurrentLevel]
	game.State.Zombies = make([]ZombieI, 0)
	waveticker := time.NewTicker(level.TimeBetweenWaves)
	zombieticker := time.NewTicker(level.TimeBetweenZombies)
	fmt.Println("num zom", level.NumberOfZombiesPerWave, "num waves", level.NumberOfWaves)

	for i := 0; i < level.NumberOfWaves; i++ {
		fmt.Println("i:", i)
		<-waveticker.C
		for j := 0; j < level.NumberOfZombiesPerWave; j++ {
			fmt.Println("j:", j)
			fmt.Println("i*level.NumberOfZombiesPerWave+j", i*level.NumberOfZombiesPerWave+j)
			fmt.Println(len(game.State.Zombies))
			var ZOM Being
			switch rand.Intn(4) {
			case 1:
				ZOM = ZOMBIE
			case 0:
				ZOM = BOMBZOMBIE
			case 2:
				ZOM = FASTZOMBIE
			case 3:
				ZOM = SLOWZOMBIE
			}
			game.NewZombie(game.CurrentMap.SpawnPoint[rand.Intn(len(game.CurrentMap.SpawnPoint))], ZOM)
			<-zombieticker.C
		}
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

var lastTime = time.Now()

func (game *Game) HandleLoot() {
	for i := len(game.State.Players) - 1; i >= 0; i-- {
		player := &game.State.Players[i]

		for j := len(game.State.Lootboxes) - 1; j >= 0; j-- {
			lootbox := game.State.Lootboxes[j]

			if PointFrom(player.Pos).Dist(PointFrom(lootbox.Pos)) < 30 {
				player.PickupLootbox(&lootbox)
				game.Remove(Entry{lootbox, j})
			}
		}
	}

	// Place lootboxes
	if lastTime.Add(time.Second * 10).Before(time.Now()) && float64(rand.Intn(100)) <= 3.8 {
		min := 0
		max := len(game.CurrentMap.LootPoints)

		lootPoint := rand.Intn(max-min) + min
		point := game.CurrentMap.LootPoints[lootPoint]

		if ! game.State.HasLootboxAt(point) {
			game.Add(NewLootbox(point.X, point.Y))

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
				zombie.Hit(shoot, &game.State)
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

func (game *Game) HandleCorpses() {
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
		if barrel.IsExploded(){
			//Remove zombie from game
			game.Remove(Entry{barrel, i})
		}
	}
}

func (game *Game) HandleBarrels() {
	for i := range game.State.Barrels {
		barrel := game.State.Barrels[i]
		for j := len(game.State.Shots) - 1; j >= 0; j-- {
			shot := game.State.Shots[j]
			if shot.GetPos().Sub(barrel.GetPos()).Len() < barrel.GetHitbox() {
				//Update objects
				barrel.Explode(&game.State)
				game.Remove(Entry{shot, j})
				break
			}
		}
	}

}

func (game *Game) Add(entities ...EntityI) {
	for _, entity := range entities {
		switch entity.(type) {
		case BarrelI:
			game.State.Barrels = append(game.State.Barrels, entity.(BarrelI))
		case Shot:
			game.State.Shots = append(game.State.Shots, entity.(Shot))
		case ZombieI:
			game.State.Zombies = append(game.State.Zombies, entity.(ZombieI))
		case Player:
			game.State.Players = append(game.State.Players, entity.(Player))
		case Lootbox:
			game.State.Lootboxes = append(game.State.Lootboxes, entity.(Lootbox))
		default: fmt.Fprintln(os.Stderr, "ADD: Unable to find:", entity, "with type", reflect.TypeOf(entity)); continue
		}
		game.Updates.Add(entity)
	}
}

func (game *Game) Remove(entries ...Entry) {
	shots := make([]Entry, 0, MinInt(len(entries), len(game.State.Shots)))
	players := make([]Entry, 0, MinInt(len(entries), len(game.State.Players)))
	zombies := make([]Entry, 0, MinInt(len(entries), len(game.State.Zombies)))
	barrels := make([]Entry, 0, MinInt(len(entries), len(game.State.Barrels)))
	lootboxes := make([]Entry, 0, MinInt(len(entries), len(game.State.Lootboxes)))
	for _, entry := range entries {
		switch entry.elem.(type) {
		case Shot: shots = append(shots, entry)
		case Player: players = append(players, entry)
		case ZombieI: zombies = append(zombies, entry)
		case BarrelI: barrels = append(barrels, entry)
		case Lootbox: lootboxes = append(lootboxes, entry)
		default:fmt.Fprintln(os.Stderr, "REMOVE: Unable to find:", entry, "with type:",reflect.TypeOf(entry.elem)); continue
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

	switch zombieType {
	case FASTZOMBIE:
		zom = &FastZombie{
			zombie,
		}
	case BOMBZOMBIE:
		barrel := NewBarrel(vec)
		game.Add(barrel)
		zom = &BombZombie{
			zombie,
			barrel,
		}
	case SLOWZOMBIE:
		zom = &SlowZombie{
			zombie,
		}
	case ZOMBIE:
		zom = &zombie
	}
	game.Add(zom)
	return zom
}
