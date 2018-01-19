package client

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/gustavfjorder/pixel-head/config"
	. "github.com/gustavfjorder/pixel-head/client/animation"
	"github.com/faiface/pixel"
	"strconv"
	"math/rand"
	"time"
	"github.com/pkg/errors"
	"fmt"
	"os"
)


type AnimationHandler struct {
	win              *pixelgl.Window
	animations       map[string]Animation
	activeAnimations []map[string]Animation
	tracked          map[string]model.EntityI
	updateChan       <-chan model.Updates
	stateChan        <-chan model.State
	center           pixel.Vec
	state            model.State
	ticker           *time.Ticker
	me               model.Player
	loaded           bool
}

func NewAnimationHandler() (ah AnimationHandler) {
	ah.Clear()
	ah.tracked = make(map[string]model.EntityI)
	ah.center = pixel.ZV
	ah.me = model.NewPlayer(config.ID, pixel.V(0,0))
	ah.ticker = time.NewTicker(config.Conf.AnimationSpeed)
	return
}

func (ah *AnimationHandler) Load() {
	spritePath := "assets/sprites/"
	ah.animations = LoadAll(spritePath+"animations", spritePath+"images")
	ah.animations["explosion"] = NewAnimation("explosion",
		LoadSpriteSheet(1024/8, 1024/8, 8*8, spritePath+"images/explosion/explosion.png"),
		Terminal)

	ah.loaded = true
}

func Layer(entityType model.EntityType) int{
	switch entityType {
	case model.BarrelE: return 3
	case model.PlayerE: return 4
	case model.ZombieE: return 1
	case model.LootboxE: return 0
	case model.ShotE: return 2
	}
	return 0
}

func (ah *AnimationHandler) SetWindow(win *pixelgl.Window) {
	ah.win = win
}

func (ah *AnimationHandler) SetUpdateChan(ch chan model.Updates) {
	ah.updateChan = ch
}

func (ah AnimationHandler) Draw(state model.State) {
	if ! ah.loaded {
		return
	}

	ah.state = state
	GetPlayer(state.Players, &ah.me)
	ah.handleUpdates()
	ah.collectZombies()
	ah.collectPlayers()
	ah.handleTracked()
	for i, animations := range ah.activeAnimations {
		for id, animation := range animations {
			animation.Draw(ah.win)
			next := animation.Next()
			if next == nil {
				delete(ah.activeAnimations[i], id)
			} else {
				ah.activeAnimations[i][id] = next
			}
		}
	}
	ah.DrawAbilities()
	ah.DrawHealthbar()
}

func (ah *AnimationHandler) Clear(){
	ah.activeAnimations = make([]map[string]Animation,5)
	for i := range ah.activeAnimations {
		ah.activeAnimations[i] = make(map[string]Animation)
	}
}

func (ah AnimationHandler) handleUpdates() () {
	var update model.Updates
	for {
		select {
		case update = <-ah.updateChan:
			for _, entity := range update.Removed {
				switch entity.EntityType {
				case model.ShotE:
					delete(ah.activeAnimations[Layer(entity.EntityType)], entity.ID)
					delete(ah.tracked, entity.ID)

				case model.ZombieE:
					zombie, present := ah.activeAnimations[Layer(entity.EntityType)][entity.ID]
					if present {
						prefix := Prefix("zombie", "death0"+strconv.Itoa(rand.Intn(2)+1))
						animation, err := ah.Get(prefix)
						if err != nil {
							fmt.Fprint(os.Stderr, err.Error())
							delete(ah.activeAnimations[entity.EntityType], entity.ID)
							continue
						}
						ah.activeAnimations[Layer(entity.EntityType)][entity.ID] = zombie.ChangeAnimation(animation)
					}
				case model.PlayerE:
					delete(ah.activeAnimations[Layer(entity.EntityType)], entity.ID) //todo:No death animation
				case model.BarrelE:
					if barrel, present := ah.activeAnimations[Layer(entity.EntityType)][entity.ID]; present {
						exp, err := ah.Get("explosion")
						if err != nil {
							fmt.Fprint(os.Stderr, err.Error())
							delete(ah.activeAnimations[Layer(entity.EntityType)],entity.ID)
							continue
						}
						exp.SetAnimationSpeed(time.Second / 120)
						exp = barrel.ChangeAnimation(exp)
						t := exp.GetTransformation()
						y := exp.CurrentSprite().Picture().Bounds().Max.Y
						t.Scale = model.Barrel{}.GetRange()*8/y  //Times 8 as the spritesheet is bigger
						t.Pos = t.Pos.Add(pixel.V(0,model.Barrel{}.GetRange()/4))
						ah.activeAnimations[Layer(entity.EntityType)][entity.ID] = exp
					}
				case model.LootboxE:
					delete(ah.activeAnimations[Layer(entity.EntityType)], entity.ID)
				}
			}
			for _, entity := range update.Added {
				transformation := Transformation{Pos: entity.GetPos(), Scale: 1, Rotation: entity.GetDir()}
				switch entity.EntityType() {
				case model.BarrelE:
					barrel, err := ah.Get("barrel", "barrel")
					if err != nil {fmt.Fprint(os.Stderr, err.Error());continue}
					transformation.Scale = entity.GetHitbox() / barrel.CurrentSprite().Picture().Bounds().Max.X
					barrel.SetTransformation(transformation)
					ah.activeAnimations[Layer(entity.EntityType())][entity.ID()] = barrel
				case model.ShotE:
					bullet, err := ah.Get("bullet", "bullet")
					if err != nil {fmt.Fprint(os.Stderr,err.Error()); continue}
					transformation.Scale = config.BulletScale
					bullet.SetTransformation(transformation)
					ah.activeAnimations[Layer(entity.EntityType())][entity.ID()] = bullet
					ah.tracked[entity.ID()] = entity
				case model.ZombieE:
					zombie, err := ah.Get("zombie", "walk")
					if err != nil {fmt.Fprint(os.Stderr,err.Error()); continue}
					ah.activeAnimations[Layer(entity.EntityType())][entity.ID()] = zombie
				case model.PlayerE:
					player, err := ah.Get("survivor", "knife", "idle")
					if err != nil {fmt.Fprint(os.Stderr,err.Error()); continue}
					ah.activeAnimations[Layer(entity.EntityType())][entity.ID()] = player
				case model.LootboxE:
					lootbox, err := ah.Get("lootbox", "lootbox")
					if err != nil {fmt.Fprint(os.Stderr,err.Error()); continue}
					t := lootbox.GetTransformation()
					t.Pos = entity.GetPos()
					t.Scale = 0.2
					ah.activeAnimations[Layer(entity.EntityType())][entity.ID()] = lootbox
				}
			}
		default:
			return
		}
	}
}

