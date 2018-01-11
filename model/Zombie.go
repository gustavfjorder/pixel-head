package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/gustavfjorder/pixel-head/config"
	"time"
	"strconv"
)

type Zombie struct {
	Id          string
	Pos         pixel.Vec
	Dir         float64
	Attacking   bool
	Stats       Stats
	TargetId    string
}

var count = 0
var attackDelays = make(map[string]time.Duration)

func NewZombie(x, y float64) Zombie {
	count++
	return Zombie{
		Id:    strconv.Itoa(count),
		Pos:   pixel.V(x, y),
		Dir:   0,
		Stats: NewStats(ZOMBIE),
	}
}

func (zombie *Zombie) Move(players []Player) {
	if len(players) <= 0 {
		return
	}
	closestPlayer := players[0]

	for _, player := range players {
		if player.Pos.Sub(zombie.Pos).Len() < closestPlayer.Pos.Sub(zombie.Pos).Len() {
			closestPlayer = player
		}
	}
	if closestPlayer.Pos.Sub(zombie.Pos).Len() < zombie.GetRange()/2{
		return
	}

	a := angle(zombie.Pos, closestPlayer.Pos)

	dAngle := a - zombie.Dir
	if math.Abs(dAngle) > math.Pi {
		zombie.Dir += math.Copysign(math.Min(math.Abs(dAngle)-math.Pi, zombie.GetTurnSpeed()), -dAngle)
	} else {
		zombie.Dir += math.Copysign(math.Min(math.Abs(dAngle), zombie.GetTurnSpeed()), dAngle)
	}
	zombie.Dir = math.Mod(zombie.Dir, math.Pi*2)

	zombie.Pos = zombie.Pos.Add(pixel.V(zombie.Stats.GetMoveSpeed(), 0).Rotated(zombie.Dir))

}

func (zombie *Zombie) Attack(state State) {
	zombie.Attacking = false
	if zombie.TargetId != "" {
		if zombie.AttackDelay() > state.Timestamp {
			zombie.Attacking = true
			return
		} else {
			for i := range state.Players {
				player := &state.Players[i]
				if player.Id == zombie.TargetId {
					if zombie.Pos.Sub(player.Pos).Len() <= zombie.GetRange() &&
						math.Abs(zombie.angle(player.Pos)) <= zombie.GetMaxAttackAngle(){
						player.Stats.Health -= zombie.Stats.GetPower()
					}
					break
				}
			}
			zombie.TargetId = ""
		}
	}
	for _, player := range state.Players {
		if zombie.Pos.Sub(player.Pos).Len() <= zombie.GetRange() &&
			math.Abs(zombie.angle(player.Pos)) <= zombie.GetMaxAttackAngle(){
			zombie.Dir = angle(zombie.Pos, player.Pos)
			zombie.Attacking = true
			zombie.SetAttackDelay(state.Timestamp)
			zombie.TargetId = player.Id
			break
		}
	}
}

func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y-this.Y, other.X-this.X)
}

func (zombie Zombie) GetMaxAttackAngle() float64 {
	return math.Pi/3
}

func (zombie Zombie) angle(p pixel.Vec) float64{
	return angle(zombie.Pos, p) - zombie.Dir
}

func (zombie Zombie) GetRange() float64 {
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
	turnSpeed = math.Pi / 3

	return turnSpeed * config.Conf.ServerHandleSpeed.Seconds()
}

func (zombie Zombie) AttackDelay() time.Duration{
	return attackDelays[zombie.Id]
}

func (zombie *Zombie) SetAttackDelay(timestamp time.Duration) {
	attackDelays[zombie.Id] = zombie.GetAttackDelay() + timestamp
}
