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
)

type AnimationHandler struct {
	win              *pixelgl.Window
	animations       map[string]Animation
	activeAnimations map[string]Animation
	tracked          map[string]model.EntityI
	updateChan       <-chan model.Updates
	stateChan        <-chan model.State
	center           pixel.Vec
	state            model.State
	ticker           *time.Ticker
	me               model.Player
}

func NewAnimationHandler() (ah AnimationHandler) {
	spritePath := "client/sprites/"
	ah.animations = LoadAll(spritePath+"animations", spritePath+"images")
	ah.animations["explosion"] = NewAnimation("explosion",
		LoadSpriteSheet(1024/8, 1024/8, 8*8, spritePath+"images/explosion/explosion.png"),
		Terminal)
	ah.activeAnimations = make(map[string]Animation)
	ah.tracked = make(map[string]model.EntityI)
	ah.center = pixel.ZV
	ah.me = model.NewPlayer(config.ID)
	ah.ticker = time.NewTicker(config.Conf.AnimationSpeed)
	return
}

func (ah *AnimationHandler) SetWindow(win *pixelgl.Window) {
	ah.win = win
}

func (ah *AnimationHandler) SetUpdateChan(ch chan model.Updates) {
	ah.updateChan = ch
}

func (ah AnimationHandler) Draw(state model.State) {
	ah.state = state
	GetPlayer(state.Players, &ah.me)
	ah.handleUpdates()
	ah.collectZombies()
	ah.collectPlayers()
	ah.handleTracked()
	for id, animation := range ah.activeAnimations {
		animation.Draw(ah.win)
		next := animation.Next()
		if next == nil {
			delete(ah.activeAnimations, id)
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
					delete(ah.tracked, entity.ID)

				case model.ZombieE:
					v, present := ah.activeAnimations[entity.ID]
					if present {
						prefix := Prefix("zombie", "death0"+strconv.Itoa(rand.Intn(2)+1))
						ah.activeAnimations[entity.ID] = v.ChangeAnimation(ah.Get(prefix))
					}
				case model.PlayerE:
					delete(ah.activeAnimations, entity.ID)
				case model.BarrelE:
					if anim, present := ah.activeAnimations[entity.ID]; present {
						exp := ah.Get("explosion")
						exp.SetAnimationSpeed(time.Second / 120)
						anim = anim.ChangeAnimation(exp)
						anim.SetScale(model.Barrel{}.GetRange()*8/exp.CurrentSprite().Picture().Bounds().Max.X)
						ah.activeAnimations[entity.ID] = anim
					}
				}
			}
			for _, entity := range update.Added {
				transformation := Transformation{Pos: entity.GetPos(), Scale: 1, Rotation: entity.GetDir()}
				switch entity.EntityType() {
				case model.BarrelE:
					barrel := ah.Get("barrel", "barrel")
					transformation.Scale = entity.GetHitbox()*2 / barrel.CurrentSprite().Picture().Bounds().Max.X
					barrel.SetTransformation(transformation)
					ah.activeAnimations[entity.ID()] = barrel
				case model.ShotE:
					bullet := ah.Get("bullet", "bullet")
					transformation.Scale = config.BulletScale
					bullet.SetTransformation(transformation)
					ah.activeAnimations[entity.ID()] = bullet
					ah.tracked[entity.ID()] = entity
				case model.ZombieE:
					ah.activeAnimations[entity.ID()] = ah.Get("zombie", "walk")
				case model.PlayerE:
					ah.activeAnimations[entity.ID()] = ah.Get("survivor", "knife", "idle")
				}
			}
		default:
			return
		}
	}
}

func (ah AnimationHandler) handleTracked(){
	for _, entity := range ah.tracked {
		if anim, ok := ah.activeAnimations[entity.ID()]; ok {
			anim.SetPos(entity.GetPos())
		}
	}
}

func (ah AnimationHandler) collectBulllets() {
	for _, shot := range ah.state.Shots {
		if anim, present := ah.activeAnimations[shot.ID()]; present {
			anim.SetTransformation(Transformation{Scale: config.BulletScale, Pos: shot.GetPos(), Rotation: shot.GetDir()})
		}
	}
}
func (ah AnimationHandler) collectZombies() {
	for _, zombie := range ah.state.Zombies {
		v, ok := ah.activeAnimations[zombie.ID()]
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
		if prefix != v.Prefix() {
			v = v.ChangeAnimation(ah.Get(prefix))
		}
		if zombie.GetStats().Being == model.FASTZOMBIE {
			v.SetAnimationSpeed(time.Second/100)
		} else {
			v.SetAnimationSpeed(time.Second/30)
		}
		v.SetTransformation(Transformation{Pos: zombie.GetPos(), Rotation: zombie.GetDir(), Scale: config.ZombieScale})
		ah.activeAnimations[zombie.ID()] = v
	}
}
func (ah AnimationHandler) collectPlayers() {
	for _, player := range ah.state.Players {
		anim, ok := ah.activeAnimations[player.ID()]
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
			anim = anim.ChangeAnimation(ah.Get(prefix))
		}
		anim.SetTransformation(Transformation{Pos: player.Pos, Scale: config.HumanScale, Rotation: player.Dir})
		ah.activeAnimations[player.ID()] = anim
	}
}

func (ah AnimationHandler) Get(prefix ...string) (Animation) {
	return ah.animations[Prefix(prefix...)].Copy()
}
