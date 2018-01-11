package model

import (
	"github.com/faiface/pixel"
	"math"
	"github.com/rs/xid"
)

type Zombie struct {
	Id          string
	Pos         pixel.Vec
	Dir         float64
	Attacking   bool
	Stats       Stats
	AttackDelay int
	TurnDelay   int
	TargetId    string
}

func NewZombie(x, y float64) Zombie {
	return Zombie{
		Id:    xid.New().String(),
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

func (zombie *Zombie) Attack(players []Player) {
	zombie.Attacking = false
	if zombie.TargetId != "" {
		if zombie.AttackDelay > 0 {
			zombie.AttackDelay--
			zombie.Attacking = true
			return
		} else {
			for i := range players {
				player := &players[i]
				if player.Id == zombie.TargetId {
					if zombie.Pos.Sub(player.Pos).Len() <= zombie.GetRange() {
						player.Stats.Health -= zombie.Stats.GetPower()
					}
					break
				}
			}
			zombie.TargetId = ""
		}
	}
	for _, player := range players {
		if zombie.Pos.Sub(player.Pos).Len() <= zombie.GetRange() &&
			math.Abs(zombie.Dir - angle(zombie.Pos, player.Pos)) <= math.Pi/2 {
			zombie.Dir = angle(zombie.Pos, player.Pos)
			zombie.Attacking = true
			zombie.AttackDelay = zombie.GetAttackDelay()
			zombie.TargetId = player.Id
			break
		}
	}
}

func angle(this pixel.Vec, other pixel.Vec) float64 {
	return math.Atan2(other.Y-this.Y, other.X-this.X)
}

func (zombie Zombie) GetRange() float64 {
	return 60
}

func (zombie Zombie) GetHitbox() float64 {
	return 50
}

func (zombie Zombie) GetAttackDelay() int {
	return 10
}

func (zombie Zombie) GetTurnSpeed() float64 {
	return math.Pi / 50
}
