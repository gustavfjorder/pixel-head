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
}

func NewAnimationHandler(updates <-chan model.Updates) (ah AnimationHandler) {
	ah.animations = Load("client/sprites", "", ANIM)
	ah.animations["bullet"], _ = LoadAnimation(config.Conf.BulletPath)
	ah.animations["barrel"], _ = LoadAnimation(config.Conf.BarrelPath)
	ah.animations["explosion"] = LoadSpriteSheet(1024/8, 1024/8, 8*8, config.Conf.ExplosionPath)
	ah.activeAnimations = make(map[string]*Animation)
	ah.center = pixel.ZV
	ah.updateChan = updates

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
	ah.handleUpdates()
	ah.collectBarrels()
	ah.collectBulllets()
	ah.collectZombies()
	ah.collectPlayers()
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
					if _, present := ah.activeAnimations[entity.ID]; present{
						ah.activeAnimations[entity.ID].ChangeAnimation(ah.animations["explosion"], true, true)
					}
				}
			}
		default:
			return
		}
	}
}

func (ah AnimationHandler) collectBarrels() {
	//barrel := ah.animations["barrel"]
	//for _, b := range ah.updates {
	//	transofrmation := pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.5)).Moved(b.Pos)
	//	barrel.Next().Draw(ah.win, transofrmation)
	//}
}
func (ah AnimationHandler) collectBulllets() {
	bullet := ah.animations["bullet"]
	for _, shot := range ah.state.Shots {
		p := shot.GetPos()
		bullet.Transformation =  pixel.IM.Scaled(pixel.ZV, config.BulletScale).Rotated(pixel.ZV, shot.Angle-math.Pi/2).Moved(p)
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
		v.Transformation = pixel.IM.Scaled(ah.center, config.ZombieScale).Rotated(ah.center, zombie.Dir).Moved(zombie.Pos)
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
		anim.Transformation = pixel.IM.Rotated(ah.center, player.Dir).Scaled(ah.center, config.HumanScale).Moved(player.Pos)
		ah.activeAnimations[player.ID()] = anim
	}
}
