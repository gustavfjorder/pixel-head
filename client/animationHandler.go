package client

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gustavfjorder/pixel-head/model"
	"github.com/gustavfjorder/pixel-head/config"
	"github.com/faiface/pixel"
	"math"
	"fmt"
	"strconv"
	"math/rand"
	"time"
	"os"
)

type AnimationHandler struct {
	win              *pixelgl.Window
	animations       map[string]Animation
	activeAnimations map[string]*Animation
	updateChan       <-chan model.Updates
	stateChan        <-chan model.State
	center           pixel.Vec
	state            model.State
	ticker           *time.Ticker
	me               model.Player
}

func NewAnimationHandler(updates <-chan model.Updates) (ah AnimationHandler) {
	ah.animations = Load("client/sprites", "", ANIM)
	for prefix, anim := range Load(config.Conf.AbilityPath, "", IMG) {
		ah.animations[prefix] = anim
	}
	ah.animations["bullet"], _ = LoadAnimation(config.Conf.BulletPath)
	ah.animations["barrel"], _ = LoadAnimation(config.Conf.BarrelPath)
	ah.animations["explosion"] = LoadSpriteSheet(1024/8, 1024/8, 8*8, config.Conf.ExplosionPath)
	ah.activeAnimations = make(map[string]*Animation)
	ah.center = pixel.ZV
	ah.updateChan = updates
	ah.me = model.NewPlayer(config.Conf.Id)
	ah.ticker = time.NewTicker(config.Conf.AnimationSpeed)
	for k, v := range ah.animations {
		fmt.Println(k, v)
	}
	return
}

func (ah *AnimationHandler) SetWindow(win *pixelgl.Window) {
	ah.win = win
}

func (ah AnimationHandler) Draw(state model.State) {
	ah.state = state
	GetPlayer(state.Players, &ah.me)
	ah.collectBarrels()
	ah.collectBulllets()
	ah.collectZombies()
	ah.collectPlayers()
	ah.handleUpdates()
	nextFrame := false
	select {
	case <-ah.ticker.C:
		nextFrame = true
	default:
		break
	}
	for id, animation := range ah.activeAnimations {
		animation.Draw(ah.win)
		if nextFrame {
			animation.Next()
		}
		if animation.Terminal && animation.Finished {
			delete(ah.activeAnimations, id)
			continue
		}
	}
	ah.DrawAbilities()
	ah.DrawHealthbar()
}

func (ah AnimationHandler) handleUpdates() () {
	var update model.Updates
	for {
		select {
		case update = <-ah.updateChan:
			for _, entity := range update.Removed {
				switch entity.EntityType {
				case model.ShotE:
					delete(ah.activeAnimations, entity.ID)
				case model.ZombieE:
					_, present := ah.activeAnimations[entity.ID]
					if present {
						prefix := Prefix("zombie", "death0"+strconv.Itoa(rand.Intn(2)+1))
						ah.activeAnimations[entity.ID].ChangeAnimation(ah.animations[prefix], true, true)
					}
				case model.PlayerE:
					delete(ah.activeAnimations, entity.ID)
				case model.BarrelE:
					fmt.Println(entity.ID)
					if anim, present := ah.activeAnimations[entity.ID]; present {
						exp := ah.animations["explosion"]
						exp.Terminal = true
						exp.Blocking = true
						exp.Finished = false
						exp.Cur = 0
						exp.Scale = 20
						exp.Pos = anim.Pos
						ah.activeAnimations[entity.ID] = &exp
					}
				}
			}
		default:
			return
		}
	}
}

func (ah AnimationHandler) collectBarrels() {
	for _, b := range ah.state.Barrels {
		barrel := ah.animations["barrel"]
		barrel.Pos = b.Pos
		barrel.Scale = b.GetHitBox()*2/barrel.GetCurrentSprite().Picture().Bounds().Max.X
		ah.activeAnimations[b.ID()] = &barrel
	}
}
func (ah AnimationHandler) collectBulllets() {
	bullet := ah.animations["bullet"]
	for _, shot := range ah.state.Shots {
		bullet.Scale = config.BulletScale
		bullet.Pos = shot.GetPos()
		bullet.Rotation = shot.Angle-math.Pi/2
		bullet.Draw(ah.win)
	}
}
func (ah AnimationHandler) collectZombies() {
	for _, zombie := range ah.state.Zombies {
		v, ok := ah.activeAnimations[zombie.ID()]
		prefix := Prefix("zombie", "walk")
		if !ok {
			newanim, ok := ah.animations[prefix]
			if ok {
				ah.activeAnimations[zombie.ID()] = &newanim
				v = &newanim
			} else {
				fmt.Println("Did not find animation:", prefix)
				continue
			}
		}
		if zombie.Attacking {
			n := rand.Int()%3 + 1
			prefix = Prefix("zombie", "attack0"+strconv.Itoa(n))
		}

		if prefix != v.Prefix {
			if newanim, ok := ah.animations[prefix]; ok {
				ah.activeAnimations[zombie.ID()].ChangeAnimation(newanim, true, false)
				ah.activeAnimations[zombie.ID()].Prefix = prefix
			}

		}
		v.Pos = zombie.Pos
		v.Rotation = zombie.Dir
		v.Scale = config.ZombieScale
		ah.activeAnimations[zombie.ID()] = v
	}
}
func (ah AnimationHandler) collectPlayers() {
	for _, player := range ah.state.Players {
		movement := "idle"
		blocking := false
		switch player.Action {
		case model.RELOAD:
			movement = "reload"
			blocking = true
		case model.SHOOT:
			movement = "shoot"
			blocking = true
		case model.MELEE:
			movement = "melee"
			blocking = true
		case model.IDLE:
			movement = "move"
		}

		prefix := Prefix("survivor", player.WeaponType.Name(), movement)
		anim, ok := ah.activeAnimations[player.ID()]
		if !ok {
			newAnim, ok := ah.animations[prefix]
			if !ok {
				fmt.Fprint(os.Stderr, "Invalid animation present")
				continue
			}
			anim = &newAnim
		}
		if anim.Prefix != prefix {
			newAnim, found := ah.animations[prefix]
			if found {
				anim.Prefix = prefix
				anim.ChangeAnimation(newAnim, blocking, false)
			}
		}
		anim.Scale = config.HumanScale
		anim.Rotation = player.Dir
		anim.Pos = player.Pos
		ah.activeAnimations[player.ID()] = anim
	}
}
