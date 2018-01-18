package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/gustavfjorder/pixel-head/config"
	"time"
)

type Zombie struct {
	ZombieI
	Id        string
	Pos       pixel.Vec
	Dir       float64
	Attacking bool
	Stats     Stats
	TargetId  string
	Type      Being
}

type ZombieI interface {
	EntityI
	Move(game *Game)
	GetStats() Stats
	SetPos(pos pixel.Vec)
	SetHealth(health int)
	SubHealth(health int)
	Hit(shot Shot, state *State)
	IsAttacking() bool
	Attack(game *Game)
}

func (zombie Zombie) GetStats() Stats {
	return zombie.Stats
}

//func NewBombZombie(vec pixel.Vec)

func (zombie *Zombie) Move(game *Game) {
	if len(game.State.Players) <= 0 {
		return
	}
	closestPlayer := game.State.Players[0]

	for _, player := range game.State.Players {
		if player.Pos.Sub(zombie.GetPos()).Len() < closestPlayer.Pos.Sub(zombie.GetPos()).Len() {
			closestPlayer = player
		}
	}

	a := zombie.angle(closestPlayer.Pos)

	if math.Abs(a) > math.Pi {
		a = math.Copysign(math.Pi*2-math.Abs(a), -a)
	}
	zombie.Dir += math.Copysign(math.Min(math.Abs(a), zombie.GetTurnSpeed()), a)
	zombie.Dir = math.Mod(zombie.Dir, math.Pi*2)

	if closestPlayer.Pos.Sub(zombie.GetPos()).Len() > zombie.GetRange()/2 &&
		math.Abs(a) <= math.Pi/2 {
		zombie.SetPos(zombie.GetPos().Add(pixel.V(zombie.GetStats().GetMoveSpeed(), 0).Rotated(zombie.GetDir())))
	}
}

func (zombie *Zombie) Attack(game *Game) {
	zombie.Attacking = false
	if zombie.TargetId != "" {
		if zombie.AttackDelay(game) > game.State.Timestamp {
			zombie.Attacking = true
			return
		} else {
			for i := range game.State.Players {
				player := &game.State.Players[i]
				if player.Id == zombie.TargetId {
					if zombie.GetPos().Sub(player.GetPos()).Len() <= zombie.GetRange() &&
						math.Abs(zombie.angle(player.GetPos())) <= zombie.GetMaxAttackAngle() {
						player.Stats.Health -= zombie.GetStats().GetPower()
					}
					break
				}
			}
			zombie.TargetId = ""
		}
	}
	//check if any available targets
	for _, player := range game.State.Players {
		if zombie.GetPos().Sub(player.GetPos()).Len() <= zombie.GetRange() &&
			math.Abs(zombie.angle(player.GetPos())) <= zombie.GetMaxAttackAngle() {
			zombie.Dir = angle(zombie.GetPos(), player.GetPos())
			zombie.Attacking = true
			zombie.SetAttackDelay(game)
			zombie.TargetId = player.Id
			break

		}
	}
}

func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y-this.Y, other.X-this.X)
}

func (zombie Zombie) GetMaxAttackAngle() float64 {
	return math.Pi / 3
}

func (zombie Zombie) angle(p pixel.Vec) float64 {
	return angle(zombie.Pos, p) - zombie.Dir
}

func (zombie Zombie) GetRange() float64 {
	switch zombie.Type {
	case ZOMBIE:
		return 100
	case FASTZOMBIE:
		return 50
	case SLOWZOMBIE:
		return 100
	case BOMBZOMBIE:
		return 20
	}
	return 100
}

func (zombie Zombie) GetHitbox() float64 {
	return 50
}

//Time from attack is initiated till hit is calculated
func (zombie Zombie) GetAttackDelay() time.Duration {
	switch zombie.Stats.Being {
	case FASTZOMBIE: return time.Second/5
	default:return time.Second / 3
	}
}

//Radians per second
func (zombie Zombie) GetTurnSpeed() (turnSpeed float64) {
	switch zombie.Type {
	case ZOMBIE:
		turnSpeed = math.Pi / 3
	case FASTZOMBIE:
		turnSpeed = math.Pi * 2
	case SLOWZOMBIE:
		turnSpeed = math.Pi / 5
	case BOMBZOMBIE:
		turnSpeed = math.Pi/2
	}

	return turnSpeed * config.Conf.ServerHandleSpeed.Seconds()
}

func (zombie Zombie) AttackDelay(game *Game) time.Duration {
	return game.attackDelays[zombie.Id]
}

func (zombie *Zombie) SetAttackDelay(game *Game) {
	game.attackDelays[zombie.Id] = zombie.GetAttackDelay() + game.State.Timestamp
}

func (zombie Zombie) ID() string {
	return zombie.Id
}

func (zombie Zombie) EntityType() EntityType {
	return ZombieE
}

func (zombie Zombie) GetPos() pixel.Vec {
	return zombie.Pos
}

func (zombie Zombie) GetDir() float64 {
	return zombie.Dir
}

func (zombie *Zombie) SetPos(pos pixel.Vec) {
	zombie.Pos = pos
}

func (zombie *Zombie) SetHealth(health int) {
	zombie.Stats.Health = health
}

func (zombie Zombie) IsAttacking() bool {
	return zombie.Attacking
}

func (zombie *Zombie) SubHealth(health int) {
	zombie.Stats.Health -= health
}

func (zombie *Zombie) Hit(shot Shot, state *State) {
	zombie.SubHealth(shot.WeaponType.Power())
}

type FastZombie struct {
	Zombie
}
type BombZombie struct {
	Zombie
	Barrel BarrelI
}

func (zombie *BombZombie) Move(game *Game) {
	if zombie.Barrel.IsExploded(){
		zombie.SetHealth(0)
		return
	}
	zombie.Zombie.Move(game)
	zombie.Barrel.SetPos(zombie.Zombie.Pos)
}

func (zombie *BombZombie) Attack(game *Game){
	for _, player := range game.State.Players {
		if zombie.GetPos().Sub(player.GetPos()).Len() <= zombie.GetRange() {
			barrel := NewBarrel(zombie.GetPos())
			barrel.Explode(&game.State)
			zombie.SetHealth(0)
		}
	}
}

func (zombie *BombZombie) Hit(shot Shot, state *State) {
	zombie.SetHealth(0)
	zombie.Barrel.Explode(state)
}

type SlowZombie struct {
	Zombie
}