func (ah AnimationHandler) handleTracked(){
	for _, entity := range ah.tracked {
		if anim, ok := ah.activeAnimations[Layer(entity.EntityType())][entity.ID()]; ok {
			anim.SetPos(entity.GetPos())
		}
	}
}

func (ah AnimationHandler) collectZombies() {
	for _, zombie := range ah.state.Zombies {
		animation, ok := ah.activeAnimations[Layer(zombie.EntityType())][zombie.ID()]
		if !ok {
			continue
		}
		var prefix string
		if zombie.IsAttacking() {
			n := rand.Int()%3 + 1
			prefix = Prefix("zombie", "attack0"+strconv.Itoa(n))
		} else if zombie.GetStats().Being == model.FASTZOMBIE {
			prefix = Prefix("zombie", "run")
		} else {
			prefix = Prefix("zombie", "walk")
		}
		if prefix != animation.Prefix() {
			zombieAnimation, err := ah.Get(prefix)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			} else {
				animation = animation.ChangeAnimation(zombieAnimation)
			}
		}
		switch zombie.(type) {
		case *model.FastZombie:
			animation.SetAnimationSpeed(time.Second/100)
		case *model.BombZombie:
			bz := zombie.(*model.BombZombie)
			if barrel, present :=ah.activeAnimations[Layer(bz.Barrel.EntityType())][bz.Barrel.ID()];present {
				barrel.SetPos(bz.GetPos())
			}
		}
		if zombie.GetStats().Being == model.FASTZOMBIE {
		} else {
			animation.SetAnimationSpeed(time.Second/30)
		}
		animation.SetTransformation(Transformation{Pos: zombie.GetPos(), Rotation: zombie.GetDir(), Scale: config.ZombieScale})
		ah.activeAnimations[Layer(zombie.EntityType())][zombie.ID()] = animation
	}
}

func (ah AnimationHandler) collectPlayers() {
	for _, player := range ah.state.Players {
		anim, ok := ah.activeAnimations[Layer(player.EntityType())][player.ID()]
		if !ok {
			continue
		}
		movement := "idle"
		switch player.Action {
		case model.RELOAD: if player.WeaponType != model.KNIFE{	movement = "reload"	}
		case model.SHOOT: movement = "shoot"
		case model.MELEE: movement = "melee"
		case model.IDLE: movement = "move"
		}
		prefix := Prefix("survivor", player.WeaponType.Name(), movement)
		if prefix != anim.Prefix() {
			player, err := ah.Get(prefix)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			} else {
				anim = anim.ChangeAnimation(player)
			}
		}
		anim.SetTransformation(Transformation{Pos: player.Pos, Scale: config.HumanScale, Rotation: player.Dir})
		ah.activeAnimations[Layer(player.EntityType())][player.ID()] = anim
	}
}

func (ah AnimationHandler) Get(prefix ...string) (Animation, error){
	if animation, present := ah.animations[Prefix(prefix...)]; present {
		return animation.Copy(), nil
	}
	return nil, errors.New("Unable to find animation" + Prefix(prefix...))
}
