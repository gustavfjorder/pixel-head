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
	activeAnimations map[string]Animation
	updateChan       <-chan model.Updates
	stateChan        <-chan model.State
	center           pixel.Vec
	state            model.State
	ticker           *time.Ticker
	me               model.Player
}


func NewAnimationHandler(updates <-chan model.Updates) (ah AnimationHandler) {
	spritePath := "client/sprites/"
	ah.animations = LoadAll(spritePath + "animations", spritePath + "images")
	ah.animations["explosion"] = NewAnimation( "explosion",
		LoadSpriteSheet(1024/8, 1024/8, 8*8, spritePath + "images/explosion/explosion.png"),
			Terminal)
	ah.activeAnimations = make(map[string]Animation)
	ah.center = pixel.ZV
	ah.updateChan = updates
	ah.me = model.NewPlayer(config.Conf.Id)
	ah.ticker = time.NewTicker(config.Conf.AnimationSpeed)
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
	for id, animation := range ah.activeAnimations {
		animation.Draw(ah.win)
		next := animation.Next()
		if next == nil {
			delete(ah.activeAnimations,id)
		} else {
			ah.activeAnimations[id] = next
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
					v, present := ah.activeAnimations[entity.ID]
					if present {
						prefix := Prefix("zombie", "death0"+strconv.Itoa(rand.Intn(2)+1))
						ah.activeAnimations[entity.ID] = v.ChangeAnimation(ah.animations[prefix])
					}
				case model.PlayerE:
					delete(ah.activeAnimations, entity.ID)
				case model.BarrelE:
					if anim, present := ah.activeAnimations[entity.ID]; present {
						exp := ah.animations["explosion"]
						exp.SetAnimationSpeed(time.Second/120)
						anim = anim.ChangeAnimation(exp)
						ah.activeAnimations[entity.ID] = anim
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
		barrel := ah.animations[Prefix("barrel","barrel")]
		barrel.SetTransformation(Transformation{Pos:b.Pos, Scale:1, Rotation:0})
		ah.activeAnimations[b.ID()] = barrel
	}
}
func (ah AnimationHandler) collectBulllets() {
	bullet := ah.animations[Prefix("bullet","bullet")]
	for _, shot := range ah.state.Shots {
		bullet.SetTransformation(Transformation{Scale:config.BulletScale, Pos:shot.GetPos(), Rotation:shot.Angle-math.Pi/2})
		bullet.Draw(ah.win)
	}
}
func (ah AnimationHandler) collectZombies() {
	for _, zombie := range ah.state.Zombies {
		var prefix string
		if zombie.Attacking {
			n := rand.Int()%3 + 1
			prefix = Prefix("zombie", "attack0"+strconv.Itoa(n))
		} else {
			prefix = Prefix("zombie", "walk")
		}
		v, ok := ah.activeAnimations[zombie.ID()]
		if !ok {
			newanim, _ := ah.animations[prefix]
			v = newanim
		}
		if prefix != v.Prefix() {
			if newanim, ok := ah.animations[prefix]; ok {
				v = v.ChangeAnimation(newanim)
			}
		}
		v.SetTransformation(Transformation{Pos:zombie.Pos, Rotation:zombie.Dir, Scale:config.ZombieScale})
		ah.activeAnimations[zombie.ID()] = v
	}
}
func (ah AnimationHandler) collectPlayers() {
	for _, player := range ah.state.Players {
		movement := "idle"
		switch player.Action {
		case model.RELOAD:
			movement = "reload"
		case model.SHOOT:
			movement = "shoot"
		case model.MELEE:
			movement = "melee"
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
			anim = newAnim
		}
		if anim.Prefix() != prefix {
			newAnim, found := ah.animations[prefix]
			if found {
				anim = anim.ChangeAnimation(newAnim)
			}
		}

		anim.SetTransformation(Transformation{Pos:player.Pos, Scale:config.HumanScale, Rotation:player.Dir})
		ah.activeAnimations[player.ID()] = anim
	}
}
