package model

import "github.com/gustavfjorder/pixel-head/config"

type Stats struct {
	Health    int
	Being	  int
}

const (
	HUMAN  = iota
	ZOMBIE
)

func NewStats(being int) (s Stats) {
	s.Being=being
	s.Health = s.GetMaxHealth()
	return
}

func (s Stats) GetMaxHealth() int{
	switch s.Being {
	case HUMAN:
		return 100
	case ZOMBIE:
		return 20
	}
	return 0
}

//Number of units per second
func (s Stats) GetMoveSpeed() (speed float64){
	switch s.Being {
	case HUMAN:
		speed = 400
	case ZOMBIE:
		speed =  50
	}
	return speed * config.Conf.ServerHandleSpeed.Seconds()
}

func (s Stats) GetPower() int {
	switch s.Being {
	case HUMAN:
		return 5
	case ZOMBIE:
		return 3
	}
	return 0
}