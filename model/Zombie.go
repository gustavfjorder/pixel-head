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

	if dAngle > 0 {
		if dAngle > math.Pi {
			zombie.Dir += math.Max(math.Pi - dAngle, -zombie.GetTurnSpeed())
		} else {
			zombie.Dir += math.Min(dAngle, zombie.GetTurnSpeed())
		}

	} else {
		if dAngle < -math.Pi {
			zombie.Dir += math.Min(math.Pi - dAngle, zombie.GetTurnSpeed())
		} else {
			zombie.Dir += math.Max(dAngle, -zombie.GetTurnSpeed())
		}
	}
	for zombie.Dir >= math.Pi*2{
		zombie.Dir -= math.Pi*2
	}
	for zombie.Dir <= -math.Pi*2{
		zombie.Dir += math.Pi*2
	}

	zombie.Pos = zombie.Pos.Add(pixel.V(zombie.Stats.MoveSpeed, 0).Rotated(zombie.Dir))

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
						player.Stats.Health -= zombie.Stats.Power
					}
					break
				}
			}
			zombie.TargetId = ""
		}
	}
	for _, player := range players {
		if zombie.Pos.Sub(player.Pos).Len() <= zombie.GetRange() {
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
