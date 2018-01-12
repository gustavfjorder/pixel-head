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
)

type AnimationHandler struct {
	win              *pixelgl.Window
	animations       map[string]Animation
	activeAnimations map[string]*Animation
	updateChan       <-chan model.Updates
	center           pixel.Vec
	state            model.State
	updates          []model.Updates
}

func NewAnimationHandler(updates <-chan model.Updates) (ah AnimationHandler) {
	ah.animations = Load("client/sprites", "", ANIM)
	ah.animations["bullet"], _ = LoadAnimation(config.Conf.BulletPath)
	ah.animations["barrel"], _ = LoadAnimation(config.Conf.BarrelPath)
	ah.animations["explosion"] = LoadSpriteSheet(1024/8, 1024/8, 8*8, config.Conf.ExplosionPath)
	ah.activeAnimations = make(map[string]*Animation)
	ah.center = pixel.ZV
	ah.updateChan = updates
	for k, v := range ah.animations {
		fmt.Println(k,v)
	}
	return
}

func (ah *AnimationHandler) SetWindow(win *pixelgl.Window){
	ah.win = win
}

func (ah AnimationHandler) Draw(state model.State) {
	ah.state = state
	ah.collectUpdates()
	ah.drawBarrels()
	ah.drawBulllets()
	ah.drawBarrels()
	ah.drawZombies()
	ah.drawPlayers()
}

func (ah AnimationHandler) collectUpdates() () {
	ah.updates = make([]model.Updates, 0)
	var update model.Updates
	for {
		select{
		case update = <-ah.updateChan:
			ah.updates = append(ah.updates,update )
		default:
			return
		}
	}
}

func (ah AnimationHandler) drawBarrels() {
	barrel := ah.animations["barrel"]
	for _, b := range ah.state.Barrels {
		transofrmation := pixel.IM.ScaledXY(pixel.ZV, pixel.V(0.5, 0.5)).Moved(b.Pos)
		barrel.Next().Draw(ah.win, transofrmation)
	}
}
func (ah AnimationHandler) drawBulllets() {
	bullet := ah.animations["bullet"]
	for _, shot := range ah.state.Shots {
		p := shot.GetPos()
		transformation := pixel.IM.Scaled(pixel.ZV, config.BulletScale).Rotated(pixel.ZV, shot.Angle-math.Pi/2).Moved(p)
		bullet.Next().Draw(ah.win, transformation)
	}
}
func (ah AnimationHandler) drawZombies() {
	for _, zombie := range ah.state.Zombies {
		v, ok := ah.activeAnimations[zombie.Id]
		prefix := Prefix("zombie", "walk")
		if !ok {
			newanim, ok := ah.animations[prefix]
			if ok {
				ah.activeAnimations[zombie.Id] = &newanim
				newanim.Start(config.Conf.AnimationSpeed)
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
				ah.activeAnimations[zombie.Id].ChangeAnimation(newanim, true)
				ah.activeAnimations[zombie.Id].Prefix = prefix
			}

		}
		if len(v.Sprites) > 0 {
			transformation := pixel.IM.Scaled(ah.center, config.ZombieScale).Rotated(ah.center, zombie.Dir).Moved(zombie.Pos)
			v.Next().Draw(ah.win, transformation)
		}
	}
}
func (ah AnimationHandler) drawPlayers() {
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
		anim, ok := ah.activeAnimations[player.Id]
		if !ok {
			newAnim, ok := ah.animations[prefix]
			if ok {
				newAnim.Start(config.Conf.AnimationSpeed)
				ah.activeAnimations[player.Id] = &newAnim
				newAnim.Prefix = prefix
			} else {
				continue
			}
			anim = &newAnim
		}
		if anim.Prefix != prefix {
			newAnim, found := ah.animations[prefix]
			if found {
				anim.Prefix = prefix
				anim.ChangeAnimation(newAnim, blocking)
			}
		}
		if len(anim.Sprites) > 0 {
			transformation := pixel.IM.Rotated(ah.center, player.Dir).Scaled(ah.center, config.HumanScale).Moved(player.Pos)
			anim.Next().Draw(ah.win, transformation)
		}
	}
}
