package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/gustavfjorder/pixel-head/config"
	"time"
	"github.com/rs/xid"
	"fmt"
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
	IsAttacking() bool
	Attack(game *Game)
}

func (zombie Zombie) GetStats() Stats {
	return zombie.Stats
}

var attackDelays = make(map[string]time.Duration)

func NewZombie(vec pixel.Vec, zombieType Being) ZombieI {
	zombie := &Zombie{
		Id:    xid.New().String(),
		Pos:   vec,
		Dir:   0,
		Stats: NewStats(zombieType),
		Type:  zombieType,
	}

	return zombie
}

//func NewBombZombie(vec pixel.Vec)

func (zombie *Zombie) Move(game *Game) {
	fmt.Println("Moving", zombie.Id, "To", zombie.Pos)
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
	fmt.Println("TO", zombie.Pos)
}

func (zombie *Zombie) Attack(game *Game)  {
	zombie.Attacking = false
	if zombie.TargetId != "" {
		if zombie.AttackDelay() > Timestamp {
			zombie.Attacking = true
			return
		} else {
			for i := range game.State.Players {
				player := &game.State.Players[i]
				if player.Id == zombie.TargetId {
					if zombie.GetPos().Sub(player.GetPos()).Len() <= zombie.GetRange() &&
						math.Abs(zombie.angle(player.GetPos())) <= zombie.GetMaxAttackAngle() {
						player.Stats.Health -= zombie.GetStats().GetPower()
						fmt.Println("normal attack")
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
			if zombie.Type == BOMBZOMBIE {
				barrel := NewBarrel(zombie.GetPos())
				barrel.Explode(&game.State)
				fmt.Println("bomb attack")
				fmt.Println(zombie.Id)
				zombie.SetHealth(0)
				break
			} else {
				zombie.Dir = angle(zombie.GetPos(), player.GetPos())
				zombie.Attacking = true
				zombie.SetAttackDelay()
				zombie.TargetId = player.Id
				break
			}
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
	return time.Second / 5
}

//Radians per second
func (zombie Zombie) GetTurnSpeed() (turnSpeed float64) {
	switch zombie.Type {
	case ZOMBIE:
		turnSpeed = math.Pi * 50
	case FASTZOMBIE:
		turnSpeed = math.Pi * 2000
	case SLOWZOMBIE:
		turnSpeed = math.Pi / 3
	case BOMBZOMBIE:
		turnSpeed = math.Pi * 50
	}

	return turnSpeed * config.Conf.ServerHandleSpeed.Seconds()
}

func (zombie Zombie) AttackDelay() time.Duration {
	return attackDelays[zombie.Id]
}

func (zombie *Zombie) SetAttackDelay() {
	attackDelays[zombie.Id] = zombie.GetAttackDelay() + Timestamp
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

func (zombie *Zombie) SubHealth(health int){
	zombie.Stats.Health-=health
}